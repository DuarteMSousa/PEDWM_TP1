import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../domain/entities/game_summary.dart';

class ReplayPlayerController extends ChangeNotifier {
  ReplayPlayerController({required this.game});

  final GameSummary game;

  int _currentIndex = 0;
  bool _isPlaying = false;
  Timer? _timer;

  int get currentIndex => _currentIndex;
  bool get isPlaying => _isPlaying;
  int get totalEvents => game.events.length;
  bool get isAtEnd => _currentIndex >= totalEvents;
  bool get isAtStart => _currentIndex <= 0;

  List<GameEvent> get visibleEvents =>
      game.events.sublist(0, _currentIndex.clamp(0, totalEvents));

  GameEvent? get currentEvent =>
      _currentIndex > 0 && _currentIndex <= totalEvents
          ? game.events[_currentIndex - 1]
          : null;

  void play() {
    if (isAtEnd) return;
    _isPlaying = true;
    notifyListeners();
    _timer = Timer.periodic(const Duration(milliseconds: 800), (_) {
      if (_currentIndex >= totalEvents) {
        pause();
        return;
      }
      _currentIndex++;
      notifyListeners();
    });
  }

  void pause() {
    _isPlaying = false;
    _timer?.cancel();
    _timer = null;
    notifyListeners();
  }

  void stepForward() {
    if (_currentIndex < totalEvents) {
      _currentIndex++;
      notifyListeners();
    }
  }

  void stepBackward() {
    if (_currentIndex > 0) {
      _currentIndex--;
      notifyListeners();
    }
  }

  void seekTo(int index) {
    _currentIndex = index.clamp(0, totalEvents);
    notifyListeners();
  }

  void reset() {
    pause();
    _currentIndex = 0;
    notifyListeners();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }
}
