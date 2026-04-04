import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/profile.dart';

class ProfileRemoteDataSource {
  ProfileRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<Profile> fetchProfile(String userId) async {
    final response = await _graphqlService.query(
      document: '''
        query Profile(\$userId: ID!) {
          user(id: \$userId) {
            id
            username
          }
          userStats(userId: \$userId) {
            games
            wins
          }
        }
      ''',
      variables: <String, dynamic>{'userId': userId},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for profile query.');
    }

    final userPayload = data['user'];
    final statsPayload = data['userStats'];

    if (userPayload is! Map<String, dynamic>) {
      throw AppException('User profile was not returned by GraphQL.');
    }

    final nickname = userPayload['username']?.toString() ?? userId;
    final matchesPlayed = (statsPayload is Map<String, dynamic>)
        ? (statsPayload['games'] as num?)?.toInt() ?? 0
        : 0;
    final wins = (statsPayload is Map<String, dynamic>)
        ? (statsPayload['wins'] as num?)?.toInt() ?? 0
        : 0;

    return Profile(
      userId: userId,
      nickname: nickname,
      matchesPlayed: matchesPlayed,
      wins: wins,
    );
  }
}
