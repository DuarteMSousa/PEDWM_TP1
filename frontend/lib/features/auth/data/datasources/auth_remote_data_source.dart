import '../../../../core/error/app_exception.dart';
import '../../../../core/network/graphql/graphql_service.dart';
import '../../domain/entities/user.dart';

class AuthRemoteDataSource {
  AuthRemoteDataSource({required GraphqlService graphqlService})
    : _graphqlService = graphqlService;

  final GraphqlService _graphqlService;

  Future<User> login({
    required String username,
    required String password,
  }) async {
    final cleanedUsername = username.trim();
    if (cleanedUsername.isEmpty) {
      throw AppException('Nickname is required.');
    }
    if (password.isEmpty) {
      throw AppException('Password is required.');
    }

    final response = await _graphqlService.mutation(
      document: '''
        mutation Login(\$username: String!, \$password: String!) {
          login(input: { username: \$username, password: \$password }) {
            user {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{
        'username': cleanedUsername,
        'password': password,
      },
    );
    return _userFromAuthPayload(response, operationName: 'login');
  }

  Future<User> register({
    required String username,
    required String password,
  }) async {
    final cleanedUsername = username.trim();
    if (cleanedUsername.isEmpty) {
      throw AppException('Nickname is required.');
    }
    if (password.isEmpty) {
      throw AppException('Password is required.');
    }
    if (password.length < 6) {
      throw AppException('Password must be at least 6 characters.');
    }

    final response = await _graphqlService.mutation(
      document: '''
        mutation Register(\$username: String!, \$password: String!) {
          register(input: { username: \$username, password: \$password }) {
            user {
              id
              username
            }
          }
        }
      ''',
      variables: <String, dynamic>{
        'username': cleanedUsername,
        'password': password,
      },
    );
    return _userFromAuthPayload(response, operationName: 'register');
  }

  User _userFromAuthPayload(
    Map<String, dynamic> response, {
    required String operationName,
  }) {
    final data = response['data'];
    if (data is! Map<String, dynamic>) {
      throw AppException('Invalid GraphQL response for $operationName.');
    }

    final payload = data[operationName];
    if (payload is! Map<String, dynamic>) {
      throw AppException('Auth payload missing for $operationName.');
    }

    final userPayload = payload['user'];
    if (userPayload is! Map<String, dynamic>) {
      throw AppException('User payload missing for $operationName.');
    }

    final id = userPayload['id']?.toString() ?? '';
    final username = userPayload['username']?.toString() ?? '';
    if (id.isEmpty || username.isEmpty) {
      throw AppException('Invalid user payload returned by $operationName.');
    }

    return User(id: id, nickname: username);
  }
}
