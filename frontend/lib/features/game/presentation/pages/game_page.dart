import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../../../../app/app_routes.dart';
import '../../../../core/shared_widgets/motion.dart';
import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../../auth/domain/entities/user.dart';
import '../../../lobby/domain/entities/room.dart';
import '../../../lobby/presentation/pages/lobby_page.dart';
import '../../domain/entities/card.dart';
import '../../domain/entities/game_phase.dart';
import '../../domain/entities/player.dart';
import '../../domain/entities/sueca_game_state.dart';
import '../../domain/entities/suit.dart';
import '../../domain/repositories/game_repository.dart';
import '../state/game_controller.dart';

class GamePageArgs {
  const GamePageArgs({required this.room, required this.currentPlayerId});

  final Room room;
  final String currentPlayerId;
}

enum _PostGameAction { rematch, lobby }

class GamePage extends StatefulWidget {
  const GamePage({super.key, required this.gameRepository, required this.args});

  final GameRepository gameRepository;
  final GamePageArgs args;

  @override
  State<GamePage> createState() => _GamePageState();
}

class _GamePageState extends State<GamePage> {
  static const List<String> _cardBackCandidates = <String>[
    'assets/cards/back/back.svg',
    'assets/cards/back/back_blue.svg',
    'assets/cards/back/back_red.svg',
    'assets/cards/back/back.png',
    'assets/cards/back/back_blue.png',
    'assets/cards/back/back_red.png',
  ];

  late final GameController _controller;
  String? _cardBackAssetPath;
  bool _didHandleMatchEnd = false;

  @override
  void initState() {
    super.initState();
    _controller = GameController(
      gameRepository: widget.gameRepository,
      roomId: widget.args.room.id,
      currentPlayerId: widget.args.currentPlayerId,
    );
    _controller.initialize();
    _resolveCardBackAsset();
  }

