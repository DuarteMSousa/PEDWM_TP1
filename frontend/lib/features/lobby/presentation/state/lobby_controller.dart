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

  Future<void> loadRooms({required String playerId}) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      await _lobbyRepository.connectLobby(playerId: playerId);
      rooms = await _lobbyRepository.fetchRooms();
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<Room?> createRoom({
    required String name,
    required String hostPlayerId,
    bool isPrivate = false,
  }) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final createdRoom = await _lobbyRepository.createRoom(
        name: name,
        hostPlayerId: hostPlayerId,
        isPrivate: isPrivate,
      );
      rooms = await _lobbyRepository.fetchRooms();
      return createdRoom;
    } catch (error) {
      errorMessage = error.toString();
      return null;
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<Room?> joinRoom({
    required String roomId,
    required String playerId,
  }) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final room = await _lobbyRepository.joinRoom(
        roomId: roomId,
        playerId: playerId,
      );
      rooms = await _lobbyRepository.fetchRooms();
      return room;
    } catch (error) {
      errorMessage = error.toString();
      return null;
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> refreshRooms({required String playerId}) =>
      loadRooms(playerId: playerId);

  @override
  void dispose() {
    unawaited(_lobbyRepository.disconnect());
    super.dispose();
  }
}
