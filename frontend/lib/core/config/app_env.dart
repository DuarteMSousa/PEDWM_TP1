class AppEnv {
  AppEnv._();

  static const apiBaseEndpoint = 'http://localhost:4000';
  static const graphqlEndpoint = '$apiBaseEndpoint/graphql';
  static const websocketEndpoint = 'ws://localhost:4000/ws';
}
