import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../domain/entities/lobby_realtime_event.dart';
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
  StreamSubscription<LobbyRealtimeEvent>? _eventsSubscription;

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
      await _startRealtimeSync();
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
        errorMessage = null;
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

  Future<void> _startRealtimeSync() async {
    await _eventsSubscription?.cancel();
    _eventsSubscription = _lobbyRepository.watchRealtimeEvents().listen(
      (event) {
        final sameRoom = event.roomId == roomId;
        final roomDeleted =
            event.isRoomDeleted &&
            event.payload['roomId']?.toString() == roomId;
        if (!sameRoom && !roomDeleted) {
          return;
        }
        unawaited(refreshRoom());
      },
      onError: (Object error) {
        errorMessage = error.toString();
        notifyListeners();
      },
    );
  }

  @override
  void dispose() {
    _eventsSubscription?.cancel();
    unawaited(_lobbyRepository.disconnect());
    super.dispose();
  }
}
