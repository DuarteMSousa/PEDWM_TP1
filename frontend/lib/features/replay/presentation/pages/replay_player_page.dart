import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../../game/domain/entities/card.dart';
import '../../domain/entities/game_summary.dart';
import '../../domain/repositories/replay_repository.dart';
import '../state/replay_player_controller.dart';

class ReplayPlayerPage extends StatefulWidget {
  const ReplayPlayerPage({
    super.key,
    required this.game,
    required this.userId,
    required this.replayRepository,
  });

  final GameSummary game;
  final String userId;
  final ReplayRepository replayRepository;

  @override
  State<ReplayPlayerPage> createState() => _ReplayPlayerPageState();
}

class _ReplayPlayerPageState extends State<ReplayPlayerPage> {
  late final ReplayPlayerController _controller;

  @override
  void initState() {
    super.initState();
    _controller = ReplayPlayerController(
      replayRepository: widget.replayRepository,
      userId: widget.userId,
      initialGame: widget.game,
    );
    _controller.load();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Replay ${widget.game.id.length > 8 ? widget.game.id.substring(0, 8) : widget.game.id}',
        ),
      ),
      body: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          if (_controller.isLoading) {
            return const TableBackground(
              child: Center(
                child: CircularProgressIndicator(color: Color(0xFFF8F0DB)),
              ),
            );
          }

          if (_controller.errorMessage != null) {
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
                          onPressed: _controller.load,
                          child: const Text('Tentar novamente'),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            );
          }

          final frame = _controller.currentFrame;

          return TableBackground(
            child: SafeArea(
              child: Column(
                children: [
                  // Cabeçalho de Informação
                  Padding(
                    padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
                    child: SectionCard(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Expanded(
                                child: Wrap(
                                  spacing: 8,
                                  runSpacing: 8,
                                  children: [
                                    _InfoBadge(
                                      icon: Icons.movie_filter_outlined,
                                      label:
                                          'Evento ${_controller.currentIndex}/${_controller.totalEvents}',
                                    ),
                                    _InfoBadge(
                                      icon: Icons.flag_outlined,
                                      label: frame.phaseLabel,
                                    ),
                                    _InfoBadge(
                                      icon: Icons.style_outlined,
                                      label: frame.trumpSuit == null
                                          ? 'Trunfo por revelar'
                                          : 'Trunfo: ${_suitLabel(frame.trumpSuit!.name)} ${_suitSymbol(frame.trumpSuit!.name)}',
                                    ),
                                  ],
                                ),
                              ),
                              const SizedBox(width: 12),
                              _SpeedChipGroup(
                                current: _controller.playbackSpeed,
                                onSelected: _controller.setPlaybackSpeed,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                  // Tabuleiro Central
                  Expanded(
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(12, 10, 12, 0),
                      child: _ReplayBoard(frame: frame),
                    ),
                  ),
                  // Controlos da Timeline
                  Padding(
                    padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
                    child: SectionCard(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Row(
                          //   children: [
                          //     const Icon(Icons.timeline_rounded),
                          //     const SizedBox(width: 8),
                          //     Expanded(
                          //       child: Text(
                          //         'Timeline do replay',
                          //         style: Theme.of(context).textTheme.titleSmall,
                          //       ),
                          //     ),
                          //     Text(
                          //       _controller.isPlaying
                          //           ? '${_controller.playbackSpeed.toStringAsFixed(_controller.playbackSpeed.truncateToDouble() == _controller.playbackSpeed ? 0 : 1)}x'
                          //           : 'Pausado',
                          //       style: Theme.of(context).textTheme.labelMedium,
                          //     ),
                          //   ],
                          // ),
                          SliderTheme(
                            data: SliderTheme.of(context).copyWith(
                              activeTrackColor: const Color(0xFF155B42),
                              thumbColor: const Color(0xFFD7B46A),
                              inactiveTrackColor: const Color(0x336A4A2D),
                            ),
                            child: Slider(
                              value: _controller.currentIndex.toDouble(),
                              min: 0,
                              max: _controller.totalEvents.toDouble() == 0
                                  ? 1.0
                                  : _controller.totalEvents.toDouble(),
                              onChanged: (value) =>
                                  _controller.seekTo(value.toInt()),
                            ),
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              IconButton(
                                onPressed: _controller.isAtStart
                                    ? null
                                    : _controller.reset,
                                icon: const Icon(Icons.skip_previous_rounded),
                              ),
                              IconButton(
                                onPressed: _controller.isAtStart
                                    ? null
                                    : _controller.stepBackward,
                                icon: const Icon(Icons.fast_rewind_rounded),
                              ),
                              IconButton(
                                iconSize: 42,
                                onPressed: _controller.isPlaying
                                    ? _controller.pause
                                    : (_controller.isAtEnd
                                          ? null
                                          : _controller.play),
                                icon: Icon(
                                  _controller.isPlaying
                                      ? Icons.pause_circle_filled_rounded
                                      : Icons.play_circle_fill_rounded,
                                ),
                              ),
                              IconButton(
                                onPressed: _controller.isAtEnd
                                    ? null
                                    : _controller.stepForward,
                                icon: const Icon(Icons.fast_forward_rounded),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class _ReplayBoard extends StatelessWidget {
  const _ReplayBoard({required this.frame});
  final ReplayFrame frame;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Wrap(
          spacing: 10,
          runSpacing: 10,
          alignment: WrapAlignment.center,
          children: frame.teams
              .map((team) => _TeamScorePill(team: team))
              .toList(),
        ),
        const SizedBox(height: 12),
        Expanded(
          child: DecoratedBox(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(24),
              gradient: const LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [Color(0xFF6A4A2D), Color(0xFF4A301C)],
              ),
              boxShadow: const [
                BoxShadow(
                  color: Color(0x55000000),
                  blurRadius: 20,
                  offset: Offset(0, 12),
                ),
              ],
            ),
            child: Padding(
              padding: const EdgeInsets.all(10),
              child: DecoratedBox(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(18),
                  gradient: const RadialGradient(
                    colors: [Color(0xFF1B6A4E), Color(0xFF0F3D2E)],
                    radius: 0.95,
                  ),
                  border: Border.all(
                    color: const Color(0x58D7B46A),
                    width: 1.2,
                  ),
                ),
                child: Stack(
              children: [
                // Top
                Align(
                  alignment: Alignment.topCenter,
                  child: Padding(
                    padding: const EdgeInsets.only(top: 18),
                    child: _SeatCluster(
                      seat: frame.seats.length > 2 ? frame.seats[2] : null,
                      isCurrent:
                          frame.seats.length > 2 &&
                          frame.seats[2]?.id == frame.currentPlayerId,
                      isWinner:
                          frame.seats.length > 2 &&
                          frame.seats[2]?.id == frame.winnerId,
                      verticalFan: false,
                    ),
                  ),
                ),
                // Left
                Align(
                  alignment: Alignment.centerLeft,
                  child: Padding(
                    padding: const EdgeInsets.only(left: 14),
                    child: _SeatCluster(
                      seat: frame.seats.length > 1 ? frame.seats[1] : null,
                      isCurrent:
                          frame.seats.length > 1 &&
                          frame.seats[1]?.id == frame.currentPlayerId,
                      isWinner:
                          frame.seats.length > 1 &&
                          frame.seats[1]?.id == frame.winnerId,
                      verticalFan: true,
                    ),
                  ),
                ),
                // Right
                Align(
                  alignment: Alignment.centerRight,
                  child: Padding(
                    padding: const EdgeInsets.only(right: 14),
                    child: _SeatCluster(
                      seat: frame.seats.length > 3 ? frame.seats[3] : null,
                      isCurrent:
                          frame.seats.length > 3 &&
                          frame.seats[3]?.id == frame.currentPlayerId,
                      isWinner:
                          frame.seats.length > 3 &&
                          frame.seats[3]?.id == frame.winnerId,
                      verticalFan: true,
                      reverseFan: true,
                    ),
                  ),
                ),
                // Bottom
                Align(
                  alignment: Alignment.bottomCenter,
                  child: Padding(
                    padding: const EdgeInsets.only(bottom: 18),
                    child: _SeatCluster(
                      seat: frame.seats.isNotEmpty ? frame.seats[0] : null,
                      isCurrent:
                          frame.seats.isNotEmpty &&
                          frame.seats[0]?.id == frame.currentPlayerId,
                      isWinner:
                          frame.seats.isNotEmpty &&
                          frame.seats[0]?.id == frame.winnerId,
                      verticalFan: false,
                    ),
                  ),
                ),
                Center(child: _TrickCenter(frame: frame)),
              ],
            ),
          ),
            ),
          ),
        ),
      ],
    );
  }
}

class _SeatCluster extends StatelessWidget {
  const _SeatCluster({
    required this.seat,
    required this.isCurrent,
    required this.isWinner,
    required this.verticalFan,
    this.reverseFan = false,
  });

  final ReplaySeat? seat;
  final bool isCurrent;
  final bool isWinner;
  final bool verticalFan;
  final bool reverseFan;

  @override
  Widget build(BuildContext context) {
    if (seat == null) return const SizedBox(width: 118, height: 88);

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        AnimatedContainer(
          duration: const Duration(milliseconds: 220),
          width: 118,
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
          decoration: BoxDecoration(
            color: isCurrent
                ? const Color(0xFFF8F0DB)
                : const Color(0xE61A4D3A),
            borderRadius: BorderRadius.circular(18),
            border: Border.all(
              color: isWinner
                  ? const Color(0xFFD7B46A)
                  : (isCurrent
                        ? const Color(0xFF155B42)
                        : const Color(0x66F8F0DB)),
              width: isWinner || isCurrent ? 2 : 1,
            ),
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                seat!.name,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(
                  color: isCurrent
                      ? const Color(0xFF155B42)
                      : const Color(0xFFF8F0DB),
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                seat!.teamId ?? 'Sem equipa',
                style: TextStyle(
                  fontSize: 10,
                  color: isCurrent
                      ? const Color(0xFF155B42).withValues(alpha: 0.7)
                      : const Color(0xCCF8F0DB),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 6),
          _MiniCardFan(
            cards: seat!.handCards,
            vertical: verticalFan,
          )
      ],
    );
  }
}

class _TrickCenter extends StatelessWidget {
  const _TrickCenter({required this.frame});
  final ReplayFrame frame;

  @override
  Widget build(BuildContext context) {
    final pile = <_PileCard>[
      if (frame.seats.length > 2 &&
          _buildPileCard(frame.seats[2], const Offset(0, -30), 0, 'top') !=
              null)
        _buildPileCard(frame.seats[2], const Offset(0, -30), 0, 'top')!,
      if (frame.seats.length > 1 &&
          _buildPileCard(frame.seats[1], const Offset(-35, 0), -0.2, 'left') !=
              null)
        _buildPileCard(frame.seats[1], const Offset(-35, 0), -0.2, 'left')!,
      if (frame.seats.length > 3 &&
          _buildPileCard(frame.seats[3], const Offset(35, 0), 0.2, 'right') !=
              null)
        _buildPileCard(frame.seats[3], const Offset(35, 0), 0.2, 'right')!,
      if (frame.seats.isNotEmpty &&
          _buildPileCard(frame.seats[0], const Offset(0, 30), 0, 'bottom') !=
              null)
        _buildPileCard(frame.seats[0], const Offset(0, 30), 0, 'bottom')!,
    ];

    return SizedBox(
      width: 200,
      height: 200,
      child: Stack(
        children: [
          if (pile.isEmpty)
            const Center(
              child: Text(
                'Vaza vazia',
                style: TextStyle(color: Colors.white54),
              ),
            ),
          ...pile.map(
            (item) => Align(
              alignment: Alignment.center,
              child: Transform.translate(
                offset: item.offset,
                child: _TablePlayedCard(
                  card: item.card,
                  angle: item.angle,
                  seatTag: item.seatTag,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  _PileCard? _buildPileCard(
    ReplaySeat? seat,
    Offset offset,
    double angle,
    String seatTag,
  ) {
    if (seat == null) return null;
    final card = frame.tableCards[seat.id];
    if (card == null) return null;
    return _PileCard(
      seatTag: seatTag,
      card: card,
      offset: offset,
      angle: angle,
    );
  }
}

// Widgets de suporte e Helpers (Pills, Badges, Chips)
class _TeamScorePill extends StatelessWidget {
  const _TeamScorePill({required this.team});
  final ReplayTeamSnapshot team;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: const Color(0xEEF8F0DB),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          Text(
            team.id,
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 12),
          ),
          Text(
            '${team.score} - ${team.roundScore} pts',
            style: const TextStyle(fontSize: 11),
          ),
        ],
      ),
    );
  }
}

class _TablePlayedCard extends StatelessWidget {
  const _TablePlayedCard({
    required this.card,
    required this.angle,
    required this.seatTag,
  });
  final SuecaCard card;
  final double angle;
  final String seatTag;

  @override
  Widget build(BuildContext context) {
    return Transform.rotate(
      angle: angle,
      child: SizedBox(
        width: 60,
        height: 85,
        child: DecoratedBox(
          decoration: BoxDecoration(
            color: const Color(0xFFFFFAEE),
            borderRadius: BorderRadius.circular(8),
            border: Border.all(color: const Color(0xFFB08D49), width: 1),
            boxShadow: const [
              BoxShadow(
                color: Color(0x35000000),
                blurRadius: 8,
                offset: Offset(0, 5),
              ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(7),
            child: SvgPicture.asset(
              _cardFrontAssetPath(card),
              fit: BoxFit.cover,
            ),
          ),
        ),
      ),
    );
  }
}

class _MiniCardFan extends StatelessWidget {
  const _MiniCardFan({
    required this.cards,
    required this.vertical,
  });

  final List<SuecaCard> cards;
  final bool vertical;

  @override
  Widget build(BuildContext context) {
    final safeCards = cards.length > 10 ? cards.sublist(0, 10) : cards;
    if (safeCards.isEmpty) return const SizedBox.shrink();

    const cardWidth = 28.0;
    const cardHeight = 40.0;
    const spread = 7.0;

    Widget fan = SizedBox(
      width: cardWidth + ((safeCards.length - 1) * spread),
      height: cardHeight,
      child: Stack(
        children: List.generate(safeCards.length, (index) {
          final card = safeCards[index];
          return Positioned(
            left: index * spread,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(3),
              child: SizedBox(
                width: cardWidth,
                height: cardHeight,
                child: SvgPicture.asset(
                  _cardFrontAssetPath(card),
                  fit: BoxFit.cover,
                ),
              ),
            ),
          );
        }),
      ),
    );

    return vertical ? RotatedBox(quarterTurns: 1, child: fan) : fan;
  }
}

class _InfoBadge extends StatelessWidget {
  const _InfoBadge({required this.icon, required this.label});
  final IconData icon;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: const Color(0xFFF8F0DB),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 14),
          const SizedBox(width: 4),
          Text(label, style: const TextStyle(fontSize: 11)),
        ],
      ),
    );
  }
}

class _SpeedChipGroup extends StatelessWidget {
  const _SpeedChipGroup({required this.current, required this.onSelected});
  final double current;
  final ValueChanged<double> onSelected;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [1.0, 1.5, 2.0].map((s) {
        return Padding(
          padding: const EdgeInsets.only(left: 4),
          child: ChoiceChip(
            label: Text('${s}x'),
            selected: current == s,
            onSelected: (_) => onSelected(s),
          ),
        );
      }).toList(),
    );
  }
}

// Modelos de dados para organização do Trick (Vaza)
class _PileCard {
  const _PileCard({
    required this.seatTag,
    required this.card,
    required this.offset,
    required this.angle,
  });
  final String seatTag;
  final SuecaCard card;
  final Offset offset;
  final double angle;
}

// --- Métodos de Mapeamento e Labels ---

String _cardFrontAssetPath(SuecaCard card) {
  final rankToken = switch (card.rank) {
    1 => 'ace',
    11 => 'jack',
    12 => 'queen',
    13 => 'king',
    _ => card.rank.toString(),
  };
  return 'assets/cards/svg-cards/${rankToken}_of_${card.suit.name.toLowerCase()}.svg';
}

String _suitSymbol(String suit) {
  return switch (suit.toUpperCase()) {
    'HEARTS' => '♥',
    'DIAMONDS' => '♦',
    'CLUBS' => '♣',
    'SPADES' => '♠',
    _ => suit,
  };
}

String _suitLabel(String suit) {
  return switch (suit.toUpperCase()) {
    'HEARTS' => 'Copas',
    'DIAMONDS' => 'Ouros',
    'CLUBS' => 'Paus',
    'SPADES' => 'Espadas',
    _ => suit,
  };
}
