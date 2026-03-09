class Room {
  const Room({
    required this.id,
    required this.name,
    required this.playersCount,
    required this.maxPlayers,
    this.isPrivate = false,
  });

  final String id;
  final String name;
  final int playersCount;
  final int maxPlayers;
  final bool isPrivate;

  bool get isFull => playersCount >= maxPlayers;
  String get occupancyLabel => '$playersCount/$maxPlayers';
}
