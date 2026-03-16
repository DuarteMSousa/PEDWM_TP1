import '../entities/room.dart';

abstract class LobbyRepository {
  Future<void> connect();
  Future<void> disconnect();
  Future<List<Room>> fetchRooms();
  Future<Room> joinRoom({required String roomId, required String playerId});
}
