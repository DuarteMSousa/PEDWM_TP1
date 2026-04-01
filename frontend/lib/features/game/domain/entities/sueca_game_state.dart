import 'card.dart';
import 'game_phase.dart';
import 'player.dart';
import 'suit.dart';

class SuecaGameState {
  const SuecaGameState({
    required this.roomId,
    required this.phase,
    required this.players,
    required this.hand,
    required this.tableCards,
    required this.trumpSuit,
    required this.currentPlayerId,
    required this.myPlayerId,
    required this.teamAScore,
    required this.teamBScore,
    required this.roundTeamAScore,
    required this.roundTeamBScore,
  });

  final String roomId;
  final GamePhase phase;
  final List<Player> players;
  final List<SuecaCard> hand;
  final Map<String, SuecaCard> tableCards;
  final Suit trumpSuit;
  final String currentPlayerId;
  final String myPlayerId;
  final int teamAScore;
  final int teamBScore;
  final int roundTeamAScore;
  final int roundTeamBScore;

  SuecaGameState copyWith({
    GamePhase? phase,
    List<Player>? players,
    List<SuecaCard>? hand,
    Map<String, SuecaCard>? tableCards,
    Suit? trumpSuit,
    String? currentPlayerId,
    String? myPlayerId,
    int? teamAScore,
    int? teamBScore,
    int? roundTeamAScore,
    int? roundTeamBScore,
  }) {
    return SuecaGameState(
      roomId: roomId,
      phase: phase ?? this.phase,
      players: players ?? this.players,
      hand: hand ?? this.hand,
      tableCards: tableCards ?? this.tableCards,
      trumpSuit: trumpSuit ?? this.trumpSuit,
      currentPlayerId: currentPlayerId ?? this.currentPlayerId,
      myPlayerId: myPlayerId ?? this.myPlayerId,
      teamAScore: teamAScore ?? this.teamAScore,
      teamBScore: teamBScore ?? this.teamBScore,
      roundTeamAScore: roundTeamAScore ?? this.roundTeamAScore,
      roundTeamBScore: roundTeamBScore ?? this.roundTeamBScore,
    );
  }
}
