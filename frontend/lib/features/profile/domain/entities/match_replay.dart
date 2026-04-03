class ReplayEvent {
  const ReplayEvent({
    required this.id,
    required this.type,
    required this.sequence,
    required this.timestamp,
    required this.payload,
  });

  final String id;
  final String type;
  final int sequence;
  final DateTime timestamp;
  final String payload;
}

class MatchReplay {
  const MatchReplay({
    required this.gameId,
    required this.roomId,
    required this.events,
  });

  final String gameId;
  final String roomId;
  final List<ReplayEvent> events;
}
