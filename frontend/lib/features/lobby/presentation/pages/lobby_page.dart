import 'package:flutter/material.dart';

import '../../../../app/app_routes.dart';
import '../../../../core/shared_widgets/motion.dart';
import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../../auth/domain/entities/user.dart';
import '../../domain/entities/room.dart';
import '../../domain/repositories/lobby_repository.dart';
import 'room_waiting_page.dart';
import '../state/lobby_controller.dart';

class LobbyPage extends StatefulWidget {
  const LobbyPage({
    super.key,
    required this.lobbyRepository,
    required this.currentUser,
  });

  final LobbyRepository lobbyRepository;
  final User currentUser;

  @override
  State<LobbyPage> createState() => _LobbyPageState();
}

class _LobbyPageState extends State<LobbyPage> {
  late final LobbyController _controller;

  @override
  void initState() {
    super.initState();
    _controller = LobbyController(lobbyRepository: widget.lobbyRepository);
    _controller.loadRooms(playerId: widget.currentUser.id);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _openProfile() {
    Navigator.of(
      context,
    ).pushNamed(AppRoutes.profile, arguments: widget.currentUser.id);
  }

  Future<void> _refreshRooms() {
    return _controller.refreshRooms(playerId: widget.currentUser.id);
  }

  void _openRoomWaiting({required String roomId}) {
    Navigator.of(context).pushNamed(
      AppRoutes.roomWaiting,
      arguments: RoomWaitingArgs(
        currentUser: widget.currentUser,
        roomId: roomId,
      ),
    );
  }

  Future<void> _joinRoom(Room room) async {
    final joinedRoom = await _controller.joinRoom(
      roomId: room.id,
      playerId: widget.currentUser.id,
    );

    if (!mounted) {
      return;
    }
    if (joinedRoom != null) {
      _openRoomWaiting(roomId: joinedRoom.id);
      return;
    }

    final message =
        _controller.errorMessage ?? 'Nao foi possivel entrar na sala.';
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }

  Future<void> _createRoom() async {
    final createdRoom = await _controller.createRoom(
      hostPlayerId: widget.currentUser.id,
    );

    if (!mounted) {
      return;
    }
    if (createdRoom != null) {
      _openRoomWaiting(roomId: createdRoom.id);
      return;
    }

    final message =
        _controller.errorMessage ?? 'Nao foi possivel criar a sala.';
    ScaffoldMessenger.of(
      context,
    ).showSnackBar(SnackBar(content: Text(message)));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Lobby de Sueca'),
        actions: [
          IconButton(
            onPressed: _createRoom,
            icon: const Icon(Icons.add_circle_outline_rounded),
            tooltip: 'Criar sala',
          ),
          IconButton(
            onPressed: _openProfile,
            icon: const Icon(Icons.account_circle_outlined),
            tooltip: 'Perfil',
          ),
        ],
      ),
      body: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          if (_controller.isLoading && _controller.rooms.isEmpty) {
            return const TableBackground(
              child: Center(
                child: CircularProgressIndicator(color: Color(0xFFF8F0DB)),
              ),
            );
          }

          if (_controller.errorMessage != null && _controller.rooms.isEmpty) {
            return TableBackground(
              child: Center(
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: 420),
                  child: Padding(
                    padding: const EdgeInsets.all(20),
                    child: SectionCard(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Icon(Icons.wifi_off_rounded, size: 40),
                          const SizedBox(height: 12),
                          Text(
                            _controller.errorMessage!,
                            textAlign: TextAlign.center,
                          ),
                          const SizedBox(height: 14),
                          ElevatedButton(
                            onPressed: _refreshRooms,
                            child: const Text('Tentar novamente'),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            );
          }

          return TableBackground(
            child: RefreshIndicator(
              onRefresh: _refreshRooms,
              color: const Color(0xFF155B42),
              backgroundColor: const Color(0xFFF8F0DB),
              child: ListView(
                padding: const EdgeInsets.fromLTRB(16, 12, 16, 24),
                children: [
                  RevealSlideFade(
                    delay: const Duration(milliseconds: 60),
                    beginOffset: const Offset(0, 0.035),
                    child: _LobbyHeader(
                      user: widget.currentUser,
                      roomCount: _controller.rooms.length,
                    ),
                  ),
                  const SizedBox(height: 14),
                  LayoutBuilder(
                    builder: (context, constraints) {
                      if (_controller.rooms.isEmpty) {
                        return const SectionCard(
                          child: Text(
                            'Sem salas disponiveis. Atualiza para procurar mesas ativas.',
                          ),
                        );
                      }

                      final maxWidth = constraints.maxWidth;
                      final columnCount = maxWidth > 980
                          ? 3
                          : (maxWidth > 660 ? 2 : 1);
                      final cardWidth =
                          (maxWidth - ((columnCount - 1) * 12)) / columnCount;

                      return Wrap(
                        spacing: 12,
                        runSpacing: 12,
                        children: _controller.rooms.asMap().entries.map((
                          entry,
                        ) {
                          final index = entry.key;
                          final room = entry.value;
                          return SizedBox(
                            width: cardWidth,
                            child: RevealSlideFade(
                              delay: Duration(milliseconds: 120 + (index * 70)),
                              beginOffset: const Offset(0, 0.06),
                              child: _RoomCard(
                                room: room,
                                onJoin: room.isFull
                                    ? null
                                    : () => _joinRoom(room),
                              ),
                            ),
                          );
                        }).toList(),
                      );
                    },
                  ),
                ],
              ),
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _refreshRooms,
        child: const Icon(Icons.refresh_rounded),
      ),
    );
  }
}

class _LobbyHeader extends StatelessWidget {
  const _LobbyHeader({required this.user, required this.roomCount});

  final User user;
  final int roomCount;

  @override
  Widget build(BuildContext context) {
    final subtitleStyle = Theme.of(
      context,
    ).textTheme.bodyMedium?.copyWith(color: const Color(0xB9F8F0DB));
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(22),
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0x6B224C3A), Color(0x811D362B)],
        ),
        border: Border.all(color: const Color(0x55D7B46A)),
      ),
      child: Padding(
        padding: const EdgeInsets.all(18),
        child: Row(
          children: [
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                color: const Color(0x27E6C57C),
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: const Color(0x91D7B46A)),
              ),
              child: const Icon(
                Icons.casino_outlined,
                color: Color(0xFFF8F0DB),
                size: 30,
              ),
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Bem-vindo, ${user.nickname}',
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      color: const Color(0xFFF8F0DB),
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '$roomCount salas disponiveis para jogar agora',
                    style: subtitleStyle,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _RoomCard extends StatelessWidget {
  const _RoomCard({required this.room, required this.onJoin});

  final Room room;
  final VoidCallback? onJoin;

  @override
  Widget build(BuildContext context) {
    final occupancy = room.maxPlayers == 0
        ? 0.0
        : room.playersCount / room.maxPlayers;
    return HoverLift(
      child: SectionCard(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    room.name,
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text('Jogadores: ${room.occupancyLabel}'),
            const SizedBox(height: 8),
            ClipRRect(
              borderRadius: BorderRadius.circular(999),
              child: LinearProgressIndicator(
                minHeight: 8,
                value: occupancy.clamp(0.0, 1.0).toDouble(),
                backgroundColor: const Color(0x1A6A4A2D),
                valueColor: const AlwaysStoppedAnimation<Color>(
                  Color(0xFF155B42),
                ),
              ),
            ),
            const SizedBox(height: 14),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: onJoin,
                child: Text(room.isFull ? 'Sala cheia' : 'Entrar na mesa'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
