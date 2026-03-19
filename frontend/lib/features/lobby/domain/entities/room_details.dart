import 'room_member.dart';

class RoomDetails {
  const RoomDetails({
    required this.id,
    required this.name,
    required this.hostPlayerId,
    required this.status,
    required this.maxPlayers,
    required this.isPrivate,
    required this.players,
  });

  final String id;
  final String name;
  final String hostPlayerId;
  final String status;
  final int maxPlayers;
  final bool isPrivate;
  final List<RoomMember> players;

  int get playersCount => players.length;
  bool get isFull => playersCount >= maxPlayers;
  bool get hasRequiredPlayers => playersCount >= 4;
}
