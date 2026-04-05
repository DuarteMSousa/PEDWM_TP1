import '../../domain/entities/card.dart';
import '../../domain/entities/sueca_game_state.dart';
import '../../domain/repositories/game_repository.dart';
import '../datasources/game_remote_data_source.dart';

class GameRepositoryImpl implements GameRepository {
  GameRepositoryImpl({required GameRemoteDataSource remoteDataSource})
    : _remoteDataSource = remoteDataSource;

  final GameRemoteDataSource _remoteDataSource;

  @override
  Future<SuecaGameState> loadGame({
    required String roomId,
    required String playerId,
  }) {
    return _remoteDataSource.loadGame(roomId: roomId, playerId: playerId);
  }

  @override
  Future<SuecaGameState> playCard({
    required String roomId,
    required SuecaCard card,
  }) {
    return _remoteDataSource.playCard(roomId: roomId, card: card);
  }

  @override
  Stream<SuecaGameState> watchGame(String roomId) {
    return _remoteDataSource.watchGame(roomId);
  }

  @override
  Future<void> disconnect(String roomId) {
    return _remoteDataSource.disconnect(roomId);
  }
}
