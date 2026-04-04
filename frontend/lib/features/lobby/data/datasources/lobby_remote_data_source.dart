import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../../../core/network/websocket/websocket_service.dart';
import '../../domain/entities/lobby_realtime_event.dart';
import '../../domain/entities/room.dart';
import '../../domain/entities/room_details.dart';
import '../../domain/entities/room_member.dart';

class LobbyRemoteDataSource {
  LobbyRemoteDataSource({
    required GraphqlService graphqlService,
    required WebSocketService webSocketService,
  }) : _graphqlService = graphqlService,
       _webSocketService = webSocketService;

  final GraphqlService _graphqlService;
  final WebSocketService _webSocketService;

  Future<void> connectLobby({required String playerId}) async {
    await _webSocketService.connect(roomId: 'lobby', playerId: playerId);
  }

  Future<void> connectRoom({
    required String roomId,
    required String playerId,
  }) async {
    await _webSocketService.connect(roomId: roomId, playerId: playerId);
  }

  Future<void> disconnect() async {
    await _webSocketService.disconnect();
  }

  Stream<LobbyRealtimeEvent> watchRealtimeEvents() {
    return _webSocketService.events
        .map(_mapRealtimeEvent)
        .where((event) => event != null)
        .cast<LobbyRealtimeEvent>();
  }

  Future<List<Room>> fetchRooms() async {
    final response = await _graphqlService.query(
      document: '''
        query Rooms {
          rooms {
            id
            hostId
            status
            players {
              id
              username
            }
          }
        }
      ''',
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for rooms query.');
    }

    final rawRooms = data['rooms'];
    if (rawRooms is! List) {
      throw AppException('Rooms were not returned by GraphQL.');
    }

    return rawRooms
        .whereType<Map<String, dynamic>>()
        .map(_roomFromJson)
        .toList(growable: false);
  }

  Future<Room> createRoom({required String hostPlayerId}) async {
    final response = await _graphqlService.mutation(
      document: '''
        mutation CreateRoom(\$hostId: ID!) {
          createRoom(input: { hostId: \$hostId }) {
            id
            hostId
            status
            players {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{'hostId': hostPlayerId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for createRoom.');
    }

    final payload = data['createRoom'];
    if (payload is! Map<String, dynamic>) {
      throw AppException('Room was not returned by GraphQL createRoom.');
    }

    return _roomFromJson(payload);
  }

  Future<RoomDetails?> fetchRoomDetails({required String roomId}) async {
    final response = await _graphqlService.query(
      document: '''
        query RoomDetails(\$roomId: ID!) {
          room(id: \$roomId) {
            id
            hostId
            status
            players {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{'roomId': roomId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for room query.');
    }

    final rawRoom = data['room'];
    if (rawRoom == null) {
      return null;
    }
    if (rawRoom is! Map<String, dynamic>) {
      throw AppException('Room was not returned by GraphQL room query.');
    }

    return _roomDetailsFromJson(rawRoom);
  }

  Future<Room> joinRoom({
    required String roomId,
    required String playerId,
  }) async {
    final response = await _graphqlService.mutation(
      document: '''
        mutation JoinRoom(\$roomId: ID!, \$userId: ID!) {
          joinRoom(input: { roomId: \$roomId, userId: \$userId }) {
            id
            hostId
            status
            players {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{'roomId': roomId, 'userId': playerId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for joinRoom.');
    }

    final payload = data['joinRoom'];
    if (payload is! Map<String, dynamic>) {
      throw AppException('Room was not returned by GraphQL joinRoom.');
    }

    return _roomFromJson(payload);
  }

  Future<RoomDetails> leaveRoom({
    required String roomId,
    required String playerId,
  }) async {
    final response = await _graphqlService.mutation(
      document: '''
        mutation LeaveRoom(\$roomId: ID!, \$userId: ID!) {
          leaveRoom(input: { roomId: \$roomId, userId: \$userId }) {
            id
            hostId
            status
            players {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{'roomId': roomId, 'userId': playerId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for leaveRoom.');
    }

    final payload = data['leaveRoom'];
    if (payload is! Map<String, dynamic>) {
      throw AppException('Room was not returned by GraphQL leaveRoom.');
    }

    return _roomDetailsFromJson(payload);
  }

  Future<RoomDetails> startGame({required String roomId}) async {
    final response = await _graphqlService.mutation(
      document: '''
        mutation StartGame(\$roomId: ID!) {
          startGame(input: { roomId: \$roomId }) {
            id
            hostId
            status
            players {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{'roomId': roomId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for startGame.');
    }

    final payload = data['startGame'];
    if (payload is! Map<String, dynamic>) {
      throw AppException('Room was not returned by GraphQL startGame.');
    }

    return _roomDetailsFromJson(payload);
  }

  Room _roomFromJson(Map<String, dynamic> json) {
    final players = _parsePlayers(json['players']);
    final id = json['id']?.toString() ?? '';
    return Room(
      id: id,
      name: _roomName(id),
      playersCount: players.length,
      maxPlayers: 4,
      status: json['status']?.toString() ?? 'OPEN',
    );
  }

  RoomDetails _roomDetailsFromJson(Map<String, dynamic> json) {
    final id = json['id']?.toString() ?? '';
    final players = _parsePlayers(json['players']);

    return RoomDetails(
      id: id,
      name: _roomName(id),
      hostPlayerId: json['hostId']?.toString() ?? '',
      status: json['status']?.toString() ?? 'OPEN',
      players: players,
    );
  }

  List<RoomMember> _parsePlayers(dynamic rawPlayers) {
    if (rawPlayers is! List) {
      return const <RoomMember>[];
    }

    return rawPlayers
        .whereType<Map<String, dynamic>>()
        .map(
          (item) => RoomMember(
            id: item['id']?.toString() ?? '',
            nickname: item['username']?.toString() ?? 'Guest',
          ),
        )
        .toList(growable: false);
  }

  String _roomName(String id) {
    final clean = id.trim();
    if (clean.isEmpty) {
      return 'Mesa';
    }
    final token = clean.length <= 6 ? clean : clean.substring(0, 6);
    return 'Mesa $token';
  }

  LobbyRealtimeEvent? _mapRealtimeEvent(Map<String, dynamic> rawEvent) {
    final type = rawEvent['type']?.toString();
    final roomId =
        rawEvent['roomId']?.toString() ?? rawEvent['gameId']?.toString();
    if (type == null || roomId == null || type.isEmpty || roomId.isEmpty) {
      return null;
    }

    final payload = rawEvent['payload'];
    final payloadMap = payload is Map
        ? Map<String, dynamic>.from(payload)
        : const <String, dynamic>{};

    return LobbyRealtimeEvent(
      id: rawEvent['id']?.toString() ?? '',
      type: type,
      roomId: roomId,
      timestamp: rawEvent['timestamp']?.toString() ?? '',
      payload: payloadMap,
    );
  }
}
