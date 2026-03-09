import 'package:flutter/foundation.dart';

class Logger {
  Logger._();

  static void info(String message) {
    debugPrint('[INFO] ${DateTime.now().toIso8601String()} $message');
  }

  static void warn(String message) {
    debugPrint('[WARN] ${DateTime.now().toIso8601String()} $message');
  }
}
