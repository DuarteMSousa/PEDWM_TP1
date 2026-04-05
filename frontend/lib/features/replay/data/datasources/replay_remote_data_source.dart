import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/game_summary.dart';

class ReplayRemoteDataSource {
  ReplayRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<List<GameSummary>> fetchUserGames(String userId) async {
    final response = await _graphqlService.query(
      document: '''
        query UserGames(\$userId: ID!) {
          userGames(userId: \$userId) {
            id
            roomId
            players {
              id
              username
            }
            createdAt
          }
        }
      ''',
      variables: <String, dynamic>{'userId': userId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for userGames query.');
    }

    final rawGames = data['userGames'];
    if (rawGames is! List) {
      return const [];
    }

    return rawGames
        .whereType<Map<String, dynamic>>()
        .map(_gameSummaryFromJson)
        .toList(growable: false);
  }

  Future<GameSummary> fetchGameReplay({
    required String userId,
    required String gameId,
  }) async {
    final response = await _graphqlService.query(
      document: '''
        query ReplayGame(\$userId: ID!) {
          userGames(userId: \$userId) {
            id
            roomId
            players {
              id
              username
            }
            events {
              id
              type
              sequence
              timestamp
              payload {
                ... on CardPlayedEventPayload {
                  playerId
                  card {
                    id
                    suit
                    rank
                  }
                }
                ... on TrickStartedEventPayload {
                  leaderId
                }
                ... on TurnChangedEventPayload {
                  playerId
                }
                ... on TrickEndedEventPayload {
                  winnerId
                  points
                }
                ... on RoundStartedEventPayload {
                  roundNumber
                  dealerId
                }
                ... on RoundEndedEventPayload {
                  winnerTeam
                  score {
                    teamId
                    points
                  }
                }
                ... on TrumpRevealedEventPayload {
                  suit
                  card {
                    id
                    suit
                    rank
                  }
                }
                ... on GameStartedEventPayload {
                  teams {
                    id
                    players {
                      id
                      name
                      sequence
                    }
                  }
                }
                ... on GameEndedEventPayload {
                  winner
                  finalScores {
                    teamId
                    points
                  }
                  teams {
                    id
                    players {
                      id
                      name
                      sequence
                    }
                  }
                }
                ... on GameScoreUpdatedEventPayload {
                  score {
                    teamId
                    points
                  }
                }
              }
            }
            createdAt
          }
        }
      ''',
      variables: <String, dynamic>{'userId': userId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for replay query.');
    }

    final rawGames = data['userGames'];
    if (rawGames is! List) {
      throw AppException('Replay nao encontrado.');
    }

    for (final rawGame in rawGames) {
      if (rawGame is! Map<String, dynamic>) {
        continue;
      }
      if (rawGame['id']?.toString() == gameId) {
        return _gameSummaryFromJson(rawGame);
      }
    }

    throw AppException('Replay nao encontrado.');
  }

  GameSummary _gameSummaryFromJson(Map<String, dynamic> json) {
    final players = <GameSummaryPlayer>[];
    final rawPlayers = json['players'];
    if (rawPlayers is List) {
      for (final p in rawPlayers) {
        if (p is Map<String, dynamic>) {
          players.add(GameSummaryPlayer(
            id: p['id']?.toString() ?? '',
            username: p['username']?.toString() ?? 'Jogador',
          ));
        }
      }
    }

    final events = <GameEvent>[];
    final rawEvents = json['events'];
    if (rawEvents is List) {
      for (final e in rawEvents) {
        if (e is Map<String, dynamic>) {
          events.add(_gameEventFromJson(e));
        }
      }
    }

    DateTime createdAt;
    try {
      createdAt = DateTime.parse(json['createdAt']?.toString() ?? '');
    } catch (_) {
      createdAt = DateTime.now();
    }

    return GameSummary(
      id: json['id']?.toString() ?? '',
      roomId: json['roomId']?.toString(),
      players: players,
      createdAt: createdAt,
      events: events
        ..sort((a, b) {
          final sequenceCompare = a.sequence.compareTo(b.sequence);
          if (sequenceCompare != 0) {
            return sequenceCompare;
          }
          return a.timestamp.compareTo(b.timestamp);
        }),
    );
  }

  GameEvent _gameEventFromJson(Map<String, dynamic> json) {
    DateTime timestamp;
    try {
      timestamp = DateTime.parse(json['timestamp']?.toString() ?? '');
    } catch (_) {
      timestamp = DateTime.now();
    }

    final rawPayload = json['payload'];
    Map<String, dynamic>? payload;
    if (rawPayload is Map<String, dynamic>) {
      payload = rawPayload;
    } else if (rawPayload is Map) {
      payload = Map<String, dynamic>.from(rawPayload);
    }

    return GameEvent(
      id: json['id']?.toString() ?? '',
      type: json['type']?.toString() ?? '',
      sequence: (json['sequence'] as num?)?.toInt() ?? 0,
      timestamp: timestamp,
      payload: payload,
    );
  }
}
