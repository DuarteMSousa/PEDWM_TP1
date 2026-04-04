class Profile {
  const Profile({
    required this.userId,
    required this.nickname,
    required this.matchesPlayed,
    required this.wins,
  });

  final String userId;
  final String nickname;
  final int matchesPlayed;
  final int wins;

  double get winRate {
    if (matchesPlayed == 0) {
      return 0;
    }
    return (wins / matchesPlayed) * 100;
  }
}
