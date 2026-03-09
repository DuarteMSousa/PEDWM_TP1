import 'dart:async';

import '../../error/app_exception.dart';
import '../../utils/logger.dart';

class WebSocketService {
  WebSocketService({required this.endpoint});

  final String endpoint;

  final StreamController<Map<String, dynamic>> _eventsController =
      StreamController<Map<String, dynamic>>.broadcast();
  bool _isConnected = false;

  bool get isConnected => _isConnected;
  Stream<Map<String, dynamic>> get events => _eventsController.stream;

  Future<void> connect() async {
    if (_isConnected) {
      return;
    }
    if (endpoint.isEmpty) {
      throw AppException('WebSocket endpoint is not configured.');
    }
    Logger.info('WebSocket connect -> $endpoint');
    await Future<void>.delayed(const Duration(milliseconds: 180));
    _isConnected = true;
  }

  Future<void> disconnect() async {
    if (!_isConnected) {
      return;
    }
    Logger.info('WebSocket disconnect');
    await Future<void>.delayed(const Duration(milliseconds: 120));
    _isConnected = false;
  }

  Future<void> send(String channel, Map<String, dynamic> payload) async {
    if (!_isConnected) {
      throw AppException('WebSocket is not connected.');
    }
    Logger.info('WebSocket send [$channel] payload=$payload');
    await Future<void>.delayed(const Duration(milliseconds: 70));
  }

  void pushMockEvent(Map<String, dynamic> event) {
    if (_eventsController.isClosed) {
      return;
    }
    _eventsController.add(event);
  }

  void dispose() {
    _eventsController.close();
  }
}
