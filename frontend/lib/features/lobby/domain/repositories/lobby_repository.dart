import '../entities/room_details.dart';
import '../entities/room.dart';

abstract class LobbyRepository {
  Future<void> connectLobby({required String playerId});
  Future<void> connectRoom({required String roomId, required String playerId});
  Future<void> disconnect();
  Future<List<Room>> fetchRooms();
  Future<Room> createRoom({
    required String name,
    required String hostPlayerId,
    int maxPlayers,
    bool isPrivate,
  });
  Future<RoomDetails?> fetchRoomDetails({required String roomId});
  Future<Room> joinRoom({required String roomId, required String playerId});
  Future<RoomDetails> leaveRoom({
    required String roomId,
    required String playerId,
  });
  Future<bool> deleteRoom({
    required String roomId,
    required String requesterId,
  });
}
