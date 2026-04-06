import 'package:sueca_pedwm/features/game/domain/entities/team.dart';

import 'card.dart';
import 'game_phase.dart';
import 'player.dart';
import 'suit.dart';

class SuecaGameState {
  const SuecaGameState({
    required this.roomId,
    required this.phase,
    required this.players,
    required this.teams,
    required this.hand,
    required this.tableCards,
    required this.trumpSuit,
    required this.currentPlayerId,
    required this.myPlayerId,
  });

  final String roomId;
  final GamePhase phase;
  final List<Player> players;
  final List<Team> teams;
  final List<SuecaCard> hand;
  final Map<String, SuecaCard> tableCards;
  final Suit trumpSuit;
  final String currentPlayerId;
  final String myPlayerId;

  SuecaGameState copyWith({
    GamePhase? phase,
    List<Player>? players,
    List<Team>? teams,
    List<SuecaCard>? hand,
    Map<String, SuecaCard>? tableCards,
    Suit? trumpSuit,
    String? currentPlayerId,
    String? myPlayerId
  }) {

    final sortedHand = List<SuecaCard>.from(hand ?? this.hand);


    sortedHand.sort((a, b) {
      if (a.suit != b.suit) {
        return a.suit.index.compareTo(b.suit.index);
      }
      return (b.rank).compareTo(a.rank); 
    });

    return SuecaGameState(
      roomId: roomId,
      phase: phase ?? this.phase,
      players: players ?? this.players,
      teams: teams ?? this.teams,
      hand: sortedHand,
      tableCards: tableCards ?? this.tableCards,
      trumpSuit: trumpSuit ?? this.trumpSuit,
      currentPlayerId: currentPlayerId ?? this.currentPlayerId,
      myPlayerId: myPlayerId ?? this.myPlayerId,
    );
  }
}
