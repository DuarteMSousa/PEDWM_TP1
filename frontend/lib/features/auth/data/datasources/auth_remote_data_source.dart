import '../../../../core/error/app_exception.dart';
import '../../domain/entities/user.dart';

class AuthRemoteDataSource {
  AuthRemoteDataSource();

  Future<User> enterWithNickname(String nickname) async {
    final cleanedNickname = nickname.trim();
    if (cleanedNickname.isEmpty) {
      throw AppException('Nickname is required.');
    }

    return User(
      id: cleanedNickname.toLowerCase().replaceAll(' ', '_'),
      nickname: cleanedNickname,
    );
  }
}
