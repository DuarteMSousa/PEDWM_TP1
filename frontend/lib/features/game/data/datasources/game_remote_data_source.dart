import 'dart:async';

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
    final room = await _fetchRoom(roomId: roomId);

    final players = _parsePlayers(room['players']);
    final currentPlayer = players.firstWhere(
      (player) => player.id == playerId,
      orElse: () => Player(id: playerId, nickname: 'Tu'),
    );
    if (players.every((player) => player.id != playerId)) {
      players.add(currentPlayer);
    }

    final existing = _cachedGames[roomId];
    final roomStatus = room['status']?.toString() ?? 'OPEN';
    final state = SuecaGameState(
      roomId: roomId,
      phase: _phaseFromRoomStatus(roomStatus),
      players: players,
      hand: existing?.hand ?? const <SuecaCard>[],
      tableCards: existing?.tableCards ?? const <String, SuecaCard>{},
      trumpSuit: existing?.trumpSuit ?? Suit.hearts,
      currentPlayerId: existing?.currentPlayerId ?? currentPlayer.id,
      myPlayerId: playerId,
      teamAScore: existing?.teamAScore ?? 0,
      teamBScore: existing?.teamBScore ?? 0,
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

    final updatedHand = List<SuecaCard>.from(currentState.hand)..remove(card);
    final updatedTableCards = Map<String, SuecaCard>.from(
      currentState.tableCards,
    )..[currentState.myPlayerId] = card;

    final nextState = currentState.copyWith(
      hand: updatedHand,
      tableCards: updatedTableCards,
      phase: updatedHand.isEmpty ? GamePhase.finished : GamePhase.playingTrick,
    );

    _cachedGames[roomId] = nextState;
    _emit(roomId, nextState);
    return nextState;
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
