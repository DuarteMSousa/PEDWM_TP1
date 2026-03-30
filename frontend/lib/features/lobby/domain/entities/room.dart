class Room {
  const Room({
    required this.id,
    required this.name,
    required this.playersCount,
    required this.maxPlayers,
    required this.status,
    this.isPrivate = false,
  });

  final String id;
  final String name;
  final int playersCount;
  final int maxPlayers;
  final String status;
  final bool isPrivate;

  bool get isFull => playersCount >= maxPlayers;
  String get occupancyLabel => '$playersCount/$maxPlayers';
}
