import '../../../../core/network/graphql/graphql_service.dart';
import '../../../../core/network/websocket/websocket_service.dart';
import '../../domain/entities/room.dart';

class LobbyRemoteDataSource {
  LobbyRemoteDataSource({
    required GraphqlService graphqlService,
    required WebSocketService webSocketService,
  }) : _graphqlService = graphqlService,
       _webSocketService = webSocketService;

  final GraphqlService _graphqlService;
  final WebSocketService _webSocketService;

  Future<void> connect() async {
    await _webSocketService.connect();
  }

  Future<void> disconnect() async {
    await _webSocketService.disconnect();
  }

  Future<List<Room>> fetchRooms() async {
    await _graphqlService.query(
      document:
          'query Rooms { rooms { id name playersCount maxPlayers isPrivate } }',
    );

    // Mocked lobby while backend is pending.
    return const [
      Room(id: 'room_1', name: 'Mesa 1', playersCount: 2, maxPlayers: 4),
      Room(id: 'room_2', name: 'Mesa 2', playersCount: 4, maxPlayers: 4),
      Room(
        id: 'room_3',
        name: 'Treino',
        playersCount: 1,
        maxPlayers: 4,
        isPrivate: true,
      ),
    ];
  }
}
