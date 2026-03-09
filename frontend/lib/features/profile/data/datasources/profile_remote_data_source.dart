import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/profile.dart';

class ProfileRemoteDataSource {
  ProfileRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<Profile> fetchProfile(String userId) async {
    await _graphqlService.query(
      document:
          'query Profile(\$userId: ID!) { profile(userId: \$userId) { id nickname matches wins } }',
      variables: <String, dynamic>{'userId': userId},
    );

    return Profile(
      userId: userId,
      nickname: userId == 'guest' ? 'Guest' : userId,
      matchesPlayed: 24,
      wins: 14,
    );
  }
}
