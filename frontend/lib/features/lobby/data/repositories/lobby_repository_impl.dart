import '../../domain/entities/room.dart';
import '../../domain/repositories/lobby_repository.dart';
import '../datasources/lobby_remote_data_source.dart';

class LobbyRepositoryImpl implements LobbyRepository {
  LobbyRepositoryImpl({required LobbyRemoteDataSource remoteDataSource})
    : _remoteDataSource = remoteDataSource;

  final LobbyRemoteDataSource _remoteDataSource;

  @override
  Future<void> connect() => _remoteDataSource.connect();

  @override
  Future<void> disconnect() => _remoteDataSource.disconnect();

  @override
  Future<List<Room>> fetchRooms() => _remoteDataSource.fetchRooms();

  @override
  Future<Room> joinRoom({required String roomId, required String playerId}) {
    return _remoteDataSource.joinRoom(roomId: roomId, playerId: playerId);
  }
}
