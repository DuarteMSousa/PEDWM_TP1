import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../domain/entities/card.dart';
import '../../domain/entities/game_phase.dart';
import '../../domain/entities/sueca_game_state.dart';
import '../../domain/repositories/game_repository.dart';

class GameController extends ChangeNotifier {
  GameController({
    required GameRepository gameRepository,
    required this.roomId,
    required this.currentPlayerId,
  }) : _gameRepository = gameRepository;

  final GameRepository _gameRepository;
  final String roomId;
  final String currentPlayerId;

  bool isLoading = false;
  String? errorMessage;
  SuecaGameState? gameState;
  StreamSubscription<SuecaGameState>? _subscription;

  bool get canPlayCard {
    final state = gameState;
    if (state == null) {
      return false;
    }
    return state.phase == GamePhase.playingTrick &&
        state.currentPlayerId == state.myPlayerId &&
        !isLoading;
  }

  Future<void> initialize() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      gameState = await _gameRepository.loadGame(
        roomId: roomId,
        playerId: currentPlayerId,
      );
      _subscription = _gameRepository
          .watchGame(roomId)
          .listen(
            (nextState) {
              gameState = nextState;
              notifyListeners();
            },
            onError: (Object error) {
              errorMessage = error.toString();
              notifyListeners();
            },
          );
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> playCard(SuecaCard card) async {
    if (!canPlayCard) {
      errorMessage = 'Ainda não é a tua vez.';
      notifyListeners();
      return;
    }

    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      gameState = await _gameRepository.playCard(roomId: roomId, card: card);
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
  }
}
