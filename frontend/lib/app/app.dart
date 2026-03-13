import 'package:flutter/material.dart';

import 'app_dependencies.dart';
import 'app_routes.dart';
import 'app_theme.dart';

class SuecaApp extends StatefulWidget {
  const SuecaApp({super.key, required this.dependencies});

  final AppDependencies dependencies;

  @override
  State<SuecaApp> createState() => _SuecaAppState();
}

class _SuecaAppState extends State<SuecaApp> {
  late final AppRouter _router;

  @override
  void initState() {
    super.initState();
    _router = AppRouter(widget.dependencies);
  }

  @override
  void dispose() {
    widget.dependencies.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Sueca Online',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light(),
      initialRoute: AppRoutes.nickname,
      onGenerateRoute: _router.onGenerateRoute,
    );
  }
}
