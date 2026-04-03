import 'dart:async';

import 'package:flutter/material.dart';

import '../../../../app/app_routes.dart';
import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../../auth/domain/entities/user.dart';
import '../../../game/presentation/pages/game_page.dart';
import '../../domain/entities/room.dart';
import '../../domain/entities/room_member.dart';
import '../../domain/repositories/lobby_repository.dart';
import '../state/room_waiting_controller.dart';

class RoomWaitingArgs {
  const RoomWaitingArgs({required this.currentUser, required this.roomId});

  final User currentUser;
  final String roomId;
}

class RoomWaitingPage extends StatefulWidget {
  const RoomWaitingPage({
    super.key,
    required this.lobbyRepository,
    required this.args,
  });

  final LobbyRepository lobbyRepository;
  final RoomWaitingArgs args;

  @override
  State<RoomWaitingPage> createState() => _RoomWaitingPageState();
}

class _RoomWaitingPageState extends State<RoomWaitingPage> {
  late final RoomWaitingController _controller;
  bool _didNavigateToGame = false;

  @override
  void initState() {
    super.initState();
    _controller = RoomWaitingController(
      lobbyRepository: widget.lobbyRepository,
      roomId: widget.args.roomId,
      currentPlayerId: widget.args.currentUser.id,
    );
    _controller.initialize();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _leaveRoom() async {
    final left = await _controller.leaveRoom();
    if (!mounted) {
      return;
    }

    if (left) {
      Navigator.of(context).pop();
      return;
    }

    final message =
        _controller.errorMessage ?? 'Nao foi possivel sair da sala.';
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }

  Future<void> _startGame() async {
    final started = await _controller.startGame();
    if (!mounted) {
      return;
    }
    if (started) {
      await _controller.refreshRoom();
      _maybeNavigateToGame();
      return;
    }
    final message =
        _controller.errorMessage ?? 'Nao foi possivel iniciar o jogo.';
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }

  void _maybeNavigateToGame() {
    if (!mounted || _didNavigateToGame || !_controller.hasGameStarted) {
      return;
    }

    final room = _controller.room;
    if (room == null) {
      return;
    }

    _didNavigateToGame = true;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) {
        return;
      }

      Navigator.of(context).pushReplacementNamed(
        AppRoutes.game,
        arguments: GamePageArgs(
          room: Room(
            id: room.id,
            name: room.name,
            playersCount: room.playersCount,
            maxPlayers: room.maxPlayers,
            status: room.status,
            isPrivate: room.isPrivate,
          ),
          currentPlayerId: widget.args.currentUser.id,
        ),
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    return PopScope(
      onPopInvokedWithResult: (didPop, result) {
        if (!didPop || _controller.isActionLoading) {
          return;
        }
        unawaited(_controller.leaveRoom());
      },
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Sala de Espera'),
          leading: IconButton(
            tooltip: 'Sair da sala',
            onPressed: _controller.isActionLoading ? null : _leaveRoom,
            icon: const Icon(Icons.arrow_back_rounded),
          ),
        ),
        body: AnimatedBuilder(
          animation: _controller,
          builder: (context, _) {
            _maybeNavigateToGame();

            if (_controller.isLoading && _controller.room == null) {
              return const TableBackground(
                child: Center(child: CircularProgressIndicator()),
              );
            }

            if (_controller.roomUnavailable) {
              return TableBackground(
                child: Center(
                  child: Padding(
                    padding: const EdgeInsets.all(20),
                    child: SectionCard(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Icon(Icons.info_outline_rounded, size: 36),
                          const SizedBox(height: 12),
                          Text(
                            _controller.errorMessage ?? 'A sala foi removida.',
                            textAlign: TextAlign.center,
                          ),
                          const SizedBox(height: 14),
                          ElevatedButton(
                            onPressed: () => Navigator.of(context).pop(),
                            child: const Text('Voltar ao lobby'),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              );
            }

            final room = _controller.room;
            if (room == null) {
              return const TableBackground(
                child: Center(child: Text('Nao foi possivel carregar a sala.')),
              );
            }

            final players = _orderPlayers(
              room.players,
              hostPlayerId: room.hostPlayerId,
            );

            return TableBackground(
              child: RefreshIndicator(
                onRefresh: _controller.refreshRoom,
                child: ListView(
                  padding: const EdgeInsets.fromLTRB(16, 14, 16, 24),
                  children: [
                    SectionCard(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            room.name,
                            style: Theme.of(context).textTheme.titleLarge,
                          ),
                          const SizedBox(height: 8),
                          Text('ID da sala: ${room.id}'),
                          const SizedBox(height: 6),
                          Text(
                            'Jogadores: ${room.playersCount}/${room.maxPlayers}',
                          ),
                          const SizedBox(height: 12),
                          LinearProgressIndicator(
                            minHeight: 8,
                            value: (room.playersCount / room.maxPlayers).clamp(
                              0.0,
                              1.0,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 12),
                    if (_controller.hasAllPlayers)
                      const SectionCard(
                        child: Row(
                          children: [
                            Icon(
                              Icons.check_circle_outline_rounded,
                              color: Color(0xFF155B42),
                            ),
                            SizedBox(width: 10),
                            Expanded(
                              child: Text(
                                'Sala completa (4/4). Todos os jogadores estao conectados.',
                              ),
                            ),
                          ],
                        ),
                      )
                    else
                      const SectionCard(
                        child: Row(
                          children: [
                            Icon(Icons.hourglass_bottom_rounded),
                            SizedBox(width: 10),
                            Expanded(
                              child: Text(
                                'A aguardar mais jogadores para completar a mesa.',
                              ),
                            ),
                          ],
                        ),
                      ),
                    const SizedBox(height: 12),
                    ...List.generate(4, (index) {
                      final player = index < players.length
                          ? players[index]
                          : null;
                      final seatLabel = 'Lugar ${index + 1}';
                      return Padding(
                        padding: const EdgeInsets.only(bottom: 10),
                        child: _SeatCard(
                          seatLabel: seatLabel,
                          player: player,
                          isHost: player?.id == room.hostPlayerId,
                          isCurrentUser:
                              player?.id == widget.args.currentUser.id,
                        ),
                      );
                    }),
                    const SizedBox(height: 8),
                    if (_controller.errorMessage != null)
                      Text(
                        _controller.errorMessage!,
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Colors.red.shade200,
                        ),
                      ),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: _controller.isActionLoading
                                ? null
                                : _leaveRoom,
                            icon: const Icon(Icons.logout_rounded),
                            label: const Text('Sair da sala'),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: ElevatedButton.icon(
                            onPressed: _controller.hasGameStarted
                                ? _maybeNavigateToGame
                                : _startGame,
                            icon: const Icon(Icons.play_arrow_rounded),
                            label: Text(
                              _controller.hasGameStarted
                                  ? 'Entrar no jogo'
                                  : (_controller.isHost
                                        ? 'Iniciar jogo'
                                        : 'Aguardar host'),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  List<RoomMember> _orderPlayers(
    List<RoomMember> players, {
    required String hostPlayerId,
  }) {
    final copy = List<RoomMember>.from(players);
    copy.sort((a, b) {
      if (a.id == hostPlayerId && b.id != hostPlayerId) {
        return -1;
      }
      if (b.id == hostPlayerId && a.id != hostPlayerId) {
        return 1;
      }
      return a.nickname.toLowerCase().compareTo(b.nickname.toLowerCase());
    });
    return copy;
  }
}

class _SeatCard extends StatelessWidget {
  const _SeatCard({
    required this.seatLabel,
    required this.player,
    required this.isHost,
    required this.isCurrentUser,
  });

  final String seatLabel;
  final RoomMember? player;
  final bool isHost;
  final bool isCurrentUser;

  @override
  Widget build(BuildContext context) {
    return SectionCard(
      child: Row(
        children: [
          CircleAvatar(
            radius: 20,
            child: Text(
              seatLabel.split(' ').last,
              style: Theme.of(context).textTheme.labelLarge,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: player == null
                ? Text(
                    '$seatLabel: A aguardar jogador...',
                    style: Theme.of(context).textTheme.bodyMedium,
                  )
                : Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        player!.nickname,
                        style: Theme.of(context).textTheme.titleSmall,
                      ),
                      const SizedBox(height: 2),
                      Text(
                        _buildRoleLabel(),
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ],
                  ),
          ),
        ],
      ),
    );
  }

  String _buildRoleLabel() {
    final roles = <String>[];
    if (isHost) {
      roles.add('Host');
    }
    if (isCurrentUser) {
      roles.add('Tu');
    }
    if (roles.isEmpty) {
      return 'Jogador ligado';
    }
    return roles.join(' - ');
  }
}
