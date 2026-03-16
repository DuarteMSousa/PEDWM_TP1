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

  Future<SuecaGameState> loadGame(String roomId) async {
    await _webSocketService.connect(roomId: roomId, playerId: 'p1');
    await _graphqlService.query(
      document:
          'query Room(\$id: ID!) { room(id: \$id) { id phase trumpSuit } }',
      variables: <String, dynamic>{'id': roomId},
    );

    final state = _cachedGames.putIfAbsent(roomId, () => _initialState(roomId));
    _emit(roomId, state);
    return state;
  }

  Future<SuecaGameState> playCard({
    required String roomId,
    required SuecaCard card,
  }) async {
    final currentState = _cachedGames.putIfAbsent(
      roomId,
      () => _initialState(roomId),
    );

    if (!currentState.hand.contains(card)) {
      return currentState;
    }

    final updatedHand = List<SuecaCard>.from(currentState.hand)..remove(card);
    final updatedTableCards = Map<String, SuecaCard>.from(
      currentState.tableCards,
    )..[currentState.myPlayerId] = card;

    final nextPhase = updatedHand.isEmpty
        ? GamePhase.finished
        : GamePhase.playingTrick;
    final nextState = currentState.copyWith(
      hand: updatedHand,
      tableCards: updatedTableCards,
      phase: nextPhase,
      teamAScore: currentState.teamAScore + card.points,
    );

    _cachedGames[roomId] = nextState;
    _emit(roomId, nextState);

    await _webSocketService.send('game.play_card', <String, dynamic>{
      'roomId': roomId,
      'card': <String, dynamic>{'suit': card.suit.name, 'rank': card.rank},
    });

    return nextState;
  }

  Stream<SuecaGameState> watchGame(String roomId) {
    return _controllerFor(roomId).stream;
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

  SuecaGameState _initialState(String roomId) {
    const players = [
      Player(id: 'p1', nickname: 'You'),
      Player(id: 'p2', nickname: 'Mate'),
      Player(id: 'p3', nickname: 'Opp A'),
      Player(id: 'p4', nickname: 'Opp B'),
    ];

    const hand = [
      SuecaCard(suit: Suit.hearts, rank: 1),
      SuecaCard(suit: Suit.spades, rank: 7),
      SuecaCard(suit: Suit.clubs, rank: 13),
      SuecaCard(suit: Suit.diamonds, rank: 11),
      SuecaCard(suit: Suit.hearts, rank: 3),
    ];

    return SuecaGameState(
      roomId: roomId,
      phase: GamePhase.playingTrick,
      players: players,
      hand: hand,
      tableCards: <String, SuecaCard>{},
      trumpSuit: Suit.hearts,
      currentPlayerId: 'p1',
      myPlayerId: 'p1',
      teamAScore: 0,
      teamBScore: 0,
    );
  }

  void dispose() {
    for (final controller in _controllers.values) {
      controller.close();
    }
    _controllers.clear();
  }
}
