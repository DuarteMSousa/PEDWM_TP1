import 'package:flutter/foundation.dart';

import '../../domain/entities/profile.dart';
import '../../domain/repositories/profile_repository.dart';

class ProfileController extends ChangeNotifier {
  ProfileController({
    required ProfileRepository profileRepository,
    required this.userId,
  }) : _profileRepository = profileRepository;

  final ProfileRepository _profileRepository;
  final String userId;

  bool isLoading = false;
  String? errorMessage;
  Profile? profile;

  Future<void> loadProfile() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      profile = await _profileRepository.fetchProfile(userId);
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }
}
