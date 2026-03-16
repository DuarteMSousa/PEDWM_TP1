import 'dart:async';
import 'dart:convert';

import 'package:web_socket_channel/web_socket_channel.dart';

import '../../error/app_exception.dart';
import '../../utils/logger.dart';

class WebSocketService {
  WebSocketService({required this.endpoint});

  final String endpoint;

  final StreamController<Map<String, dynamic>> _eventsController =
      StreamController<Map<String, dynamic>>.broadcast();

  WebSocketChannel? _channel;
  StreamSubscription<dynamic>? _subscription;
  bool _isConnected = false;
  String? _connectedRoomId;
  String? _connectedPlayerId;

  bool get isConnected => _isConnected;
  Stream<Map<String, dynamic>> get events => _eventsController.stream;

  Future<void> connect({
    String roomId = 'lobby',
    String playerId = 'guest',
  }) async {
    if (endpoint.isEmpty) {
      throw AppException('WebSocket endpoint is not configured.');
    }

    if (_isConnected &&
        roomId == _connectedRoomId &&
        playerId == _connectedPlayerId) {
      return;
    }

    await disconnect();

    final uri = _buildUri(roomId: roomId, playerId: playerId);
    Logger.info('WebSocket connect -> $uri');

    try {
      _channel = WebSocketChannel.connect(uri);
      _subscription = _channel!.stream.listen(
        _onEvent,
        onError: (Object error) {
          _isConnected = false;
          if (!_eventsController.isClosed) {
            _eventsController.addError(error);
          }
        },
        onDone: () {
          _isConnected = false;
        },
      );

      _isConnected = true;
      _connectedRoomId = roomId;
      _connectedPlayerId = playerId;
    } catch (error) {
      throw AppException('Failed to connect WebSocket: $error');
    }
  }

  Future<void> disconnect() async {
    await _subscription?.cancel();
    _subscription = null;

    await _channel?.sink.close();
    _channel = null;

    _isConnected = false;
    _connectedRoomId = null;
    _connectedPlayerId = null;
  }

  Future<void> send(String channel, Map<String, dynamic> payload) async {
    if (!_isConnected || _channel == null) {
      throw AppException('WebSocket is not connected.');
    }

    final message = <String, dynamic>{'channel': channel, 'payload': payload};

    Logger.info('WebSocket send [$channel] payload=$payload');
    _channel!.sink.add(jsonEncode(message));
  }

  void _onEvent(dynamic data) {
    if (_eventsController.isClosed) {
      return;
    }

    try {
      if (data is String) {
        final decoded = jsonDecode(data);
        if (decoded is Map<String, dynamic>) {
          _eventsController.add(decoded);
        }
      }
    } catch (_) {
      // Ignore malformed messages and keep connection alive.
    }
  }

  Uri _buildUri({required String roomId, required String playerId}) {
    final uri = Uri.parse(endpoint);
    final params = Map<String, String>.from(uri.queryParameters);
    params['roomId'] = roomId;
    params['playerId'] = playerId;
    return uri.replace(queryParameters: params);
  }

  void dispose() {
    unawaited(disconnect());
    _eventsController.close();
  }
}
