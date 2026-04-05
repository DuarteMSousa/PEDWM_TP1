class Profile {
  const Profile({
    required this.userId,
    required this.nickname,
    required this.matchesPlayed,
    required this.wins,
    required this.elo,
  });

  final String userId;
  final String nickname;
  final int matchesPlayed;
  final int wins;
  final int elo;

  double get winRate {
    if (matchesPlayed == 0) {
      return 0;
    }
    return (wins / matchesPlayed) * 100;
  }
}
