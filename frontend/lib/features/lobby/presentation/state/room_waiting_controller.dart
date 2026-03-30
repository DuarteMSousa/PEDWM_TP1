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
  bool _isRefreshingFromEvent = false;
  bool _hasPendingEventRefresh = false;

  bool get hasAllPlayers => room?.hasRequiredPlayers == true;
  bool get isHost => room?.hostPlayerId == currentPlayerId;
  bool get hasGameStarted => room?.status == 'IN_GAME';

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

  Future<bool> startGame() async {
    if (!isHost) {
      errorMessage = 'Apenas o host pode iniciar o jogo.';
      notifyListeners();
      return false;
    }
    if (!hasAllPlayers) {
      errorMessage = 'Sao precisos 4 jogadores para iniciar.';
      notifyListeners();
      return false;
    }

    isActionLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      room = await _lobbyRepository.startGame(roomId: roomId);
      roomUnavailable = false;
      return true;
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
        if (event.roomId != roomId) {
          return;
        }
        unawaited(_refreshFromRealtimeEvent());
      },
      onError: (Object error) {
        errorMessage = error.toString();
        notifyListeners();
      },
    );
  }

  Future<void> _refreshFromRealtimeEvent() async {
    if (_isRefreshingFromEvent) {
      _hasPendingEventRefresh = true;
      return;
    }

    _isRefreshingFromEvent = true;
    try {
      do {
        _hasPendingEventRefresh = false;
        await refreshRoom();
      } while (_hasPendingEventRefresh);
    } finally {
      _isRefreshingFromEvent = false;
    }
  }

  @override
  void dispose() {
    _eventsSubscription?.cancel();
    super.dispose();
  }
}
