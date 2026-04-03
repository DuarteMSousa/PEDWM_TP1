import 'package:flutter/foundation.dart';

import '../../domain/entities/match_replay.dart';
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
  bool isReplayLoading = false;
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

  Future<MatchReplay?> loadReplay(String gameId) async {
    isReplayLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      return await _profileRepository.fetchReplay(gameId);
    } catch (error) {
      errorMessage = error.toString();
      return null;
    } finally {
      isReplayLoading = false;
      notifyListeners();
    }
  }
}
