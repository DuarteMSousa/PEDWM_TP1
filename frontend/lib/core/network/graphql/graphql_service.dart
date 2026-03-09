import '../../error/app_exception.dart';
import '../../utils/logger.dart';

class GraphqlService {
  GraphqlService({required this.endpoint});

  final String endpoint;

  Future<Map<String, dynamic>> query({
    required String document,
    Map<String, dynamic> variables = const {},
  }) {
    return _execute(
      operation: 'query',
      document: document,
      variables: variables,
    );
  }

  Future<Map<String, dynamic>> mutation({
    required String document,
    Map<String, dynamic> variables = const {},
  }) {
    return _execute(
      operation: 'mutation',
      document: document,
      variables: variables,
    );
  }

  Future<Map<String, dynamic>> _execute({
    required String operation,
    required String document,
    required Map<String, dynamic> variables,
  }) async {
    if (endpoint.isEmpty) {
      throw AppException('GraphQL endpoint is not configured.');
    }

    Logger.info(
      'GraphQL $operation -> $endpoint | vars=${variables.keys.toList()}',
    );

    await Future<void>.delayed(const Duration(milliseconds: 140));

    // Stub response while backend is not connected.
    return <String, dynamic>{'data': <String, dynamic>{}};
  }
}
