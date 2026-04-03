import 'dart:convert';

import 'package:flutter/material.dart';

import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../domain/entities/match_history_entry.dart';
import '../../domain/entities/match_replay.dart';
import '../../domain/repositories/profile_repository.dart';
import '../state/profile_controller.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({
    super.key,
    required this.profileRepository,
    required this.userId,
  });

  final ProfileRepository profileRepository;
  final String userId;

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  late final ProfileController _controller;
  String? _loadingReplayGameId;

  @override
  void initState() {
    super.initState();
    _controller = ProfileController(
      profileRepository: widget.profileRepository,
      userId: widget.userId,
    );
    _controller.loadProfile();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _openReplay(MatchHistoryEntry entry) async {
    setState(() => _loadingReplayGameId = entry.gameId);
    final replay = await _controller.loadReplay(entry.gameId);
    if (!mounted) {
      return;
    }
    setState(() => _loadingReplayGameId = null);

    if (replay == null) {
      final message = _controller.errorMessage ?? 'Replay indisponivel.';
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text(message)));
      return;
    }

    await showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      backgroundColor: const Color(0xFF0F4B37),
      useSafeArea: true,
      builder: (_) => _ReplaySheet(replay: replay),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Perfil')),
      body: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          if (_controller.isLoading && _controller.profile == null) {
            return const TableBackground(
              child: Center(
                child: CircularProgressIndicator(color: Color(0xFFF8F0DB)),
              ),
            );
          }

          if (_controller.errorMessage != null && _controller.profile == null) {
            return TableBackground(
              child: Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: SectionCard(
                    child: Text(
                      _controller.errorMessage!,
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
              ),
            );
          }

          final profile = _controller.profile;
          if (profile == null) {
            return const TableBackground(
              child: Center(
                child: Text(
                  'Perfil indisponivel.',
                  style: TextStyle(color: Color(0xFFF8F0DB)),
                ),
              ),
            );
          }

          return TableBackground(
            child: SafeArea(
              child: Center(
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: 920),
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      SectionCard(
                        child: Row(
                          children: [
                            Container(
                              width: 68,
                              height: 68,
                              decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(20),
                                gradient: const LinearGradient(
                                  colors: [
                                    Color(0xFFD7B46A),
                                    Color(0xFFB08D49),
                                  ],
                                ),
                              ),
                              child: const Icon(
                                Icons.person_rounded,
                                size: 36,
                                color: Color(0xFF2A241A),
                              ),
                            ),
                            const SizedBox(width: 14),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    profile.nickname,
                                    style: Theme.of(
                                      context,
                                    ).textTheme.headlineSmall,
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    'Jogador: ${profile.userId}',
                                    style: Theme.of(
                                      context,
                                    ).textTheme.bodyMedium,
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 12),
                      LayoutBuilder(
                        builder: (context, constraints) {
                          final width = constraints.maxWidth;
                          final cardsPerRow = width > 760 ? 3 : 1;
                          final cardWidth =
                              (width - ((cardsPerRow - 1) * 10)) / cardsPerRow;

                          final cards = [
                            _StatCard(
                              label: 'Partidas',
                              value: '${profile.matchesPlayed}',
                              icon: Icons.sports_esports_outlined,
                            ),
                            _StatCard(
                              label: 'Vitorias',
                              value: '${profile.wins}',
                              icon: Icons.emoji_events_outlined,
                            ),
                            _StatCard(
                              label: 'Taxa',
                              value: '${profile.winRate.toStringAsFixed(1)}%',
                              icon: Icons.trending_up_rounded,
                            ),
                          ];

                          return Wrap(
                            spacing: 10,
                            runSpacing: 10,
                            children: cards
                                .map(
                                  (card) =>
                                      SizedBox(width: cardWidth, child: card),
                                )
                                .toList(growable: false),
                          );
                        },
                      ),
                      const SizedBox(height: 12),
                      SectionCard(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Eficiencia',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            const SizedBox(height: 10),
                            ClipRRect(
                              borderRadius: BorderRadius.circular(999),
                              child: LinearProgressIndicator(
                                minHeight: 10,
                                value: (profile.winRate / 100).clamp(0.0, 1.0),
                                backgroundColor: const Color(0x1A6A4A2D),
                                valueColor: const AlwaysStoppedAnimation<Color>(
                                  Color(0xFF155B42),
                                ),
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              '${profile.winRate.toStringAsFixed(1)}% de partidas vencidas',
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 12),
                      SectionCard(
                        child: Row(
                          children: [
                            Text(
                              'Historico de partidas',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            const Spacer(),
                            if (_controller.isReplayLoading)
                              const SizedBox(
                                width: 16,
                                height: 16,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                ),
                              ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 8),
                      if (profile.matchHistory.isEmpty)
                        const SectionCard(
                          child: Text(
                            'Ainda nao tens partidas terminadas para mostrar.',
                          ),
                        )
                      else
                        ...profile.matchHistory.map((entry) {
                          final isLoading =
                              _loadingReplayGameId == entry.gameId;
                          return Padding(
                            padding: const EdgeInsets.only(bottom: 10),
                            child: _MatchHistoryCard(
                              entry: entry,
                              isReplayLoading: isLoading,
                              onReplay: isLoading
                                  ? null
                                  : () => _openReplay(entry),
                            ),
                          );
                        }),
                    ],
                  ),
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  const _StatCard({
    required this.label,
    required this.value,
    required this.icon,
  });

  final String label;
  final String value;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    return SectionCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: const Color(0xFF6A4A2D)),
          const SizedBox(height: 6),
          Text(
            value,
            style: Theme.of(
              context,
            ).textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.w800),
          ),
          const SizedBox(height: 2),
          Text(label),
        ],
      ),
    );
  }
}

