import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../domain/entities/room_details.dart';
import '../../domain/repositories/lobby_repository.dart';

class RoomWaitingController extends ChangeNotifier {
  RoomWaitingController({
    required LobbyRepository lobbyRepository,
    required this.roomId,
    required this.currentPlayerId,
  }) : _lobbyRepository = lobbyRepository;

  final LobbyRepository _lobbyRepository;
  final String roomId;
  final String currentPlayerId;

  bool isLoading = true;
  bool isActionLoading = false;
  String? errorMessage;
  RoomDetails? room;
  bool roomUnavailable = false;

  Timer? _pollingTimer;

  bool get hasAllPlayers => room?.hasRequiredPlayers == true;
  bool get isHost => room?.hostPlayerId == currentPlayerId;

  Future<void> initialize() async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      await _lobbyRepository.connectRoom(
        roomId: roomId,
        playerId: currentPlayerId,
      );
      await refreshRoom();
      _startPolling();
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> refreshRoom() async {
    try {
      final updatedRoom = await _lobbyRepository.fetchRoomDetails(
        roomId: roomId,
      );
      if (updatedRoom == null) {
        roomUnavailable = true;
        errorMessage = 'A sala foi removida.';
      } else {
        room = updatedRoom;
        roomUnavailable = false;
      }
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      notifyListeners();
    }
  }

  Future<bool> leaveRoom() async {
    isActionLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      await _lobbyRepository.leaveRoom(
        roomId: roomId,
        playerId: currentPlayerId,
      );
      return true;
    } catch (error) {
      errorMessage = error.toString();
      return false;
    } finally {
      isActionLoading = false;
      notifyListeners();
    }
  }

  Future<bool> deleteRoom() async {
    isActionLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final deleted = await _lobbyRepository.deleteRoom(
        roomId: roomId,
        requesterId: currentPlayerId,
      );
      if (!deleted) {
        errorMessage = 'Nao foi possivel eliminar a sala.';
      }
      return deleted;
    } catch (error) {
      errorMessage = error.toString();
      return false;
    } finally {
      isActionLoading = false;
      notifyListeners();
    }
  }

  void _startPolling() {
    _pollingTimer?.cancel();
    _pollingTimer = Timer.periodic(const Duration(seconds: 2), (_) {
      unawaited(refreshRoom());
    });
  }

  @override
  void dispose() {
    _pollingTimer?.cancel();
    unawaited(_lobbyRepository.disconnect());
    super.dispose();
  }
}
