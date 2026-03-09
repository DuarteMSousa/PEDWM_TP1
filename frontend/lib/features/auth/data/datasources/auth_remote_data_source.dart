import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/user.dart';

class AuthRemoteDataSource {
  AuthRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<User> login(String nickname) async {
    final cleanedNickname = nickname.trim();
    if (cleanedNickname.isEmpty) {
      throw AppException('Nickname is required.');
    }

    await _graphqlService.mutation(
      document:
          'mutation Login(\$nickname: String!) { login(nickname: \$nickname) { id nickname } }',
      variables: <String, dynamic>{'nickname': cleanedNickname},
    );

    return User(
      id: cleanedNickname.toLowerCase().replaceAll(' ', '_'),
      nickname: cleanedNickname,
    );
  }
}
