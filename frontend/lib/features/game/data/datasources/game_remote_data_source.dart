import 'dart:async';

import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../../../core/network/websocket/websocket_service.dart';
import '../../domain/entities/card.dart';
import '../../domain/entities/game_phase.dart';
import '../../domain/entities/player.dart';
import '../../domain/entities/sueca_game_state.dart';
import '../../domain/entities/suit.dart';

class GameRemoteDataSource {
  GameRemoteDataSource({
    required GraphqlService graphqlService,
    required WebSocketService webSocketService,
  }) : _graphqlService = graphqlService,
       _webSocketService = webSocketService;

  final GraphqlService _graphqlService;
  final WebSocketService _webSocketService;

  final Map<String, SuecaGameState> _cachedGames = <String, SuecaGameState>{};
  final Map<String, StreamController<SuecaGameState>> _controllers =
      <String, StreamController<SuecaGameState>>{};
  final Map<String, StreamSubscription<dynamic>> _eventSubscriptions =
      <String, StreamSubscription<dynamic>>{};

  Future<SuecaGameState> loadGame({
    required String roomId,
    required String playerId,
  }) async {
    await _webSocketService.connect(roomId: roomId, playerId: playerId);
    final roomFuture = _fetchRoom(roomId: roomId);
    final snapshotFuture = _fetchGameSnapshot(
      roomId: roomId,
      playerId: playerId,
    );
    final room = await roomFuture;
    final snapshot = await snapshotFuture;

    final players = _parsePlayers(room['players']);
    final currentPlayer = players.firstWhere(
      (player) => player.id == playerId,
      orElse: () => Player(id: playerId, nickname: 'Tu'),
    );
    if (players.every((player) => player.id != playerId)) {
      players.add(currentPlayer);
    }

    final state = _stateFromSnapshot(
      roomId: roomId,
      playerId: playerId,
      roomStatus: room['status']?.toString() ?? 'OPEN',
      players: players,
      snapshot: snapshot,
      fallbackCurrentPlayerId: currentPlayer.id,
    );

    _cachedGames[roomId] = state;
    _emit(roomId, state);
    _startRealtimeSync(roomId);
    return state;
  }

  Future<SuecaGameState> playCard({
    required String roomId,
    required SuecaCard card,
  }) async {
    final currentState = _cachedGames[roomId];
    if (currentState == null || !currentState.hand.contains(card)) {
      return currentState ?? _initialState(roomId: roomId, playerId: 'guest');
    }
    if (currentState.currentPlayerId != currentState.myPlayerId) {
      throw AppException('Ainda não é a tua vez.');
    }
    if (currentState.phase != GamePhase.playingTrick) {
      throw AppException('A ronda ainda não está pronta para jogar carta.');
    }

    await _webSocketService.send('play_card', <String, dynamic>{
      'cardId': card.backendId,
    });
    await _awaitCommandResult('play_card');

    // O estado é atualizado por eventos do backend (CARD_PLAYED, TURN_CHANGED, etc).
    return _cachedGames[roomId] ?? currentState;
  }

  Stream<SuecaGameState> watchGame(String roomId) {
    return _controllerFor(roomId).stream;
  }

  void _startRealtimeSync(String roomId) {
    if (_eventSubscriptions.containsKey(roomId)) {
      return;
    }

    _eventSubscriptions[roomId] = _webSocketService.events.listen(
      (rawEvent) => _onRealtimeEvent(roomId: roomId, rawEvent: rawEvent),
      onError: (Object error, StackTrace stackTrace) {
        final controller = _controllers[roomId];
        if (controller != null && !controller.isClosed) {
          controller.addError(error, stackTrace);
        }
      },
    );
  }

  void _onRealtimeEvent({
    required String roomId,
    required Map<String, dynamic> rawEvent,
  }) {
    final eventRoomId =
        rawEvent['roomId']?.toString() ?? rawEvent['gameId']?.toString();
    if (eventRoomId != roomId) {
      return;
    }

    final current = _cachedGames[roomId];
    if (current == null) {
      return;
    }

    final next = _applyEvent(current, rawEvent);
    _cachedGames[roomId] = next;
    _emit(roomId, next);
  }

