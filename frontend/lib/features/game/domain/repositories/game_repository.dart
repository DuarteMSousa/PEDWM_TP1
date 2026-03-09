import '../entities/card.dart';
import '../entities/sueca_game_state.dart';

abstract class GameRepository {
  Future<SuecaGameState> loadGame(String roomId);
  Future<SuecaGameState> playCard({
    required String roomId,
    required SuecaCard card,
  });
  Stream<SuecaGameState> watchGame(String roomId);
}
