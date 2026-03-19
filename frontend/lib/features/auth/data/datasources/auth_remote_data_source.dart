import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/user.dart';

class AuthRemoteDataSource {
  AuthRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<User> enterWithNickname(String nickname) async {
    final cleanedNickname = nickname.trim();
    if (cleanedNickname.isEmpty) {
      throw AppException('Nickname is required.');
    }

    final response = await _graphqlService.mutation(
      document: '''
        mutation CreatePlayer(\$nickname: String!) {
          createPlayer(nickname: \$nickname) {
            id
            nickname
          }
        }
      ''',
      variables: <String, dynamic>{'nickname': cleanedNickname},
    );

    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for createPlayer.');
    }

    final payload = data['createPlayer'];
    if (payload is! Map<String, dynamic>) {
      throw AppException('Player not returned by GraphQL.');
    }

    return User(
      id: payload['id']?.toString() ?? '',
      nickname: payload['nickname']?.toString() ?? cleanedNickname,
    );
  }
}
