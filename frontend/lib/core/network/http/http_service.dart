import 'dart:convert';

import 'package:http/http.dart' as http;

import '../../error/app_exception.dart';
import '../../utils/logger.dart';

class HttpService {
  HttpService({required this.baseUrl});

  final String baseUrl;

  Future<List<dynamic>> getList(String path) async {
    final data = await _request('GET', path);
    if (data is! List<dynamic>) {
      throw AppException('Expected a JSON array from $path');
    }
    return data;
  }

  Future<Map<String, dynamic>> post(
    String path, {
    Map<String, dynamic> body = const {},
  }) async {
    final data = await _request('POST', path, body: body);
    if (data is! Map<String, dynamic>) {
      throw AppException('Expected a JSON object from $path');
    }
    return data;
  }

  Future<dynamic> _request(
    String method,
    String path, {
    Map<String, dynamic>? body,
  }) async {
    final uri = _buildUri(path);
    Logger.info('HTTP $method -> $uri');

    http.Response response;
    try {
      switch (method) {
        case 'GET':
          response = await http.get(uri);
          break;
        case 'POST':
          response = await http.post(
            uri,
            headers: {'Content-Type': 'application/json'},
            body: jsonEncode(body ?? <String, dynamic>{}),
          );
          break;
        default:
          throw AppException('Unsupported HTTP method: $method');
      }
    } catch (error) {
      throw AppException('Network error calling $uri: $error');
    }

    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw AppException(
        'HTTP ${response.statusCode} on $uri: ${response.body}',
      );
    }

    if (response.body.isEmpty) {
      return <String, dynamic>{};
    }

    try {
      return jsonDecode(response.body);
    } catch (_) {
      throw AppException('Invalid JSON response from $uri');
    }
  }

  Uri _buildUri(String path) {
    final normalizedBase = baseUrl.endsWith('/')
        ? baseUrl.substring(0, baseUrl.length - 1)
        : baseUrl;
    final normalizedPath = path.startsWith('/') ? path : '/$path';
    return Uri.parse('$normalizedBase$normalizedPath');
  }
}
