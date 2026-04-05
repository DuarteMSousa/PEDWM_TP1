import '../../domain/entities/game_summary.dart';
import '../../domain/repositories/replay_repository.dart';
import '../datasources/replay_remote_data_source.dart';

class ReplayRepositoryImpl implements ReplayRepository {
  ReplayRepositoryImpl({required ReplayRemoteDataSource remoteDataSource})
    : _remoteDataSource = remoteDataSource;

  final ReplayRemoteDataSource _remoteDataSource;

  @override
  Future<List<GameSummary>> fetchUserGames(String userId) {
    return _remoteDataSource.fetchUserGames(userId);
  }
}
