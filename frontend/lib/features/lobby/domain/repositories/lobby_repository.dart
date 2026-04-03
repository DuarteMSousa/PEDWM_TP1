import '../entities/lobby_realtime_event.dart';
import '../entities/room_details.dart';
import '../entities/room.dart';

abstract class LobbyRepository {
  Future<void> connectLobby({required String playerId});
  Future<void> connectRoom({required String roomId, required String playerId});
  Stream<LobbyRealtimeEvent> watchRealtimeEvents();
  Future<void> disconnect();
  Future<List<Room>> fetchRooms();
  Future<Room> createRoom({required String hostPlayerId});
  Future<RoomDetails?> fetchRoomDetails({required String roomId});
  Future<Room> joinRoom({required String roomId, required String playerId});
  Future<RoomDetails> leaveRoom({
    required String roomId,
    required String playerId,
  });
  Future<RoomDetails> startGame({
    required String roomId,
    required String requesterId,
  });
}
