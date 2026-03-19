import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../../../core/network/websocket/websocket_service.dart';
import '../../domain/entities/lobby_realtime_event.dart';
import '../../domain/entities/room_details.dart';
import '../../domain/entities/room_member.dart';
import '../../domain/entities/room.dart';

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
            name
            playersCount
            maxPlayers
            isPrivate
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

  Future<Room> createRoom({
    required String name,
    required String hostPlayerId,
    int maxPlayers = 4,
    bool isPrivate = false,
    String? password,
  }) async {
    final cleanedPassword = password?.trim();

    final response = await _graphqlService.mutation(
      document: '''
        mutation CreateRoom(
          \$name: String!
          \$hostPlayerId: ID!
          \$maxPlayers: Int
          \$isPrivate: Boolean
          \$password: String
        ) {
          createRoom(
            input: {
              name: \$name
              hostPlayerId: \$hostPlayerId
              maxPlayers: \$maxPlayers
              isPrivate: \$isPrivate
              password: \$password
            }
          ) {
            id
            name
            playersCount
            maxPlayers
            isPrivate
          }
        }
      ''',
      variables: <String, dynamic>{
        'name': name.trim(),
        'hostPlayerId': hostPlayerId,
        'maxPlayers': maxPlayers,
        'isPrivate': isPrivate,
        'password': cleanedPassword?.isEmpty == true ? null : cleanedPassword,
      },
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
            name
            hostPlayerId
            status
            maxPlayers
            isPrivate
            playersCount
            players {
              id
              nickname
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
    String? password,
  }) async {
    final cleanedPassword = password?.trim();

    final response = await _graphqlService.mutation(
      document: '''
        mutation JoinRoom(\$roomId: ID!, \$playerId: ID!, \$password: String) {
          joinRoom(roomId: \$roomId, playerId: \$playerId, password: \$password) {
            id
            name
            playersCount
            maxPlayers
            isPrivate
          }
        }
      ''',
      variables: <String, dynamic>{
        'roomId': roomId,
        'playerId': playerId,
        'password': cleanedPassword?.isEmpty == true ? null : cleanedPassword,
      },
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
        mutation LeaveRoom(\$roomId: ID!, \$playerId: ID!) {
          leaveRoom(roomId: \$roomId, playerId: \$playerId) {
            id
            name
            hostPlayerId
            status
            maxPlayers
            isPrivate
            playersCount
            players {
              id
              nickname
            }
          }
        }
      ''',
      variables: <String, dynamic>{'roomId': roomId, 'playerId': playerId},
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

  Future<bool> deleteRoom({
    required String roomId,
    required String requesterId,
  }) async {
    final response = await _graphqlService.mutation(
      document: '''
        mutation DeleteRoom(\$roomId: ID!, \$requesterId: ID!) {
          deleteRoom(roomId: \$roomId, requesterId: \$requesterId)
        }
      ''',
      variables: <String, dynamic>{
        'roomId': roomId,
        'requesterId': requesterId,
      },
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for deleteRoom.');
    }

    return data['deleteRoom'] == true;
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

  RoomDetails _roomDetailsFromJson(Map<String, dynamic> json) {
    final rawPlayers = json['players'];
    final players = rawPlayers is List
        ? rawPlayers
              .whereType<Map<String, dynamic>>()
              .map(
                (item) => RoomMember(
                  id: item['id']?.toString() ?? '',
                  nickname: item['nickname']?.toString() ?? 'Guest',
                ),
              )
              .toList(growable: false)
        : const <RoomMember>[];

    return RoomDetails(
      id: json['id']?.toString() ?? '',
      name: json['name']?.toString() ?? 'Sala',
      hostPlayerId: json['hostPlayerId']?.toString() ?? '',
      status: json['status']?.toString() ?? 'OPEN',
      maxPlayers: (json['maxPlayers'] as num?)?.toInt() ?? 4,
      isPrivate: json['isPrivate'] == true,
      players: players,
    );
  }

  LobbyRealtimeEvent? _mapRealtimeEvent(Map<String, dynamic> rawEvent) {
    final type = rawEvent['type']?.toString();
    final roomId = rawEvent['roomId']?.toString();
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
