import '../../domain/entities/lobby_realtime_event.dart';
import '../../domain/entities/room_details.dart';
import '../../domain/entities/room.dart';
import '../../domain/repositories/lobby_repository.dart';
import '../datasources/lobby_remote_data_source.dart';

class LobbyRepositoryImpl implements LobbyRepository {
  LobbyRepositoryImpl({required LobbyRemoteDataSource remoteDataSource})
    : _remoteDataSource = remoteDataSource;

  final LobbyRemoteDataSource _remoteDataSource;

  @override
  Future<void> connectLobby({required String playerId}) {
    return _remoteDataSource.connectLobby(playerId: playerId);
  }

  @override
  Future<void> connectRoom({required String roomId, required String playerId}) {
    return _remoteDataSource.connectRoom(roomId: roomId, playerId: playerId);
  }

  @override
  Stream<LobbyRealtimeEvent> watchRealtimeEvents() {
    return _remoteDataSource.watchRealtimeEvents();
  }

  @override
  Future<void> disconnect() => _remoteDataSource.disconnect();

  @override
  Future<List<Room>> fetchRooms() => _remoteDataSource.fetchRooms();

  @override
  Future<Room> createRoom({
    required String name,
    required String hostPlayerId,
    int maxPlayers = 4,
    bool isPrivate = false,
    String? password,
  }) {
    return _remoteDataSource.createRoom(
      name: name,
      hostPlayerId: hostPlayerId,
      maxPlayers: maxPlayers,
      isPrivate: isPrivate,
      password: password,
    );
  }

  @override
  Future<RoomDetails?> fetchRoomDetails({required String roomId}) {
    return _remoteDataSource.fetchRoomDetails(roomId: roomId);
  }

  @override
  Future<Room> joinRoom({
    required String roomId,
    required String playerId,
    String? password,
  }) {
    return _remoteDataSource.joinRoom(
      roomId: roomId,
      playerId: playerId,
      password: password,
    );
  }

  @override
  Future<RoomDetails> leaveRoom({
    required String roomId,
    required String playerId,
  }) {
    return _remoteDataSource.leaveRoom(roomId: roomId, playerId: playerId);
  }

  @override
  Future<bool> deleteRoom({
    required String roomId,
    required String requesterId,
  }) {
    return _remoteDataSource.deleteRoom(
      roomId: roomId,
      requesterId: requesterId,
    );
  }
}
