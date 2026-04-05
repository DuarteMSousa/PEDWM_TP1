import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../../game/domain/entities/card.dart';
// Removido import de suit.dart se não for usado explicitamente como tipo
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
          final currentEvent = _controller.currentEvent;
          final playerNames = _buildPlayerNames(_controller.game);

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
                      padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
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
                  // Lista de Eventos Recentes
                  Padding(
                    padding: const EdgeInsets.fromLTRB(16, 12, 16, 24),
                    child: SectionCard(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Momentos recentes',
                            style: Theme.of(context).textTheme.titleSmall,
                          ),
                          const SizedBox(height: 10),
                          if (_controller.recentEvents.isEmpty)
                            const Text('Ainda não há eventos visíveis.')
                          else
                            SizedBox(
                              height: 50,
                              child: ListView.separated(
                                separatorBuilder: (_, __) =>
                                    const SizedBox(height: 8),
                                itemCount: _controller.recentEvents.length,
                                itemBuilder: (context, index) {
                                  final event = _controller.recentEvents[index];
                                  return _EventTile(
                                    event: event,
                                    isHighlighted: event.id == currentEvent?.id,
                                    playerNames: playerNames,
                                  );
                                },
                              ),
                            ),
                        ],
                      ),
                    ),
                  ),
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
          child: Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(28),
              gradient: const RadialGradient(
                colors: [Color(0xFF1B6A4E), Color(0xFF0F3D2E)],
                radius: 0.95,
              ),
              border: Border.all(color: const Color(0x66D7B46A), width: 1.5),
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
                      ? const Color(0xFF155B42).withOpacity(0.7)
                      : const Color(0xCCF8F0DB),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 8),
        _CardBackFan(
          count: seat!.handCount,
          vertical: verticalFan,
          reverse: reverseFan,
          cardBackAssetPath: 'assets/cards/back/back_blue.svg',
        ),
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

class _CardBackFan extends StatelessWidget {
  const _CardBackFan({
    required this.count,
    required this.vertical,
    required this.cardBackAssetPath,
    this.reverse = false,
  });

  final int count;
  final bool vertical;
  final bool reverse;
  final String cardBackAssetPath;

  @override
  Widget build(BuildContext context) {
    final safeCount = count.clamp(0, 10);
    if (safeCount == 0) return const SizedBox.shrink();

    const cardWidth = 24.0;
    const cardHeight = 36.0;
    const spread = 6.0;

    Widget fan = SizedBox(
      width: cardWidth + ((safeCount - 1) * spread),
      height: cardHeight,
      child: Stack(
        children: List.generate(safeCount, (index) {
          return Positioned(
            left: index * spread,
            child: SvgPicture.asset(
              cardBackAssetPath,
              width: cardWidth,
              height: cardHeight,
            ),
          );
        }),
      ),
    );

    return vertical ? RotatedBox(quarterTurns: 1, child: fan) : fan;
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
      child: SvgPicture.asset(_cardFrontAssetPath(card), width: 60, height: 85),
    );
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

class _EventTile extends StatelessWidget {
  const _EventTile({
    required this.event,
    required this.isHighlighted,
    required this.playerNames,
  });
  final GameEvent event;
  final bool isHighlighted;
  final Map<String, String> playerNames;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: isHighlighted
            ? const Color(0xFFD7B46A).withOpacity(0.2)
            : Colors.black12,
        borderRadius: BorderRadius.circular(8),
        border: isHighlighted
            ? Border.all(color: const Color(0xFFD7B46A))
            : null,
      ),
      child: Row(
        children: [
          Icon(_iconForType(event.type), size: 16),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  _labelForType(event.type),
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 12,
                  ),
                ),
                Text(
                  _descriptionForEvent(event, playerNames),
                  style: TextStyle(fontSize: 11, color: Colors.grey[800]),
                ),
              ],
            ),
          ),
        ],
      ),
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

Map<String, String> _buildPlayerNames(GameSummary game) {
  final names = <String, String>{};
  for (final p in game.players) {
    names[p.id] = p.username;
  }
  return names;
}

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

String _labelForType(String type) {
  return switch (type) {
    'PLAYER_JOINED' => 'Entrou',
    'GAME_STARTED' => 'Início',
    'CARD_PLAYED' => 'Jogada',
    'TRICK_ENDED' => 'Fim da Vaza',
    'ROUND_ENDED' => 'Fim da Ronda',
    _ => type.replaceAll('_', ' '),
  };
}

IconData _iconForType(String type) {
  return switch (type) {
    'CARD_PLAYED' => Icons.play_arrow,
    'TRICK_ENDED' => Icons.check_circle,
    'GAME_ENDED' => Icons.emoji_events,
    _ => Icons.info_outline,
  };
}

String _descriptionForEvent(GameEvent event, Map<String, String> names) {
  final payload = event.payload ?? {};
  String getName(String? id) => names[id] ?? id?.substring(0, 4) ?? '?';

  String describeCard(dynamic cardRaw) {
    if (cardRaw is! Map) {
      return 'uma carta';
    }
    final rank = cardRaw['rank']?.toString() ?? '';
    final suit = cardRaw['suit']?.toString() ?? '';
    if (rank.isEmpty && suit.isEmpty) {
      return 'uma carta';
    }
    return '${_rankLabel(rank)} de ${_suitLabel(suit)} ${_suitSymbol(suit)}';
  }

  return switch (event.type) {
    'PLAYER_JOINED' =>
      '${getName(payload['playerId']?.toString())} entrou na partida',
    'PLAYER_LEFT' =>
      '${getName(payload['playerId']?.toString())} saiu da partida',
    'GAME_STARTED' => 'Partida iniciada',
    'ROUND_STARTED' =>
      'Ronda ${payload['roundNumber']?.toString() ?? '?'} (dealer: ${getName(payload['dealerId']?.toString())})',
    'TRICK_STARTED' =>
      'Nova vaza - lider ${getName(payload['leaderId']?.toString())}',
    'TURN_CHANGED' =>
      'Vez de ${getName(payload['playerId']?.toString())}',
    'TRUMP_REVEALED' =>
      'Trunfo: ${describeCard(payload['card'])}',
    'CARD_DEALT' =>
      '${getName(payload['playerId']?.toString())} recebeu ${describeCard(payload['card'])}',
    'CARD_PLAYED' =>
      '${getName(payload['playerId']?.toString())} jogou ${describeCard(payload['card'])}',
    'TRICK_ENDED' =>
      'Vencedor: ${getName(payload['winnerId']?.toString())} (+${payload['points']?.toString() ?? '0'} pts)',
    'ROUND_ENDED' =>
      'Ronda terminada - vencedor ${payload['winnerTeam']?.toString() ?? '?'}',
    'GAME_SCORE_UPDATED' => 'Pontuação atualizada',
    'GAME_ENDED' => 'Jogo terminado - vencedor ${payload['winner']?.toString() ?? '?'}',
    _ => '',
  };
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

String _rankLabel(String rank) {
  return switch (rank.toUpperCase()) {
    '1' || 'A' => 'Ás',
    '11' || 'J' => 'Valete',
    '12' || 'Q' => 'Dama',
    '13' || 'K' => 'Rei',
    '7' => 'Bisca',
    _ => rank,
  };
}
