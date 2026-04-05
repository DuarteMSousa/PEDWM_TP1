import 'package:flutter/material.dart';

import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../domain/entities/game_summary.dart';
import '../../domain/repositories/replay_repository.dart';
import '../state/replay_list_controller.dart';
import 'replay_player_page.dart';

class ReplayListPage extends StatefulWidget {
  const ReplayListPage({
    super.key,
    required this.replayRepository,
    required this.userId,
  });

  final ReplayRepository replayRepository;
  final String userId;

  @override
  State<ReplayListPage> createState() => _ReplayListPageState();
}

class _ReplayListPageState extends State<ReplayListPage> {
  late final ReplayListController _controller;

  @override
  void initState() {
    super.initState();
    _controller = ReplayListController(
      replayRepository: widget.replayRepository,
      userId: widget.userId,
    );
    _controller.loadGames();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _openReplay(GameSummary game) {
    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => ReplayPlayerPage(game: game),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Historico de Jogos')),
      body: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          if (_controller.isLoading && _controller.games.isEmpty) {
            return const TableBackground(
              child: Center(
                child: CircularProgressIndicator(color: Color(0xFFF8F0DB)),
              ),
            );
          }

          if (_controller.errorMessage != null && _controller.games.isEmpty) {
            return TableBackground(
              child: Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: SectionCard(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(Icons.error_outline_rounded, size: 36),
                        const SizedBox(height: 12),
                        Text(
                          _controller.errorMessage!,
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 14),
                        ElevatedButton(
                          onPressed: _controller.loadGames,
                          child: const Text('Tentar novamente'),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            );
          }

          if (_controller.games.isEmpty) {
            return const TableBackground(
              child: Center(
                child: Padding(
                  padding: EdgeInsets.all(24),
                  child: SectionCard(
                    child: Text(
                      'Ainda nao tens jogos registados.',
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
              ),
            );
          }

          return TableBackground(
            child: ListView.builder(
              padding: const EdgeInsets.fromLTRB(16, 12, 16, 24),
              itemCount: _controller.games.length,
              itemBuilder: (context, index) {
                final game = _controller.games[index];
                return Padding(
                  padding: const EdgeInsets.only(bottom: 10),
                  child: _GameHistoryCard(
                    game: game,
                    onReplay: () => _openReplay(game),
                  ),
                );
              },
            ),
          );
        },
      ),
    );
  }
}

class _GameHistoryCard extends StatelessWidget {
  const _GameHistoryCard({required this.game, required this.onReplay});

  final GameSummary game;
  final VoidCallback onReplay;

  @override
  Widget build(BuildContext context) {
    final playerNames =
        game.players.map((p) => p.username).join(', ');
    final dateLabel =
        '${game.createdAt.day.toString().padLeft(2, '0')}/'
        '${game.createdAt.month.toString().padLeft(2, '0')}/'
        '${game.createdAt.year} '
        '${game.createdAt.hour.toString().padLeft(2, '0')}:'
        '${game.createdAt.minute.toString().padLeft(2, '0')}';
    final eventCount = game.events.length;

    final gameEndEvent = game.events
        .where((e) => e.type == 'GAME_ENDED')
        .toList();
    String? winnerLabel;
    if (gameEndEvent.isNotEmpty) {
      final payload = gameEndEvent.first.payload;
      if (payload != null) {
        winnerLabel = payload['winner']?.toString();
      }
    }

    return SectionCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.sports_esports_outlined,
                  color: Color(0xFF6A4A2D)),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  'Jogo ${game.id.length > 8 ? game.id.substring(0, 8) : game.id}',
                  style: Theme.of(context).textTheme.titleSmall,
                ),
              ),
            ],
          ),
          const SizedBox(height: 6),
          Text(dateLabel, style: Theme.of(context).textTheme.bodySmall),
          const SizedBox(height: 4),
          Text(
            'Jogadores: $playerNames',
            style: Theme.of(context).textTheme.bodySmall,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
          if (winnerLabel != null) ...[
            const SizedBox(height: 4),
            Text(
              'Vencedor: $winnerLabel',
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: const Color(0xFF155B42),
                  ),
            ),
          ],
          const SizedBox(height: 4),
          Text(
            '$eventCount eventos',
            style: Theme.of(context).textTheme.bodySmall,
          ),
          const SizedBox(height: 10),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: onReplay,
              icon: const Icon(Icons.replay_rounded, size: 18),
              label: const Text('Ver Replay'),
            ),
          ),
        ],
      ),
    );
  }
}
