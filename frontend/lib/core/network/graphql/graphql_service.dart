import 'dart:convert';

import 'package:http/http.dart' as http;

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

    final uri = Uri.parse(endpoint);
    final payload = <String, dynamic>{
      'query': document,
      'variables': variables,
    };

    http.Response response;
    try {
      response = await http.post(
        uri,
        headers: <String, String>{'Content-Type': 'application/json'},
        body: jsonEncode(payload),
      );
    } catch (error) {
      throw AppException('GraphQL network error calling $uri: $error');
    }

    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw AppException(
        'GraphQL HTTP ${response.statusCode} on $uri: ${response.body}',
      );
    }

    if (response.body.isEmpty) {
      throw AppException('GraphQL empty response from $uri');
    }

    dynamic decoded;
    try {
      decoded = jsonDecode(response.body);
    } catch (_) {
      throw AppException('GraphQL invalid JSON response from $uri');
    }

    if (decoded is! Map<String, dynamic>) {
      throw AppException('GraphQL invalid payload from $uri');
    }

    final errors = decoded['errors'];
    if (errors is List && errors.isNotEmpty) {
      final first = errors.first;
      if (first is Map<String, dynamic>) {
        final message = first['message']?.toString();
        if (message != null && message.isNotEmpty) {
          throw AppException(message);
        }
      }
      throw AppException('GraphQL request failed.');
    }

    return decoded;
  }
}