  SuecaGameState _applyEvent(
    SuecaGameState current,
    Map<String, dynamic> rawEvent,
  ) {
    final type = rawEvent['type']?.toString() ?? '';
    final payload = rawEvent['payload'];
    final payloadMap = payload is Map<String, dynamic>
        ? payload
        : (payload is Map
              ? Map<String, dynamic>.from(payload)
              : <String, dynamic>{});

    switch (type) {
      case 'GAME_STARTED':
        return current.copyWith(phase: GamePhase.playingTrick);
      case 'ROUND_STARTED':
        return current.copyWith(
          phase: GamePhase.playingTrick,
          tableCards: <String, SuecaCard>{},
        );
      case 'TRICK_STARTED':
        final leaderId = payloadMap['leaderId']?.toString();
        return current.copyWith(
          phase: GamePhase.playingTrick,
          tableCards: <String, SuecaCard>{},
          currentPlayerId: leaderId ?? current.currentPlayerId,
        );
      case 'TURN_CHANGED':
        final playerId = payloadMap['playerId']?.toString();
        if (playerId == null || playerId.isEmpty) {
          return current;
        }
        return current.copyWith(
          phase: GamePhase.playingTrick,
          currentPlayerId: playerId,
        );
      case 'TRUMP_REVEALED':
        final trumpSuitRaw =
            payloadMap['suit']?.toString() ??
            (payloadMap['card'] is Map
                ? (payloadMap['card'] as Map)['suit']?.toString()
                : null);
        final trumpSuit = _parseSuit(trumpSuitRaw) ?? current.trumpSuit;
        return current.copyWith(
          trumpSuit: trumpSuit,
          phase: GamePhase.dealingCards,
        );
      case 'CARD_DEALT':
        final targetPlayerId = payloadMap['playerId']?.toString();
        final dealtCard = _parseCard(payloadMap['card']);
        if (targetPlayerId != current.myPlayerId || dealtCard == null) {
          return current;
        }
        final updatedHand = List<SuecaCard>.from(current.hand);
        if (!updatedHand.contains(dealtCard)) {
          updatedHand.add(dealtCard);
        }
        return current.copyWith(
          hand: updatedHand,
          phase: GamePhase.playingTrick,
        );
      case 'CARD_PLAYED':
        final playerId = payloadMap['playerId']?.toString();
        final playedCard = _parseCard(payloadMap['card']);
        if (playerId == null || playedCard == null) {
          return current;
        }

        final updatedTableCards = Map<String, SuecaCard>.from(
          current.tableCards,
        )..[playerId] = playedCard;

        final updatedHand = List<SuecaCard>.from(current.hand);
        if (playerId == current.myPlayerId) {
          updatedHand.remove(playedCard);
        }

        return current.copyWith(
          tableCards: updatedTableCards,
          hand: updatedHand,
          phase: GamePhase.playingTrick,
        );
      case 'TRICK_ENDED':
        return current.copyWith(phase: GamePhase.scoring);
      case 'ROUND_ENDED':
        final scoreTuple = _parseScoreTuple(payloadMap['score']);
        return current.copyWith(
          phase: GamePhase.scoring,
          teamAScore: scoreTuple.$1,
          teamBScore: scoreTuple.$2,
          tableCards: <String, SuecaCard>{},
        );
      case 'GAME_ENDED':
        final scoreTuple = _parseScoreTuple(payloadMap['finalScores']);
        return current.copyWith(
          phase: GamePhase.finished,
          teamAScore: scoreTuple.$1,
          teamBScore: scoreTuple.$2,
        );
      default:
        return current;
    }
  }

