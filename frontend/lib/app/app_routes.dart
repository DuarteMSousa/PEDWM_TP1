import 'package:flutter/material.dart';

import '../features/auth/domain/entities/user.dart';
import '../features/auth/presentation/pages/nickname_page.dart';
import '../features/game/presentation/pages/game_page.dart';
import '../features/lobby/domain/entities/room.dart';
import '../features/lobby/presentation/pages/lobby_page.dart';
import '../features/profile/presentation/pages/profile_page.dart';
import 'app_dependencies.dart';

class AppRoutes {
  AppRoutes._();

  static const nickname = '/';
  static const lobby = '/lobby';
  static const game = '/game';
  static const profile = '/profile';
}

class AppRouter {
  AppRouter(this._dependencies);

  final AppDependencies _dependencies;

  Route<dynamic> onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case AppRoutes.nickname:
        return _buildAnimatedRoute(
          settings: settings,
          beginOffset: const Offset(-0.03, 0),
          builder: (_) =>
              NicknamePage(authRepository: _dependencies.authRepository),
        );
      case AppRoutes.lobby:
        final user = settings.arguments is User
            ? settings.arguments! as User
            : const User(id: 'guest', nickname: 'Guest');
        return _buildAnimatedRoute(
          settings: settings,
          beginOffset: const Offset(0.06, 0),
          builder: (_) => LobbyPage(
            lobbyRepository: _dependencies.lobbyRepository,
            currentUser: user,
          ),
        );
      case AppRoutes.game:
        if (settings.arguments is! Room) {
          return _errorRoute('Missing Room argument for /game');
        }
        return _buildAnimatedRoute(
          settings: settings,
          beginOffset: const Offset(0, 0.07),
          builder: (_) => GamePage(
            gameRepository: _dependencies.gameRepository,
            room: settings.arguments! as Room,
          ),
        );
      case AppRoutes.profile:
        final userId = settings.arguments is String
            ? settings.arguments! as String
            : 'guest';
        return _buildAnimatedRoute(
          settings: settings,
          beginOffset: const Offset(0.05, 0),
          builder: (_) => ProfilePage(
            profileRepository: _dependencies.profileRepository,
            userId: userId,
          ),
        );
      default:
        return _errorRoute('Route not found: ${settings.name}');
    }
  }

  PageRouteBuilder<void> _buildAnimatedRoute({
    required RouteSettings settings,
    required WidgetBuilder builder,
    Offset beginOffset = const Offset(0.04, 0),
  }) {
    return PageRouteBuilder<void>(
      settings: settings,
      transitionDuration: const Duration(milliseconds: 340),
      reverseTransitionDuration: const Duration(milliseconds: 260),
      pageBuilder: (context, animation, secondaryAnimation) => builder(context),
      transitionsBuilder: (context, animation, secondaryAnimation, child) {
        final curved = CurvedAnimation(
          parent: animation,
          curve: Curves.easeOutCubic,
          reverseCurve: Curves.easeInCubic,
        );
        return SlideTransition(
          position: Tween<Offset>(
            begin: beginOffset,
            end: Offset.zero,
          ).animate(curved),
          child: FadeTransition(opacity: curved, child: child),
        );
      },
    );
  }

  MaterialPageRoute<void> _errorRoute(String message) {
    return MaterialPageRoute<void>(
      builder: (_) => Scaffold(
        appBar: AppBar(title: const Text('Navigation error')),
        body: Center(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Text(message, textAlign: TextAlign.center),
          ),
        ),
      ),
    );
  }
}
