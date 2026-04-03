class MatchHistoryEntry {
  const MatchHistoryEntry({
    required this.gameId,
    required this.roomId,
    required this.playedAt,
    required this.won,
    required this.winnerTeamId,
    required this.myTeamId,
    required this.myScore,
    required this.opponentScore,
    required this.finalScores,
  });

  final String gameId;
  final String roomId;
  final DateTime playedAt;
  final bool won;
  final String winnerTeamId;
  final String myTeamId;
  final int myScore;
  final int opponentScore;
  final Map<String, int> finalScores;
}
