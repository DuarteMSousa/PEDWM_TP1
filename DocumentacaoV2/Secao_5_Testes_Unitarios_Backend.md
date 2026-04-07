# 5. Validacao e Testes Unitarios (Backend)

## 5.1 Objetivo
O objetivo dos testes unitarios e validar, de forma isolada, as regras de negocio do backend, reduzindo regressos funcionais durante evolucao do sistema.

Neste trabalho, os testes focam dois niveis:
- Dominio (`internal/domain`): regras de cartas, jogadas, rondas, salas e estados.
- Servicos de aplicacao (`internal/application/services`): casos de uso de autenticacao, salas, eventos e estatisticas.

## 5.2 Estrategia de Teste
### 5.2.1 Tipo e Escopo

**Tabela 5.1 - Estrategia geral de testes unitarios**

| Campo | Descricao |
|---|---|
| Tipo de teste | Unitario |
| Ferramenta | `go test` |
| Escopo principal | `internal/domain/*` e `internal/application/services/*` |
| Dependencias externas | Substituidas por fakes em memoria (repositorios) |
| Criterio de sucesso | Todos os testes passam (`PASS`) |

### 5.2.2 Procedimento de execucao

Comandos utilizados:

```bash
go test ./...
go test ./... -cover
go test ./internal/application/services -cover
```

Para detalhe por funcao nos services:

```bash
go test ./internal/application/services -coverprofile=services.cover.out
go tool cover -func=services.cover.out
```

## 5.3 Casos de Teste Implementados
Foram implementados testes unitarios em modulos de dominio e de servicos, cobrindo cenarios de sucesso e erro.

**Tabela 5.2 - Casos de teste (lista completa)**

