import '../../../../core/network/http/http_service.dart';
import '../../../../core/network/websocket/websocket_service.dart';
import '../../domain/entities/room.dart';

class LobbyRemoteDataSource {
  LobbyRemoteDataSource({
    required HttpService httpService,
    required WebSocketService webSocketService,
  }) : _httpService = httpService,
       _webSocketService = webSocketService;

  final HttpService _httpService;
  final WebSocketService _webSocketService;

  Future<void> connect() async {
    await _webSocketService.connect(roomId: 'lobby', playerId: 'lobby_client');
  }

  Future<void> disconnect() async {
    await _webSocketService.disconnect();
  }

  Future<List<Room>> fetchRooms() async {
    final rawRooms = await _httpService.getList('/rooms');
    return rawRooms
        .whereType<Map<String, dynamic>>()
        .map(_roomFromJson)
        .toList(growable: false);
  }

  Future<Room> joinRoom({
    required String roomId,
    required String playerId,
  }) async {
    final response = await _httpService.post(
      '/rooms/$roomId/join',
      body: <String, dynamic>{'playerId': playerId},
    );

    return _roomFromJson(response);
  }

  Room _roomFromJson(Map<String, dynamic> json) {
    return Room(
      id: json['id']?.toString() ?? '',
      name: json['name']?.toString() ?? 'Sala',
      playersCount: (json['playersCount'] as num?)?.toInt() ?? 0,
      maxPlayers: (json['maxPlayers'] as num?)?.toInt() ?? 4,
      isPrivate: json['isPrivate'] == true,
    );
  }
}
