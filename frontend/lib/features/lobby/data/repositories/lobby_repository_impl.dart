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
  Future<Room> createRoom({required String hostPlayerId}) {
    return _remoteDataSource.createRoom(hostPlayerId: hostPlayerId);
  }

  @override
  Future<RoomDetails?> fetchRoomDetails({required String roomId}) {
    return _remoteDataSource.fetchRoomDetails(roomId: roomId);
  }

  @override
  Future<Room> joinRoom({required String roomId, required String playerId}) {
    return _remoteDataSource.joinRoom(roomId: roomId, playerId: playerId);
  }

  @override
  Future<RoomDetails> leaveRoom({
    required String roomId,
    required String playerId,
  }) {
    return _remoteDataSource.leaveRoom(roomId: roomId, playerId: playerId);
  }

  @override
  Future<RoomDetails> startGame({required String roomId}) {
    return _remoteDataSource.startGame(roomId: roomId);
  }

  @override
  Future<void> changeBotStrategy({required String strategy}) {
    return _remoteDataSource.changeBotStrategy(strategy: strategy);
  }
}
