class GameSummary {
  const GameSummary({
    required this.id,
    required this.roomId,
    required this.players,
    required this.createdAt,
    required this.events,
  });

  final String id;
  final String? roomId;
  final List<GameSummaryPlayer> players;
  final DateTime createdAt;
  final List<GameEvent> events;
}

class GameSummaryPlayer {
  const GameSummaryPlayer({required this.id, required this.username});

  final String id;
  final String username;
}

class GameEvent {
  const GameEvent({
    required this.id,
    required this.type,
    required this.sequence,
    required this.timestamp,
    this.payload,
  });

  final String id;
  final String type;
  final int sequence;
  final DateTime timestamp;
  final Map<String, dynamic>? payload;
}
