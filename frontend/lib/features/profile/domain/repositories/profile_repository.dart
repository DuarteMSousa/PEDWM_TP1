import '../entities/profile.dart';

abstract class ProfileRepository {
  Future<Profile> fetchProfile(String userId);
}
