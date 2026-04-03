import '../../domain/entities/match_replay.dart';
import '../../domain/entities/profile.dart';
import '../../domain/repositories/profile_repository.dart';
import '../datasources/profile_remote_data_source.dart';

class ProfileRepositoryImpl implements ProfileRepository {
  ProfileRepositoryImpl({required ProfileRemoteDataSource remoteDataSource})
    : _remoteDataSource = remoteDataSource;

  final ProfileRemoteDataSource _remoteDataSource;

  @override
  Future<Profile> fetchProfile(String userId) {
    return _remoteDataSource.fetchProfile(userId);
  }

  @override
  Future<MatchReplay?> fetchReplay(String gameId) {
    return _remoteDataSource.fetchReplay(gameId);
  }
}
