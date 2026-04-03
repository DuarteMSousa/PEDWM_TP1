import '../entities/profile.dart';
import '../entities/match_replay.dart';

abstract class ProfileRepository {
  Future<Profile> fetchProfile(String userId);
  Future<MatchReplay?> fetchReplay(String gameId);
}