| ID | Ficheiro | Caso de teste (funcao) | Resultado esperado |
|---|---|---|---|
| TU-001 | `services/user_stats_service_test.go` | `TestUserStatsServiceGetByUserID` | Devolve stats existentes; erro quando user nao existe |
| TU-002 | `services/user_stats_service_test.go` | `TestUserStatsServiceRecordGameCreatesWhenMissing` | Cria stats em falta e atualiza jogos/vitorias/ELO |
| TU-003 | `services/user_stats_service_test.go` | `TestUserStatsServiceRecordGameSaveError` | Propaga erro de persistencia |
| TU-004 | `domain/card/card_test.go` | `TestNewCardValid` | Cria carta valida com ID correto |
| TU-005 | `domain/card/card_test.go` | `TestNewCardInvalidSuit` | Erro `ErrInvalidSuit` |
| TU-006 | `domain/card/card_test.go` | `TestNewCardInvalidRank` | Erro `ErrInvalidRank` |
| TU-007 | `domain/card/card_test.go` | `TestCardValidateInvalidID` | Erro `ErrInvalidCardID` |
| TU-008 | `domain/card/card_test.go` | `TestCardIsTrump` | Avalia corretamente carta de trunfo |
| TU-009 | `services/user_service_test.go` | `TestUserServiceRegisterRejectsExistingUsername` | Rejeita username duplicado |
| TU-010 | `services/user_service_test.go` | `TestUserServiceRegisterSuccessCreatesUserAndStats` | Regista user e cria estatisticas iniciais |
| TU-011 | `services/user_service_test.go` | `TestUserServiceLogin` | Login valido passa; invalido falha |
| TU-012 | `services/user_service_test.go` | `TestUserServiceGetUserNotFound` | Erro quando user nao existe |
| TU-013 | `domain/events/event_bus_test.go` | `TestEventBusPublishToSubscribers` | Publica evento para subscritores |
| TU-014 | `domain/events/event_bus_test.go` | `TestEventBusUnsubscribeStopsReceiving` | Nao entrega eventos apos unsubscribe |
| TU-015 | `services/room_service_test.go` | `TestRoomServiceCreateRoom` | Cria sala, regista no hub e persiste |
| TU-016 | `services/room_service_test.go` | `TestRoomServiceJoinRoom` | Join valido; erros de sala/user/duplicado |
| TU-017 | `services/room_service_test.go` | `TestRoomServiceLeaveRoom` | Leave valida; persiste OPEN e salta CLOSED |
| TU-018 | `services/room_service_test.go` | `TestRoomServiceDeleteRoom` | Remove sala no hub e repositorio |
| TU-019 | `services/room_service_test.go` | `TestRoomServiceGetRoomAndGetRooms` | Resolve sala por hub/repositorio e lista salas |
| TU-020 | `services/room_service_test.go` | `TestRoomServiceStartGame` | Inicia jogo e persiste room/game |
| TU-021 | `services/room_service_test.go` | `TestRoomServiceGetGameSnapshot` | Devolve snapshot; valida erros room/game |
| TU-022 | `services/game_service_test.go` | `TestGameServiceGetUserGames` | Lista jogos por user e trata erro repo |
| TU-023 | `services/game_service_test.go` | `TestGameServiceSetGameStatus` | Atualiza status e persiste |
| TU-024 | `services/event_service_test.go` | `TestEventServiceSaveEvent` | Persiste evento e propaga erro |
| TU-025 | `services/event_service_test.go` | `TestEventServiceGetEventsByRoomAndGame` | Consulta eventos por room/game e trata erros |
| TU-026 | `domain/deck/deck_test.go` | `TestDeckFirstDrawRemainingAndEmpty` | `First/Draw/Remaining` corretos; erro em deck vazio |
| TU-027 | `domain/deck/deck_test.go` | `TestDeckResetAndIsEmpty` | `Reset` limpa deck e marca vazio |
| TU-028 | `domain/trick/trick_test.go` | `TestTrickAddPlaySetsLeadSuitAndEnforcesTurnOrder` | Define naipe inicial e valida ordem de turno |
| TU-029 | `domain/trick/trick_test.go` | `TestTrickAddPlayRejectsDuplicatePlayer` | Rejeita jogada repetida do mesmo player |
| TU-030 | `domain/trick/trick_test.go` | `TestValidatePlayRequiresFollowingLeadSuit` | Obriga a seguir naipe quando possivel |
| TU-031 | `domain/trick/trick_test.go` | `TestValidatePlayAllowsDifferentSuitWhenPlayerDoesNotHaveLeadSuit` | Permite outro naipe sem carta do naipe lider |
| TU-032 | `domain/trick/trick_test.go` | `TestWinningPlayerAndTeamTrumpBeatsLeadSuit` | Trunfo vence e equipa vencedora correta |
| TU-033 | `domain/trick/trick_test.go` | `TestSuecaTrickScoringPoints` | Soma pontos da vaza corretamente |
| TU-034 | `domain/round/sueca_round_rule_strategy_test.go` | `TestSuecaRoundRuleStrategyHasEnded` | Identifica fim de ronda por maos vazias |
| TU-035 | `domain/round/sueca_round_rule_strategy_test.go` | `TestSuecaRoundRuleStrategyWinner` | Devolve vencedor correto da ronda |
| TU-036 | `domain/round/sueca_round_rule_strategy_test.go` | `TestCalculateCurrentTrickRoundPoints` | Atribui pontos da vaza a equipa vencedora |
| TU-037 | `domain/hand/hand_test.go` | `TestHandAddGetRemoveLifecycle` | Ciclo Add/Get/Remove consistente |
| TU-038 | `domain/hand/hand_test.go` | `TestHandRemoveCardMissing` | Erro ao remover carta inexistente |
| TU-039 | `domain/hand/hand_test.go` | `TestHandHasSuit` | Deteta corretamente existencia de naipe |
| TU-040 | `domain/room/room_test.go` | `TestNewRoomValidationAndHost` | Valida dados de sala e cria host inicial |
| TU-041 | `domain/room/room_test.go` | `TestAddPlayerPublishesEventAndPreventsDuplicate` | Publica evento de join e bloqueia duplicado |
| TU-042 | `domain/room/room_test.go` | `TestRemovePlayerClosesEmptyRoomAndPublishesCloseEvent` | Fecha sala vazia e publica `ROOM_CLOSED` |
| TU-043 | `domain/room/room_test.go` | `TestCreateGameSetsInGameAndCreatesFourPlayers` | Cria jogo e completa com bots ate 4 jogadores |
| TU-044 | `domain/game/game_test.go` | `TestNewGameInitialStateAndMaps` | Inicializa game com estado/mapas corretos |
| TU-045 | `domain/game/game_test.go` | `TestAddEventFillsGameAndRoomAndSequenceAndPublishes` | Preenche `GameID/RoomID/Sequence` e publica |
| TU-046 | `domain/game/game_test.go` | `TestPlayCardErrorGuards` | Protecoes para game nulo/estado invalido |
| TU-047 | `domain/game/game_test.go` | `TestPlayCardReturnsPlayerNotFoundBeforeRoundExecution` | Falha com player inexistente antes de jogar |
| TU-048 | `domain/game/game_test.go` | `TestRemovePlayerReplacesWithBotAndAddsEvents` | Substitui por bot e gera eventos esperados |
| TU-049 | `domain/game/game_test.go` | `TestGetPlayerAndGetPlayerTeamErrors` | Erros de pesquisa para player/equipa ausentes |
| TU-050 | `domain/game/game_test.go` | `TestSuecaGameScoringStrategyThresholdsAndWinner` | Regras de pontuacao e vencedor corretas |
| TU-051 | `domain/turnorder/turnorder_test.go` | `TestNewTurnOrderRejectsInvalidSize` | Rejeita numero invalido de jogadores |
| TU-052 | `domain/turnorder/turnorder_test.go` | `TestNewTurnOrderRejectsDuplicateIDs` | Rejeita IDs duplicados |
| TU-053 | `domain/turnorder/turnorder_test.go` | `TestNewTurnOrderStartsFromLeaderAndFollowsSequence` | Ordem circular correta a partir do lider |
| TU-054 | `domain/turnorder/turnorder_test.go` | `TestTurnOrderRemoveAndContains` | Remove jogador e atualiza presenca |
| TU-055 | `domain/player/botStrategy/bot_strategy_test.go` | `TestEasyBotChooseCardWithLeadSuit` | Easy escolhe carta do naipe lider |
| TU-056 | `domain/player/botStrategy/bot_strategy_test.go` | `TestEasyBotChooseFirstWhenNoLeadSuitInHand` | Easy joga primeira carta se nao tiver naipe |
| TU-057 | `domain/player/botStrategy/bot_strategy_test.go` | `TestHardBotChooseStrongestFromLeadSuit` | Hard escolhe carta mais forte do naipe lider |
| TU-058 | `domain/player/botStrategy/bot_strategy_test.go` | `TestHardBotChooseStrongestOverallWhenNoLeadSuit` | Hard escolhe carta globalmente mais forte |

