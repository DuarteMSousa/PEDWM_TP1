import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../domain/entities/lobby_realtime_event.dart';
import '../../domain/entities/room.dart';
import '../../domain/repositories/lobby_repository.dart';

class LobbyController extends ChangeNotifier {
  LobbyController({required LobbyRepository lobbyRepository})
    : _lobbyRepository = lobbyRepository;

  final LobbyRepository _lobbyRepository;

  bool isLoading = false;
  String? errorMessage;
  List<Room> rooms = const [];
  String? _connectedPlayerId;
  bool _isRefreshingFromEvent = false;
  bool _hasPendingEventRefresh = false;
  StreamSubscription<LobbyRealtimeEvent>? _eventsSubscription;

  Future<void> loadRooms({required String playerId}) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      await _ensureConnected(playerId: playerId);
      rooms = await _lobbyRepository.fetchRooms();
    } catch (error) {
      errorMessage = error.toString();
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<Room?> createRoom({required String hostPlayerId}) async {
    isLoading = true;
    errorMessage = null;
    notifyListeners();

    try {
      final createdRoom = await _lobbyRepository.createRoom(
        hostPlayerId: hostPlayerId,
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

  Future<void> _ensureConnected({required String playerId}) async {
    await _lobbyRepository.connectLobby(playerId: playerId);

    if (_eventsSubscription != null && _connectedPlayerId == playerId) {
      return;
    }

    await _eventsSubscription?.cancel();
    _connectedPlayerId = playerId;
    _eventsSubscription = _lobbyRepository.watchRealtimeEvents().listen(
      (event) {
        if (event.roomId != 'lobby') {
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
        rooms = await _lobbyRepository.fetchRooms();
        errorMessage = null;
        notifyListeners();
      } while (_hasPendingEventRefresh);
    } catch (_) {
      // Silent by design: temporary WS bursts should not spam user-facing errors.
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
