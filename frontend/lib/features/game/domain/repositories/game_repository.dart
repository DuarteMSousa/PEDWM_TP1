import '../entities/card.dart';
import '../entities/sueca_game_state.dart';

abstract class GameRepository {
  Future<SuecaGameState> loadGame({
    required String roomId,
    required String playerId,
  });
  Future<SuecaGameState> playCard({
    required String roomId,
    required SuecaCard card,
  });
  Stream<SuecaGameState> watchGame(String roomId);
}