## 5.4 Resultados de Execucao
Todos os testes unitarios implementados passaram com sucesso.

Data da medicao: **07/04/2026**.

**Tabela 5.3 - Cobertura global e por camada**

| Indicador | Valor |
|---|---|
| Cobertura total backend (`go test ./...`) | **11.3%** |
| Cobertura `internal/application/services` | **84.9%** |

### 5.4.1 Cobertura por funcao (Services)

**Tabela 5.4 - Cobertura por funcao em `internal/application/services`**

| Ficheiro / Funcao | Cobertura |
|---|---|
| `event_service.NewEventService` | 100.0% |
| `event_service.SaveEvent` | 100.0% |
| `event_service.GetEventsByRoomID` | 100.0% |
| `event_service.GetEventsByGameID` | 100.0% |
| `game_service.NewGameService` | 100.0% |
| `game_service.GetUserGames` | 100.0% |
| `game_service.SetGameStatus` | 81.8% |
| `room_service.NewRoomService` | 100.0% |
| `room_service.CreateRoom` | 90.0% |
| `room_service.JoinRoom` | 88.2% |
| `room_service.LeaveRoom` | 71.4% |
| `room_service.DeleteRoom` | 100.0% |
| `room_service.StartGame` | 66.7% |
| `room_service.GetRoom` | 66.7% |
| `room_service.GetRooms` | 100.0% |
| `room_service.GetGameSnapshot` | 92.6% |
| `user_service.NewUserService` | 100.0% |
| `user_service.Register` | 66.7% |
| `user_service.Login` | 100.0% |
| `user_service.GetUser` | 80.0% |
| `user_service.GetUserByUsername` | 80.0% |
| `user_stats_service.NewUserStatsService` | 100.0% |
| `user_stats_service.GetByUserID` | 100.0% |
| `user_stats_service.RecordGame` | 85.7% |
| **Total da camada services** | **84.9%** |

## 5.5 Analise Critica
O resultado de cobertura da camada de servicos (84.9%) indica boa confianca na logica de casos de uso do backend.

A cobertura total do backend (11.3%) e inferior porque inclui pacotes ainda sem testes unitarios:
- infraestrutura (`graph`, `websocket`, `persistence/postgres`);
- codigo gerado (gqlgen);
- componentes sem cenarios unitarios definidos.

Assim, a metrica total nao reflete apenas a qualidade dos testes de negocio, mas tambem a presenca de camadas tecnicas nao cobertas.

## 5.6 Trabalhos Futuros
Para aumentar a cobertura global de forma util:
- adicionar testes de integracao para `graphql` e `websocket`;
- cobrir persistencia com base de dados de teste;
- incluir casos adicionais de estado de ronda/jogo (transicoes completas).
