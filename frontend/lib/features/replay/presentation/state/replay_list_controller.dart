import 'package:flutter/foundation.dart';

import '../../domain/entities/game_summary.dart';
import '../../domain/repositories/replay_repository.dart';

class ReplayListController extends ChangeNotifier {
  ReplayListController({
    required ReplayRepository replayRepository,
    required this.userId,
  }) : _replayRepository = replayRepository;

  final ReplayRepository _replayRepository;
  final String userId;

  bool isLoading = false;
  String? errorMessage;
  List<GameSummary> games = const [];

  Future<void> loadGames() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      games = await _replayRepository.fetchUserGames(userId);
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }
}