  Future<void> _resolveCardBackAsset() async {
    for (final candidate in _cardBackCandidates) {
      try {
        await rootBundle.load(candidate);
        if (!mounted) {
          return;
        }
        setState(() => _cardBackAssetPath = candidate);
        return;
      } catch (_) {
        // Continue searching for an existing back asset.
      }
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _onMatchFinished(SuecaGameState state) {
    if (_didHandleMatchEnd) {
      return;
    }
    _didHandleMatchEnd = true;

    WidgetsBinding.instance.addPostFrameCallback((_) async {
      if (!mounted) {
        return;
      }

      final action = await _showPostGameSheet(state);
      if (!mounted) {
        return;
      }

      final user = _currentUserFromState(state);
      switch (action ?? _PostGameAction.lobby) {
        case _PostGameAction.rematch:
          Navigator.of(context).pushNamedAndRemoveUntil(
            AppRoutes.lobby,
            (_) => false,
            arguments: LobbyPageArgs(currentUser: user, autoCreateRoom: true),
          );
          break;
        case _PostGameAction.lobby:
          Navigator.of(context).pushNamedAndRemoveUntil(
            AppRoutes.lobby,
            (_) => false,
            arguments: user,
          );
          break;
      }
    });
  }

  Future<_PostGameAction?> _showPostGameSheet(SuecaGameState state) {
    final winner = _winnerLabel(state);

    return showModalBottomSheet<_PostGameAction>(
      context: context,
      isDismissible: false,
      enableDrag: false,
      useSafeArea: true,
      backgroundColor: const Color(0xFF0F4B37),
      builder: (context) {
        return Padding(
          padding: const EdgeInsets.fromLTRB(20, 18, 20, 22),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Partida terminada',
                style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  color: const Color(0xFFF8F0DB),
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(height: 8),
              Text(
                winner == null ? 'Empate técnico.' : 'Vencedor: $winner.',
                style: Theme.of(
                  context,
                ).textTheme.bodyLarge?.copyWith(color: const Color(0xFFF8F0DB)),
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: _ScorePill(
                      label: 'Equipa A',
                      score: state.teamAScore,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: _ScorePill(
                      label: 'Equipa B',
                      score: state.teamBScore,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () =>
                          Navigator.of(context).pop(_PostGameAction.lobby),
                      icon: const Icon(Icons.home_outlined),
                      label: const Text('Lobby'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: ElevatedButton.icon(
                      onPressed: () =>
                          Navigator.of(context).pop(_PostGameAction.rematch),
                      icon: const Icon(Icons.replay_rounded),
                      label: const Text('Revanche'),
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      },
    );
  }

  User _currentUserFromState(SuecaGameState state) {
    Player? me;
    for (final player in state.players) {
      if (player.id == state.myPlayerId) {
        me = player;
        break;
      }
    }
    return User(id: state.myPlayerId, nickname: me?.nickname ?? 'Jogador');
  }

  String? _winnerLabel(SuecaGameState state) {
    if (state.teamAScore == state.teamBScore) {
      return null;
    }
    return state.teamAScore > state.teamBScore ? 'Equipa A' : 'Equipa B';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Mesa ${widget.args.room.name}')),
      body: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          if (_controller.isLoading && _controller.gameState == null) {
            return const TableBackground(
              child: Center(
                child: CircularProgressIndicator(color: Color(0xFFF8F0DB)),
              ),
            );
          }

          final state = _controller.gameState;
          if (state == null) {
            return TableBackground(
              child: Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: SectionCard(
                    child: Text(
                      _controller.errorMessage ??
                          'Estado de jogo indisponivel.',
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
              ),
            );
          }

          if (state.phase == GamePhase.finished) {
            _onMatchFinished(state);
          }

          final players = _orderedPlayers(state);
          final me = _playerAt(players, 0);
          final left = _playerAt(players, 1);
          final top = _playerAt(players, 2);
          final right = _playerAt(players, 3);

          return TableBackground(
            child: SafeArea(
              child: Center(
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: 1120),
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(12, 10, 12, 14),
                    child: Column(
                      children: [
                        RevealSlideFade(
                          delay: const Duration(milliseconds: 50),
                          beginOffset: const Offset(0, 0.03),
                          child: _TopHud(
                            room: widget.args.room,
                            state: state,
                            isBusy: _controller.isLoading,
                          ),
                        ),
                        const SizedBox(height: 12),
                        Expanded(
                          child: LayoutBuilder(
                            builder: (context, constraints) {
                              final handHeight = constraints.maxHeight < 700
                                  ? 154.0
                                  : 180.0;
                              final gap = 10.0;
                              final boardSize = math.max(
                                250.0,
                                math.min(
                                  constraints.maxWidth,
                                  constraints.maxHeight - handHeight - gap,
                                ),
                              );

                              return Column(
                                children: [
                                  SizedBox(
                                    width: boardSize,
                                    height: boardSize,
                                    child: RevealSlideFade(
                                      delay: const Duration(milliseconds: 120),
                                      beginOffset: const Offset(0, 0.04),
                                      child: _WoodenTable(
                                        state: state,
                                        top: top,
                                        left: left,
                                        right: right,
                                        me: me,
                                        cardBackAssetPath: _cardBackAssetPath,
                                      ),
                                    ),
                                  ),
                                  const SizedBox(height: 10),
                                  SizedBox(
                                    height: handHeight,
                                    child: RevealSlideFade(
                                      delay: const Duration(milliseconds: 200),
                                      beginOffset: const Offset(0, 0.06),
                                      child: _MyHandArea(
                                        me: me,
                                        hand: state.hand,
                                        isBusy: _controller.isLoading,
                                        canPlay: _controller.canPlayCard,
                                        onPlayCard: _controller.playCard,
                                      ),
                                    ),
                                  ),
                                ],
                              );
                            },
                          ),
                        ),
                      ],
                    ),
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

class _TopHud extends StatelessWidget {
  const _TopHud({
    required this.room,
    required this.state,
    required this.isBusy,
  });

  final Room room;
  final SuecaGameState state;
  final bool isBusy;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(16),
        color: const Color(0x3A10291F),
        border: Border.all(color: const Color(0x6ED7B46A)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
        child: Wrap(
          spacing: 8,
          runSpacing: 8,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            _HudPill(
              label: room.name,
              icon: Icons.casino_outlined,
              strong: true,
            ),
            _HudPill(
              label: _phaseLabel(state.phase),
              icon: Icons.flag_outlined,
            ),
            _HudPill(
              label: 'Trunfo ${_suitLabel(state.trumpSuit)}',
              icon: Icons.style_outlined,
            ),
            _ScorePill(label: 'Equipa A', score: state.teamAScore),
            _ScorePill(label: 'Equipa B', score: state.teamBScore),
            if (isBusy)
              const SizedBox(
                width: 16,
                height: 16,
                child: CircularProgressIndicator(strokeWidth: 2),
              ),
          ],
        ),
      ),
    );
  }
}

class _HudPill extends StatelessWidget {
  const _HudPill({
    required this.label,
    required this.icon,
    this.strong = false,
  });

  final String label;
  final IconData icon;
  final bool strong;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        color: strong ? const Color(0x52D7B46A) : const Color(0x2EEAD8A8),
        border: Border.all(color: const Color(0x8ED7B46A)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 16, color: const Color(0xFFF8F0DB)),
            const SizedBox(width: 6),
            Text(
              label,
              style: Theme.of(context).textTheme.labelLarge?.copyWith(
                color: const Color(0xFFF8F0DB),
                fontWeight: strong ? FontWeight.w700 : FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ScorePill extends StatelessWidget {
  const _ScorePill({required this.label, required this.score});

  final String label;
  final int score;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        color: const Color(0xFFFAF0D8),
        border: Border.all(color: const Color(0xFFB08D49)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        child: Text(
          '$label: $score',
          style: Theme.of(
            context,
          ).textTheme.labelLarge?.copyWith(fontWeight: FontWeight.w700),
        ),
      ),
    );
  }
}

class _WoodenTable extends StatelessWidget {
  const _WoodenTable({
    required this.state,
    required this.top,
    required this.left,
    required this.right,
    required this.me,
    required this.cardBackAssetPath,
  });

  final SuecaGameState state;
  final Player? top;
  final Player? left;
  final Player? right;
  final Player? me;
  final String? cardBackAssetPath;

  @override
  Widget build(BuildContext context) {
    final opponentCards = math.max(5, state.hand.length + 1);
    return DecoratedBox(
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
        padding: const EdgeInsets.all(12),
        child: DecoratedBox(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(18),
            gradient: const LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [Color(0xFF176348), Color(0xFF0F4B37)],
            ),
            border: Border.all(color: const Color(0x58D7B46A), width: 1.2),
          ),
          child: Stack(
            clipBehavior: Clip.none,
            children: [
              Positioned.fill(
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(18),
                    gradient: RadialGradient(
                      center: const Alignment(0, -0.15),
                      radius: 1.1,
                      colors: [
                        const Color(0x4FF4F1D3),
                        const Color(0x00155342),
                      ],
                    ),
                  ),
                ),
              ),
              if (top != null)
                Align(
                  alignment: const Alignment(0, -1.06),
                  child: _OpponentSeat(
                    player: top!,
                    seat: _Seat.top,
                    isCurrent: top!.id == state.currentPlayerId,
                    cardCount: opponentCards,
                    cardBackAssetPath: cardBackAssetPath,
                  ),
                ),
              if (left != null)
                Align(
                  alignment: const Alignment(-1.02, 0),
                  child: _OpponentSeat(
                    player: left!,
                    seat: _Seat.left,
                    isCurrent: left!.id == state.currentPlayerId,
                    cardCount: opponentCards,
                    cardBackAssetPath: cardBackAssetPath,
                  ),
                ),
              if (right != null)
                Align(
                  alignment: const Alignment(1.02, 0),
                  child: _OpponentSeat(
                    player: right!,
                    seat: _Seat.right,
                    isCurrent: right!.id == state.currentPlayerId,
                    cardCount: opponentCards,
                    cardBackAssetPath: cardBackAssetPath,
                  ),
                ),
              Center(
                child: _TrickCenter(
                  state: state,
                  me: me,
                  top: top,
                  left: left,
                  right: right,
                ),
              ),
              if (me != null)
                Align(
                  alignment: const Alignment(-0.86, 0.88),
                  child: _CurrentTurnBadge(
                    nickname: me!.nickname,
                    isCurrent: me!.id == state.currentPlayerId,
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

enum _Seat { left, top, right }

class _OpponentSeat extends StatelessWidget {
  const _OpponentSeat({
    required this.player,
    required this.seat,
    required this.isCurrent,
    required this.cardCount,
    required this.cardBackAssetPath,
  });

  final Player player;
  final _Seat seat;
  final bool isCurrent;
  final int cardCount;
  final String? cardBackAssetPath;

  @override
  Widget build(BuildContext context) {
    final avatar = _PlayerAvatar(player: player, highlight: isCurrent);
    final cards = _CardBackFan(
      count: cardCount,
      vertical: seat == _Seat.left || seat == _Seat.right,
      reverse: seat == _Seat.right,
      cardBackAssetPath: cardBackAssetPath,
    );

    switch (seat) {
      case _Seat.top:
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [avatar, const SizedBox(height: 6), cards],
        );
      case _Seat.left:
        return Row(
          mainAxisSize: MainAxisSize.min,
          children: [avatar, const SizedBox(width: 6), cards],
        );
      case _Seat.right:
        return Row(
          mainAxisSize: MainAxisSize.min,
          children: [cards, const SizedBox(width: 6), avatar],
        );
    }
  }
}

class _PlayerAvatar extends StatelessWidget {
  const _PlayerAvatar({required this.player, required this.highlight});

  final Player player;
  final bool highlight;

  @override
  Widget build(BuildContext context) {
    final initials = _initials(player.nickname);
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        DecoratedBox(
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            gradient: const LinearGradient(
              colors: [Color(0xFFF4D896), Color(0xFFD7B46A)],
            ),
            border: Border.all(
              color: highlight
                  ? const Color(0xFFFFF5D6)
                  : const Color(0x996A4A2D),
              width: highlight ? 2.2 : 1.2,
            ),
            boxShadow: const [
              BoxShadow(
                color: Color(0x33000000),
                blurRadius: 8,
                offset: Offset(0, 4),
              ),
            ],
          ),
          child: SizedBox(
            width: 42,
            height: 42,
            child: Center(
              child: Text(
                initials,
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.w800,
                  color: const Color(0xFF2A241A),
                ),
              ),
            ),
          ),
        ),
        const SizedBox(height: 3),
        ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 86),
          child: Text(
            player.nickname,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            textAlign: TextAlign.center,
            style: Theme.of(
              context,
            ).textTheme.labelMedium?.copyWith(color: const Color(0xFFF8F0DB)),
          ),
        ),
      ],
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
  final String? cardBackAssetPath;

  @override
  Widget build(BuildContext context) {
    final safeCount = count.clamp(2, 10);
    final cardWidth = 32.0;
    final cardHeight = 46.0;
    final spread = 9.0;
    final width = cardWidth + ((safeCount - 1) * spread);
    final height = cardHeight + 12;

    Widget fan = SizedBox(
      width: width,
      height: height,
      child: Stack(
        clipBehavior: Clip.none,
        children: List.generate(safeCount, (i) {
          final midpoint = (safeCount - 1) / 2;
          final orderIndex = reverse ? (safeCount - 1 - i) : i;
          final distance = orderIndex - midpoint;
          return Positioned(
            left: orderIndex * spread,
            top: distance.abs() * 0.7,
            child: Transform.rotate(
              angle: distance * 0.06,
              child: _CardBack(cardBackAssetPath: cardBackAssetPath),
            ),
          );
        }),
      ),
    );

    if (vertical) {
      fan = RotatedBox(quarterTurns: 1, child: fan);
    }
    return fan;
  }
}

class _CardBack extends StatelessWidget {
  const _CardBack({required this.cardBackAssetPath});

  final String? cardBackAssetPath;

  @override
  Widget build(BuildContext context) {
    if (cardBackAssetPath != null) {
      final isSvg = cardBackAssetPath!.toLowerCase().endsWith('.svg');
      return ClipRRect(
        borderRadius: BorderRadius.circular(7),
        child: SizedBox(
          width: 32,
          height: 46,
          child: isSvg
              ? SvgPicture.asset(cardBackAssetPath!, fit: BoxFit.cover)
              : Image.asset(cardBackAssetPath!, fit: BoxFit.cover),
        ),
      );
    }

    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(7),
        gradient: const LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0xFFD5E6FF), Color(0xFF8DB3E6)],
        ),
        border: Border.all(color: const Color(0xFFF4FBFF), width: 1),
      ),
      child: SizedBox(
        width: 32,
        height: 46,
        child: Center(
          child: DecoratedBox(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(4),
              border: Border.all(color: const Color(0x88FFFFFF)),
            ),
            child: const SizedBox(width: 22, height: 30),
          ),
        ),
      ),
    );
  }
}

class _TrickCenter extends StatelessWidget {
  const _TrickCenter({
    required this.state,
    required this.me,
    required this.top,
    required this.left,
    required this.right,
  });

  final SuecaGameState state;
  final Player? me;
  final Player? top;
  final Player? left;
  final Player? right;

  @override
  Widget build(BuildContext context) {
    final pile = <_PileCard>[
      if (top != null && state.tableCards[top!.id] != null)
        _PileCard(
          seatTag: 'top',
          card: state.tableCards[top!.id]!,
          offset: const Offset(-6, -18),
          angle: -0.11,
        ),
      if (left != null && state.tableCards[left!.id] != null)
        _PileCard(
          seatTag: 'left',
          card: state.tableCards[left!.id]!,
          offset: const Offset(-20, -2),
          angle: -0.22,
        ),
      if (right != null && state.tableCards[right!.id] != null)
        _PileCard(
          seatTag: 'right',
          card: state.tableCards[right!.id]!,
          offset: const Offset(16, 6),
          angle: 0.19,
        ),
      if (me != null && state.tableCards[me!.id] != null)
        _PileCard(
          seatTag: 'me',
          card: state.tableCards[me!.id]!,
          offset: const Offset(4, 20),
          angle: 0.06,
        ),
    ];

    return SizedBox(
      width: 230,
      height: 230,
      child: Stack(
        children: [
          Align(
            alignment: Alignment.center,
            child: DecoratedBox(
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [const Color(0x4AF8F0DB), const Color(0x0DFFFFFF)],
                ),
              ),
              child: const SizedBox(width: 128, height: 128),
            ),
          ),
          if (pile.isEmpty)
            Center(
              child: Text(
                'Aguardando cartas da vaza',
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: const Color(0xD7F8F0DB),
                ),
              ),
            ),
          ...pile.asMap().entries.map((entry) {
            final i = entry.key;
            final item = entry.value;
            return Align(
              alignment: Alignment.center,
              child: Transform.translate(
                offset: item.offset,
                child: RevealSlideFade(
                  key: ValueKey<String>(
                    'pile_${item.seatTag}_${item.card.compactLabel}',
                  ),
                  delay: Duration(milliseconds: 40 + (i * 65)),
                  beginOffset: const Offset(0, 0.06),
                  child: _TablePlayedCard(
                    card: item.card,
                    angle: item.angle,
                    seatTag: item.seatTag,
                  ),
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}

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
        width: 70,
        height: 94,
        child: DecoratedBox(
          decoration: BoxDecoration(
            color: const Color(0xFFFFFAEE),
            borderRadius: BorderRadius.circular(10),
            border: Border.all(color: const Color(0xFFB08D49), width: 1.1),
            boxShadow: const [
              BoxShadow(
                color: Color(0x35000000),
                blurRadius: 8,
                offset: Offset(0, 5),
              ),
            ],
          ),
          child: AnimatedSwitcher(
            duration: const Duration(milliseconds: 260),
            switchInCurve: Curves.easeOutBack,
            switchOutCurve: Curves.easeInCubic,
            transitionBuilder: (child, animation) {
              return FadeTransition(
                opacity: animation,
                child: ScaleTransition(scale: animation, child: child),
              );
            },
            child: ClipRRect(
              key: ValueKey<String>('${seatTag}_${card.compactLabel}'),
              borderRadius: BorderRadius.circular(9),
              child: SvgPicture.asset(
                _cardFrontAssetPath(card),
                fit: BoxFit.cover,
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _CurrentTurnBadge extends StatelessWidget {
  const _CurrentTurnBadge({required this.nickname, required this.isCurrent});

  final String nickname;
  final bool isCurrent;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        color: isCurrent ? const Color(0xA8D7B46A) : const Color(0x9A0D2F23),
        border: Border.all(
          color: isCurrent ? const Color(0xFFF8F0DB) : const Color(0x63D7B46A),
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
        child: Text(
          isCurrent ? 'Tua vez' : nickname,
          style: Theme.of(context).textTheme.labelLarge?.copyWith(
            color: isCurrent
                ? const Color(0xFF2A241A)
                : const Color(0xFFF8F0DB),
            fontWeight: FontWeight.w700,
          ),
        ),
      ),
    );
  }
}

class _MyHandArea extends StatelessWidget {
  const _MyHandArea({
    required this.me,
    required this.hand,
    required this.isBusy,
    required this.canPlay,
    required this.onPlayCard,
  });

  final Player? me;
  final List<SuecaCard> hand;
  final bool isBusy;
  final bool canPlay;
  final Future<void> Function(SuecaCard card) onPlayCard;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(16),
        color: const Color(0x400A2A1F),
        border: Border.all(color: const Color(0x66D7B46A)),
      ),
      child: Padding(
        padding: const EdgeInsets.fromLTRB(10, 8, 10, 8),
        child: Column(
          children: [
            Row(
              children: [
                if (me != null) ...[
                  _PlayerAvatar(player: me!, highlight: false),
                  const SizedBox(width: 8),
                ],
                Expanded(
                  child: Text(
                    'A tua mao (${hand.length} cartas)',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      color: const Color(0xFFF8F0DB),
                    ),
                  ),
                ),
                if (isBusy)
                  const SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
              ],
            ),
            const SizedBox(height: 8),
            Expanded(
              child: _FannedHand(
                hand: hand,
                canInteract: canPlay && !isBusy,
                onPlayCard: onPlayCard,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _FannedHand extends StatelessWidget {
  const _FannedHand({
    required this.hand,
    required this.canInteract,
    required this.onPlayCard,
  });

  final List<SuecaCard> hand;
  final bool canInteract;
  final Future<void> Function(SuecaCard card) onPlayCard;

  @override
  Widget build(BuildContext context) {
    if (hand.isEmpty) {
      return Center(
        child: Text(
          'Sem cartas na mao.',
          style: Theme.of(
            context,
          ).textTheme.bodyMedium?.copyWith(color: const Color(0xFFF8F0DB)),
        ),
      );
    }

    return LayoutBuilder(
      builder: (context, constraints) {
        final cardWidth = 84.0;
        final maxSpread = hand.length > 1
            ? (constraints.maxWidth - cardWidth) / (hand.length - 1)
            : 0.0;
        final spread = hand.length > 1
            ? math.max(16.0, math.min(42.0, maxSpread))
            : 0.0;
        final totalWidth = cardWidth + ((hand.length - 1) * spread);
        final startLeft = (constraints.maxWidth - totalWidth) / 2;

        return Stack(
          clipBehavior: Clip.none,
          children: List.generate(hand.length, (index) {
            final card = hand[index];
            final midpoint = (hand.length - 1) / 2;
            final delta = index - midpoint;
            final angle = delta * 0.07;

            return Positioned(
              left: startLeft + (index * spread),
              bottom: math.max(0.0, 6 - delta.abs() * 1.4),
              child: RevealSlideFade(
                key: ValueKey<String>('hand_${card.compactLabel}_$index'),
                delay: Duration(milliseconds: 80 + (index * 45)),
                beginOffset: const Offset(0, 0.06),
                child: Transform.rotate(
                  angle: angle,
                  child: _HandCard(
                    card: card,
                    isBusy: !canInteract,
                    onPressed: () => onPlayCard(card),
                  ),
                ),
              ),
            );
          }),
        );
      },
    );
  }
}

class _HandCard extends StatefulWidget {
  const _HandCard({
    required this.card,
    required this.isBusy,
    required this.onPressed,
  });

  final SuecaCard card;
  final bool isBusy;
  final VoidCallback onPressed;

  @override
  State<_HandCard> createState() => _HandCardState();
}

class _HandCardState extends State<_HandCard> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 160),
        curve: Curves.easeOutCubic,
        transform: Matrix4.translationValues(0, _isHovered ? -10 : 0, 0),
        child: Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: widget.isBusy ? null : widget.onPressed,
            borderRadius: BorderRadius.circular(13),
            child: Ink(
              width: 84,
              height: 122,
              decoration: BoxDecoration(
                color: const Color(0xFFFFFAEE),
                borderRadius: BorderRadius.circular(13),
                border: Border.all(color: const Color(0xFFB08D49), width: 1.2),
                boxShadow: [
                  BoxShadow(
                    color: const Color(0x2A000000),
                    blurRadius: _isHovered ? 15 : 10,
                    offset: Offset(0, _isHovered ? 9 : 6),
                  ),
                ],
              ),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(11.8),
                child: SvgPicture.asset(
                  _cardFrontAssetPath(widget.card),
                  fit: BoxFit.cover,
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

Player? _playerAt(List<Player> players, int index) {
  if (index < 0 || index >= players.length) {
    return null;
  }
  return players[index];
}

List<Player> _orderedPlayers(SuecaGameState state) {
  final myIndex = state.players.indexWhere(
    (player) => player.id == state.myPlayerId,
  );
  if (myIndex < 0) {
    return state.players;
  }
  return [...state.players.skip(myIndex), ...state.players.take(myIndex)];
}

String _phaseLabel(GamePhase phase) {
  switch (phase) {
    case GamePhase.waitingForPlayers:
      return 'A espera de jogadores';
    case GamePhase.dealingCards:
      return 'Distribuicao';
    case GamePhase.playingTrick:
      return 'A jogar a vaza';
    case GamePhase.scoring:
      return 'Contagem';
    case GamePhase.finished:
      return 'Terminada';
  }
}

String _suitLabel(Suit suit) {
  switch (suit) {
    case Suit.clubs:
      return 'Paus';
    case Suit.diamonds:
      return 'Ouros';
    case Suit.hearts:
      return 'Copas';
    case Suit.spades:
      return 'Espadas';
  }
}

String _cardFrontAssetPath(SuecaCard card) {
  final rankToken = switch (card.rank) {
    1 => 'ace',
    11 => 'jack',
    12 => 'queen',
    13 => 'king',
    _ => '${card.rank}',
  };
  return 'assets/cards/svg-cards/${rankToken}_of_${card.suit.name}.svg';
}

String _initials(String nickname) {
  final trimmed = nickname.trim();
  if (trimmed.isEmpty) {
    return '?';
  }
  final parts = trimmed
      .split(RegExp(r'\s+'))
      .where((it) => it.isNotEmpty)
      .toList();
  if (parts.length == 1) {
    return parts.first.substring(0, 1).toUpperCase();
  }
  final first = parts.first.substring(0, 1);
  final second = parts[1].substring(0, 1);
  return '$first$second'.toUpperCase();
}
