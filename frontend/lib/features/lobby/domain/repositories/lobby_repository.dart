import '../entities/room.dart';

abstract class LobbyRepository {
  Future<void> connect();
  Future<void> disconnect();
  Future<List<Room>> fetchRooms();
}