  Future<Map<String, dynamic>> _fetchRoom({required String roomId}) async {
    final response = await _graphqlService.query(
      document: '''
        query Room(\$id: ID!) {
          room(id: \$id) {
            id
            status
            players {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{'id': roomId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw Exception('Invalid GraphQL response for room query.');
    }

    final room = data['room'];
    if (room is! Map<String, dynamic>) {
      throw Exception('Room not found or invalid GraphQL payload.');
    }

    return room;
  }

  Future<Map<String, dynamic>> _fetchGameSnapshot({
    required String roomId,
    required String playerId,
  }) async {
    final response = await _graphqlService.query(
      document: '''
        query GameSnapshot(\$roomId: ID!, \$playerId: ID!) {
          gameSnapshot(roomId: \$roomId, playerId: \$playerId) {
            roomId
            gameId
            status
            trumpSuit
            currentPlayerId
            myHand {
              id
              suit
              rank
            }
            tablePlays {
              playerId
              card {
                id
                suit
                rank
              }
            }
            scores {
              teamId
              points
            }
          }
        }
      ''',
      variables: <String, dynamic>{'roomId': roomId, 'playerId': playerId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for gameSnapshot.');
    }

    final snapshot = data['gameSnapshot'];
    if (snapshot is! Map<String, dynamic>) {
      throw AppException('gameSnapshot not found.');
    }

    return snapshot;
  }

  SuecaGameState _stateFromSnapshot({
    required String roomId,
    required String playerId,
    required String roomStatus,
    required List<Player> players,
    required Map<String, dynamic> snapshot,
    required String fallbackCurrentPlayerId,
  }) {
    final myHandRaw = snapshot['myHand'];
    final tablePlaysRaw = snapshot['tablePlays'];
    final scoresRaw = snapshot['scores'];
    final trumpSuitRaw = snapshot['trumpSuit']?.toString();
    final currentPlayerRaw = snapshot['currentPlayerId']?.toString();
    final statusRaw = snapshot['status']?.toString();

    final hand = <SuecaCard>[];
    if (myHandRaw is List) {
      for (final item in myHandRaw) {
        final parsed = _parseCard(item);
        if (parsed != null) {
          hand.add(parsed);
        }
      }
    }

    final tableCards = <String, SuecaCard>{};
    if (tablePlaysRaw is List) {
      for (final item in tablePlaysRaw) {
        if (item is! Map) {
          continue;
        }
        final mapped = Map<String, dynamic>.from(item);
        final pId = mapped['playerId']?.toString() ?? '';
        final parsedCard = _parseCard(mapped['card']);
        if (pId.isEmpty || parsedCard == null) {
          continue;
        }
        tableCards[pId] = parsedCard;
      }
    }

    final scoreTuple = _parseSnapshotScores(scoresRaw);
    final trumpSuit = _parseSuit(trumpSuitRaw) ?? Suit.hearts;
    final currentPlayerId =
        (currentPlayerRaw == null || currentPlayerRaw.isEmpty)
        ? fallbackCurrentPlayerId
        : currentPlayerRaw;
    final phase = _phaseFromStatus(
      roomStatus: roomStatus,
      gameStatus: statusRaw,
      hasCards: hand.isNotEmpty || tableCards.isNotEmpty,
    );

    return SuecaGameState(
      roomId: roomId,
      phase: phase,
      players: players,
      hand: hand,
      tableCards: tableCards,
      trumpSuit: trumpSuit,
      currentPlayerId: currentPlayerId,
      myPlayerId: playerId,
      teamAScore: scoreTuple.$1,
      teamBScore: scoreTuple.$2,
    );
  }

  List<Player> _parsePlayers(dynamic rawPlayers) {
    if (rawPlayers is! List) {
      return <Player>[];
    }

    return rawPlayers
        .whereType<Map<String, dynamic>>()
        .map(
          (rawPlayer) => Player(
            id: rawPlayer['id']?.toString() ?? '',
            nickname: rawPlayer['username']?.toString() ?? 'Jogador',
          ),
        )
        .where((player) => player.id.isNotEmpty)
        .toList(growable: true);
  }

  GamePhase _phaseFromRoomStatus(String status) {
    switch (status) {
      case 'OPEN':
        return GamePhase.waitingForPlayers;
      case 'IN_GAME':
        return GamePhase.playingTrick;
      case 'CLOSED':
        return GamePhase.finished;
      default:
        return GamePhase.waitingForPlayers;
    }
  }

  GamePhase _phaseFromStatus({
    required String roomStatus,
    required String? gameStatus,
    required bool hasCards,
  }) {
    final normalizedGameStatus = gameStatus?.trim().toUpperCase() ?? '';
    if (normalizedGameStatus == 'FINISHED') {
      return GamePhase.finished;
    }
    if (normalizedGameStatus == 'IN_PROGRESS' || roomStatus == 'IN_GAME') {
      return hasCards ? GamePhase.playingTrick : GamePhase.dealingCards;
    }
    return _phaseFromRoomStatus(roomStatus);
  }

  (int, int) _parseScoreTuple(dynamic rawScore) {
    if (rawScore is! Map) {
      return (0, 0);
    }

    final mapped = Map<String, dynamic>.from(rawScore);
    final teamA = _toInt(mapped['team1']);
    final teamB = _toInt(mapped['team2']);
    if (teamA != null || teamB != null) {
      return (teamA ?? 0, teamB ?? 0);
    }

    final numericScores = mapped.values
        .map(_toInt)
        .whereType<int>()
        .toList(growable: false);
    if (numericScores.length >= 2) {
      return (numericScores[0], numericScores[1]);
    }
    if (numericScores.length == 1) {
      return (numericScores[0], 0);
    }
    return (0, 0);
  }

  (int, int) _parseSnapshotScores(dynamic rawScores) {
    if (rawScores is! List) {
      return (0, 0);
    }

    int? team1;
    int? team2;
    final collected = <int>[];

    for (final item in rawScores) {
      if (item is! Map) {
        continue;
      }
      final mapped = Map<String, dynamic>.from(item);
      final points = _toInt(mapped['points']);
      if (points == null) {
        continue;
      }
      collected.add(points);

      final teamId = mapped['teamId']?.toString().toLowerCase() ?? '';
      if (teamId.contains('1') && team1 == null) {
        team1 = points;
      } else if (teamId.contains('2') && team2 == null) {
        team2 = points;
      }
    }

    if (team1 != null || team2 != null) {
      return (team1 ?? 0, team2 ?? 0);
    }
    if (collected.length >= 2) {
      return (collected[0], collected[1]);
    }
    if (collected.length == 1) {
      return (collected[0], 0);
    }
    return (0, 0);
  }

  Future<void> _awaitCommandResult(String commandType) async {
    final Map<String, dynamic> response = await _webSocketService.events
        .firstWhere(
          (event) =>
              event['type']?.toString() == commandType &&
              event.containsKey('success'),
        )
        .timeout(
          const Duration(seconds: 4),
          onTimeout: () => throw AppException(
            'Sem confirmação do servidor para $commandType.',
          ),
        );

    final success = response['success'] == true;
    if (success) {
      return;
    }

    final error = response['error']?.toString();
    throw AppException(
      error == null || error.isEmpty ? 'A jogada falhou.' : error,
    );
  }

  SuecaCard? _parseCard(dynamic rawCard) {
    if (rawCard is! Map) {
      return null;
    }
    final map = Map<String, dynamic>.from(rawCard);
    final suit = _parseSuit(map['suit']?.toString());
    final rank = _parseRank(map['rank']);
    if (suit == null || rank == null) {
      return null;
    }
    return SuecaCard(suit: suit, rank: rank);
  }

  Suit? _parseSuit(String? rawSuit) {
    if (rawSuit == null || rawSuit.trim().isEmpty) {
      return null;
    }

    switch (rawSuit.trim().toUpperCase()) {
      case 'HEARTS':
        return Suit.hearts;
      case 'SPADES':
        return Suit.spades;
      case 'DIAMONDS':
        return Suit.diamonds;
      case 'CLUBS':
        return Suit.clubs;
      default:
        return null;
    }
  }

  int? _parseRank(dynamic rawRank) {
    if (rawRank is num) {
      return rawRank.toInt();
    }
    if (rawRank == null) {
      return null;
    }

    final token = rawRank.toString().trim().toUpperCase();
    switch (token) {
      case 'A':
        return 1;
      case 'K':
        return 13;
      case 'Q':
        return 12;
      case 'J':
        return 11;
      default:
        return int.tryParse(token);
    }
  }

  int? _toInt(dynamic value) {
    if (value is num) {
      return value.toInt();
    }
    if (value == null) {
      return null;
    }
    return int.tryParse(value.toString());
  }

  void _emit(String roomId, SuecaGameState state) {
    final controller = _controllerFor(roomId);
    if (!controller.isClosed) {
      controller.add(state);
    }
  }

  StreamController<SuecaGameState> _controllerFor(String roomId) {
    return _controllers.putIfAbsent(
      roomId,
      () => StreamController<SuecaGameState>.broadcast(),
    );
  }

  SuecaGameState _initialState({
    required String roomId,
    required String playerId,
  }) {
    return SuecaGameState(
      roomId: roomId,
      phase: GamePhase.waitingForPlayers,
      players: <Player>[Player(id: playerId, nickname: 'Tu')],
      hand: const <SuecaCard>[],
      tableCards: const <String, SuecaCard>{},
      trumpSuit: Suit.hearts,
      currentPlayerId: playerId,
      myPlayerId: playerId,
      teamAScore: 0,
      teamBScore: 0,
    );
  }

  void dispose() {
    for (final subscription in _eventSubscriptions.values) {
      subscription.cancel();
    }
    _eventSubscriptions.clear();

    for (final controller in _controllers.values) {
      controller.close();
    }
    _controllers.clear();
  }
}
