import '../entities/game_summary.dart';

abstract class ReplayRepository {
  Future<List<GameSummary>> fetchUserGames(String userId);
}
