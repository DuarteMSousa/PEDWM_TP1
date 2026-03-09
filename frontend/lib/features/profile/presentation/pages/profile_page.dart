import 'package:flutter/material.dart';

import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
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
                  constraints: const BoxConstraints(maxWidth: 880),
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
                                .toList(),
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
