import 'package:flutter/foundation.dart';

import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';

class AuthController extends ChangeNotifier {
  AuthController({required AuthRepository authRepository})
    : _authRepository = authRepository;

  final AuthRepository _authRepository;

  bool isLoading = false;
  String? errorMessage;
  User? currentUser;

  Future<void> enterWithNickname(String nickname) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      currentUser = await _authRepository.enterWithNickname(nickname);
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }
}
