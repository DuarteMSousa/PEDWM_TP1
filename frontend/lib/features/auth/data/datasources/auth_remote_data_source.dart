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

    final password = _buildDeterministicPassword(cleanedNickname);

    final loginUser = await _tryLogin(
      username: cleanedNickname,
      password: password,
    );
    if (loginUser != null) {
      return loginUser;
    }

    try {
      return await _register(username: cleanedNickname, password: password);
    } catch (error) {
      if (_isUsernameAlreadyRegistered(error)) {
        final retriedLogin = await _tryLogin(
          username: cleanedNickname,
          password: password,
        );
        if (retriedLogin != null) {
          return retriedLogin;
        }
      }
      rethrow;
    }
  }

  Future<User?> _tryLogin({
    required String username,
    required String password,
  }) async {
    try {
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
          'username': username,
          'password': password,
        },
      );
      return _userFromAuthPayload(response, operationName: 'login');
    } catch (error) {
      if (_isInvalidCredentials(error)) {
        return null;
      }
      rethrow;
    }
  }

  Future<User> _register({
    required String username,
    required String password,
  }) async {
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
      variables: <String, dynamic>{'username': username, 'password': password},
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

  bool _isInvalidCredentials(Object error) {
    final message = error.toString().toLowerCase();
    return message.contains('invalid credentials');
  }

  bool _isUsernameAlreadyRegistered(Object error) {
    final message = error.toString().toLowerCase();
    return message.contains('username already exists');
  }

  String _buildDeterministicPassword(String username) {
    final safe = username.toLowerCase().replaceAll(RegExp(r'[^a-z0-9]'), '');
    final base = safe.isEmpty ? 'player' : safe;
    return 'pedwm_$base';
  }
}
