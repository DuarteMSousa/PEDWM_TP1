class LobbyRealtimeEvent {
  const LobbyRealtimeEvent({
    required this.id,
    required this.type,
    required this.roomId,
    required this.timestamp,
    required this.payload,
  });

  final String id;
  final String type;
  final String roomId;
  final String timestamp;
  final Map<String, dynamic> payload;

  bool get isRoomDeleted => type == 'ROOM_DELETED';
  bool get isRoomUpdated =>
      type == 'ROOM_UPDATED' ||
      type == 'PLAYER_JOINED' ||
      type == 'PLAYER_LEFT';
}
