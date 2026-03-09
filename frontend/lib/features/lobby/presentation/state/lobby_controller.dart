import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../domain/entities/room.dart';
import '../../domain/repositories/lobby_repository.dart';

class LobbyController extends ChangeNotifier {
  LobbyController({required LobbyRepository lobbyRepository})
    : _lobbyRepository = lobbyRepository;

  final LobbyRepository _lobbyRepository;

  bool isLoading = false;
  String? errorMessage;
  List<Room> rooms = const [];

  Future<void> loadRooms() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      await _lobbyRepository.connect();
      rooms = await _lobbyRepository.fetchRooms();
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> refreshRooms() => loadRooms();

  @override
  void dispose() {
    unawaited(_lobbyRepository.disconnect());
    super.dispose();
  }
}
