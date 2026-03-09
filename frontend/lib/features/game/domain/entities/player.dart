class Player {
  const Player({required this.id, required this.nickname, this.tricksWon = 0});

  final String id;
  final String nickname;
  final int tricksWon;
}