class _MatchHistoryCard extends StatelessWidget {
  const _MatchHistoryCard({
    required this.entry,
    required this.isReplayLoading,
    required this.onReplay,
  });

  final MatchHistoryEntry entry;
  final bool isReplayLoading;
  final VoidCallback? onReplay;

  @override
  Widget build(BuildContext context) {
    final won = entry.won;
    final tone = won ? const Color(0xFF155B42) : const Color(0xFF8A302C);
    final resultText = won ? 'Vitoria' : 'Derrota';
    final playedAt = _formatDateTime(entry.playedAt);

    return SectionCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 10,
                  vertical: 6,
                ),
                decoration: BoxDecoration(
                  color: tone.withAlpha(26),
                  borderRadius: BorderRadius.circular(999),
                  border: Border.all(color: tone.withAlpha(120)),
                ),
                child: Text(
                  resultText,
                  style: Theme.of(context).textTheme.labelLarge?.copyWith(
                    color: tone,
                    fontWeight: FontWeight.w700,
                  ),
                ),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  'Jogo ${entry.gameId.length > 8 ? entry.gameId.substring(0, 8) : entry.gameId}',
                  style: Theme.of(
                    context,
                  ).textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w700),
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text('Score: ${entry.myScore} - ${entry.opponentScore}'),
          const SizedBox(height: 2),
          Text('Sala: ${entry.roomId}'),
          const SizedBox(height: 2),
          Text('Data: $playedAt'),
          const SizedBox(height: 10),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: onReplay,
              icon: isReplayLoading
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.replay_rounded),
              label: const Text('Ver replay'),
            ),
          ),
        ],
      ),
    );
  }
}

class _ReplaySheet extends StatelessWidget {
  const _ReplaySheet({required this.replay});

  final MatchReplay replay;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 20),
      child: SizedBox(
        height: MediaQuery.of(context).size.height * 0.82,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Replay',
              style: Theme.of(context).textTheme.titleLarge?.copyWith(
                color: const Color(0xFFF8F0DB),
                fontWeight: FontWeight.w800,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              'Jogo: ${replay.gameId}',
              style: Theme.of(
                context,
              ).textTheme.bodyMedium?.copyWith(color: const Color(0xFFF8F0DB)),
            ),
            Text(
              'Sala: ${replay.roomId}',
              style: Theme.of(
                context,
              ).textTheme.bodyMedium?.copyWith(color: const Color(0xFFF8F0DB)),
            ),
            const SizedBox(height: 10),
            Expanded(
              child: replay.events.isEmpty
                  ? Center(
                      child: Text(
                        'Sem eventos guardados para este jogo.',
                        style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: const Color(0xFFF8F0DB),
                        ),
                      ),
                    )
                  : ListView.separated(
                      itemCount: replay.events.length,
                      separatorBuilder: (context, index) =>
                          const SizedBox(height: 8),
                      itemBuilder: (context, index) {
                        final event = replay.events[index];
                        return Container(
                          padding: const EdgeInsets.all(10),
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(12),
                            color: const Color(0x3318352B),
                            border: Border.all(color: const Color(0x4CD7B46A)),
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Text(
                                    '#${event.sequence}',
                                    style: Theme.of(context)
                                        .textTheme
                                        .labelLarge
                                        ?.copyWith(
                                          color: const Color(0xFFF8F0DB),
                                          fontWeight: FontWeight.w700,
                                        ),
                                  ),
                                  const SizedBox(width: 8),
                                  Expanded(
                                    child: Text(
                                      event.type,
                                      style: Theme.of(context)
                                          .textTheme
                                          .labelLarge
                                          ?.copyWith(
                                            color: const Color(0xFFF8F0DB),
                                            fontWeight: FontWeight.w600,
                                          ),
                                    ),
                                  ),
                                  Text(
                                    _formatDateTime(event.timestamp),
                                    style: Theme.of(context).textTheme.bodySmall
                                        ?.copyWith(
                                          color: const Color(0xE6F8F0DB),
                                        ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 6),
                              SelectableText(
                                _prettyPayload(event.payload),
                                style: Theme.of(context).textTheme.bodySmall
                                    ?.copyWith(
                                      color: const Color(0xFFF8F0DB),
                                      fontFamily: 'monospace',
                                    ),
                              ),
                            ],
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }
}

String _formatDateTime(DateTime value) {
  final local = value.toLocal();
  String two(int n) => n.toString().padLeft(2, '0');
  return '${two(local.day)}/${two(local.month)}/${local.year} ${two(local.hour)}:${two(local.minute)}';
}

String _prettyPayload(String payload) {
  final trimmed = payload.trim();
  if (trimmed.isEmpty) {
    return '{}';
  }

  try {
    final decoded = jsonDecode(trimmed);
    final formatted = const JsonEncoder.withIndent('  ').convert(decoded);
    if (formatted.length > 1200) {
      return '${formatted.substring(0, 1200)}...';
    }
    return formatted;
  } catch (_) {
    if (trimmed.length > 1200) {
      return '${trimmed.substring(0, 1200)}...';
    }
    return trimmed;
  }
}
