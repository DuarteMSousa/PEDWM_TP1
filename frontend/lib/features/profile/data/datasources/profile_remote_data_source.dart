import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/match_history_entry.dart';
import '../../domain/entities/match_replay.dart';
import '../../domain/entities/profile.dart';

class ProfileRemoteDataSource {
  ProfileRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<Profile> fetchProfile(String userId) async {
    final response = await _graphqlService.query(
      document: '''
        query Profile(\$userId: ID!) {
          user(id: \$userId) {
            id
            username
          }
          userStats(userId: \$userId) {
            games
            wins
          }
          gameHistory(userId: \$userId) {
            gameId
            roomId
            playedAt
            won
            winnerTeamId
            myTeamId
            myScore
            opponentScore
            finalScores {
              teamId
              points
            }
          }
        }
      ''',
      variables: <String, dynamic>{'userId': userId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for profile query.');
    }

    final userPayload = data['user'];
    final statsPayload = data['userStats'];

    if (userPayload is! Map<String, dynamic>) {
      throw AppException('User profile was not returned by GraphQL.');
    }

    final nickname = userPayload['username']?.toString() ?? userId;
    final matchesPlayed = (statsPayload is Map<String, dynamic>)
        ? (statsPayload['games'] as num?)?.toInt() ?? 0
        : 0;
    final wins = (statsPayload is Map<String, dynamic>)
        ? (statsPayload['wins'] as num?)?.toInt() ?? 0
        : 0;
    final matchHistory = _parseMatchHistory(data['gameHistory']);

    return Profile(
      userId: userId,
      nickname: nickname,
      matchesPlayed: matchesPlayed,
      wins: wins,
      matchHistory: matchHistory,
    );
  }

  Future<MatchReplay?> fetchReplay(String gameId) async {
    final response = await _graphqlService.query(
      document: '''
        query Replay(\$gameId: ID!) {
          gameReplay(gameId: \$gameId) {
            gameId
            roomId
            events {
              id
              type
              sequence
              timestamp
              payload
            }
          }
        }
      ''',
      variables: <String, dynamic>{'gameId': gameId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for replay query.');
    }

    final replayPayload = data['gameReplay'];
    if (replayPayload == null) {
      return null;
    }
    if (replayPayload is! Map<String, dynamic>) {
      throw AppException('Replay was not returned by GraphQL.');
    }

    final events = <ReplayEvent>[];
    final rawEvents = replayPayload['events'];
    if (rawEvents is List) {
      for (final item in rawEvents) {
        if (item is! Map<String, dynamic>) {
          continue;
        }

        events.add(
          ReplayEvent(
            id: item['id']?.toString() ?? '',
            type: item['type']?.toString() ?? 'UNKNOWN',
            sequence: (item['sequence'] as num?)?.toInt() ?? 0,
            timestamp:
                DateTime.tryParse(item['timestamp']?.toString() ?? '') ??
                DateTime.fromMillisecondsSinceEpoch(0, isUtc: true),
            payload: item['payload']?.toString() ?? '{}',
          ),
        );
      }
    }

    events.sort((a, b) => a.sequence.compareTo(b.sequence));

    return MatchReplay(
      gameId: replayPayload['gameId']?.toString() ?? gameId,
      roomId: replayPayload['roomId']?.toString() ?? '',
      events: events,
    );
  }

  List<MatchHistoryEntry> _parseMatchHistory(dynamic rawHistory) {
    if (rawHistory is! List) {
      return const <MatchHistoryEntry>[];
    }

    final history = <MatchHistoryEntry>[];
    for (final item in rawHistory) {
      if (item is! Map<String, dynamic>) {
        continue;
      }

      final finalScores = <String, int>{};
      final rawFinalScores = item['finalScores'];
      if (rawFinalScores is List) {
        for (final score in rawFinalScores) {
          if (score is! Map<String, dynamic>) {
            continue;
          }
          final teamId = score['teamId']?.toString() ?? '';
          final points = (score['points'] as num?)?.toInt();
          if (teamId.isEmpty || points == null) {
            continue;
          }
          finalScores[teamId] = points;
        }
      }

      history.add(
        MatchHistoryEntry(
          gameId: item['gameId']?.toString() ?? '',
          roomId: item['roomId']?.toString() ?? '',
          playedAt:
              DateTime.tryParse(item['playedAt']?.toString() ?? '') ??
              DateTime.fromMillisecondsSinceEpoch(0, isUtc: true),
          won: item['won'] == true,
          winnerTeamId: item['winnerTeamId']?.toString() ?? '',
          myTeamId: item['myTeamId']?.toString() ?? '',
          myScore: (item['myScore'] as num?)?.toInt() ?? 0,
          opponentScore: (item['opponentScore'] as num?)?.toInt() ?? 0,
          finalScores: finalScores,
        ),
      );
    }

    history.sort((a, b) => b.playedAt.compareTo(a.playedAt));
    return history;
  }
}
