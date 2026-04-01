import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_data_source.dart';

class AuthRepositoryImpl implements AuthRepository {
  AuthRepositoryImpl({required AuthRemoteDataSource remoteDataSource})
    : _remoteDataSource = remoteDataSource;

  final AuthRemoteDataSource _remoteDataSource;

  @override
  Future<User> login({required String username, required String password}) {
    return _remoteDataSource.login(username: username, password: password);
  }

  @override
  Future<User> register({required String username, required String password}) {
    return _remoteDataSource.register(username: username, password: password);
  }
}
