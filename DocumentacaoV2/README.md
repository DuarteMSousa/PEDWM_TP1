# DocumentacaoV2

Esta pasta contem os diagramas atualizados para o estado atual do projeto (backend + frontend).

## O que foi atualizado

- Arquitetura real do backend Go com GraphQL + WebSocket + EventBus + PostgreSQL.
- Arquitetura atual do frontend Flutter (camadas, features e dependencias).
- Modelo de dados real da base de dados (`users`, `friendships`, `user_stats`, `rooms`, `room_players`, `games`, `game_players`, `events`).
- Fluxos reais de sala, inicio de jogo, jogada de carta e replay.
- Fluxos WebSocket atuais (ligacao, comandos, broadcast e ciclo de vida).

## Estrutura

- `DiagramaComponents.puml`: arquitetura global frontend/backend.
- `DiagramaFrontendArquitetura.puml`: arquitetura interna do frontend Flutter.
- `DiagramaClasses.puml`: visao geral consolidada das classes principais.
- `DiagramaClassesSockets.puml`: classes da stack WebSocket.
- `DiagramaER.puml`: modelo relacional atual de PostgreSQL.
- `DiagramaEstados.puml`: estados agregados de Room/Game/Round.
- `DiagramaSequenciaEntrarnaSala.puml`: fluxo de autenticacao/lobby/start.
- `DiagramaSequenciaJogarCarta.puml`: fluxo realtime de jogada.
- `DiagramaSequenciaReplay.puml`: fluxo do replay no perfil.
- `classes/`: divisao do diagrama de classes por categorias da secao 3.2.
  - `01_dominio_jogo.puml`
  - `02_jogadores_e_equipas.puml`
  - `03_state_pattern.puml`
  - `04_strategy_pattern.puml`
  - `05_command_pattern.puml`
  - `06_observer_pattern.puml`
- `sequencias/`: sequencias de criar sala e ronda completa.
- `estados/`: estados detalhados de `Game` e `Round`.
- `WebSockets/`: documentacao dedicada a arquitetura e fluxo WS.

## Nota

Os diagramas foram alinhados com o codigo atual em `backend/internal` e `frontend/lib`.
