import '../core/config/app_env.dart';
import '../core/network/graphql/graphql_service.dart';
import '../core/network/websocket/websocket_service.dart';
import '../features/auth/data/datasources/auth_remote_data_source.dart';
import '../features/auth/data/repositories/auth_repository_impl.dart';
import '../features/auth/domain/repositories/auth_repository.dart';
import '../features/game/data/datasources/game_remote_data_source.dart';
import '../features/game/data/repositories/game_repository_impl.dart';
import '../features/game/domain/repositories/game_repository.dart';
import '../features/lobby/data/datasources/lobby_remote_data_source.dart';
import '../features/lobby/data/repositories/lobby_repository_impl.dart';
import '../features/lobby/domain/repositories/lobby_repository.dart';
import '../features/profile/data/datasources/profile_remote_data_source.dart';
import '../features/profile/data/repositories/profile_repository_impl.dart';
import '../features/profile/domain/repositories/profile_repository.dart';
import '../features/replay/data/datasources/replay_remote_data_source.dart';
import '../features/replay/data/repositories/replay_repository_impl.dart';
import '../features/replay/domain/repositories/replay_repository.dart';

class AppDependencies {
  AppDependencies._({
    required this.graphqlService,
    required this.webSocketService,
    required this.authRepository,
    required this.lobbyRepository,
    required this.gameRepository,
    required this.profileRepository,
    required this.replayRepository,
    required GameRemoteDataSource gameRemoteDataSource,
  }) : _gameRemoteDataSource = gameRemoteDataSource;

  final GraphqlService graphqlService;
  final WebSocketService webSocketService;

  final AuthRepository authRepository;
  final LobbyRepository lobbyRepository;
  final GameRepository gameRepository;
  final ProfileRepository profileRepository;
  final ReplayRepository replayRepository;

  final GameRemoteDataSource _gameRemoteDataSource;

  factory AppDependencies.create() {
    final graphqlService = GraphqlService(endpoint: AppEnv.graphqlEndpoint);
    final webSocketService = WebSocketService(
      endpoint: AppEnv.websocketEndpoint,
    );

    final authRepository = AuthRepositoryImpl(
      remoteDataSource: AuthRemoteDataSource(graphqlService: graphqlService),
    );

    final lobbyRepository = LobbyRepositoryImpl(
      remoteDataSource: LobbyRemoteDataSource(
        graphqlService: graphqlService,
        webSocketService: webSocketService,
      ),
    );

    final gameRemoteDataSource = GameRemoteDataSource(
      graphqlService: graphqlService,
      webSocketService: webSocketService,
    );
    final gameRepository = GameRepositoryImpl(
      remoteDataSource: gameRemoteDataSource,
    );

    final profileRepository = ProfileRepositoryImpl(
      remoteDataSource: ProfileRemoteDataSource(graphqlService: graphqlService),
    );

    final replayRepository = ReplayRepositoryImpl(
      remoteDataSource: ReplayRemoteDataSource(graphqlService: graphqlService),
    );

    return AppDependencies._(
      graphqlService: graphqlService,
      webSocketService: webSocketService,
      authRepository: authRepository,
      lobbyRepository: lobbyRepository,
      gameRepository: gameRepository,
      profileRepository: profileRepository,
      replayRepository: replayRepository,
      gameRemoteDataSource: gameRemoteDataSource,
    );
  }

  void dispose() {
    _gameRemoteDataSource.dispose();
    webSocketService.dispose();
  }
}
