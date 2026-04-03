class Player {
  const Player({required this.id, required this.nickname, this.sequence = 1});

  final String id;
  final String nickname;
  final int sequence;

  factory Player.fromJson(Map<String, dynamic> json) {
    return Player(
      id: json['id']?.toString() ?? '',
      nickname: json['name']?.toString() ?? 'Jogador',
      sequence: json['sequence']?.toInt() ?? 1,
    );
  }
}
