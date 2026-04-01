import '../entities/user.dart';

abstract class AuthRepository {
  Future<User> login({required String username, required String password});
  Future<User> register({required String username, required String password});
}
