import 'package:flutter/material.dart';

import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../domain/entities/game_summary.dart';
import '../state/replay_player_controller.dart';

class ReplayPlayerPage extends StatefulWidget {
  const ReplayPlayerPage({super.key, required this.game});

  final GameSummary game;

  @override
  State<ReplayPlayerPage> createState() => _ReplayPlayerPageState();
}

class _ReplayPlayerPageState extends State<ReplayPlayerPage> {
  late final ReplayPlayerController _controller;

  @override
  void initState() {
    super.initState();
    _controller = ReplayPlayerController(game: widget.game);
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
          final visibleEvents = _controller.visibleEvents;
          final currentEvent = _controller.currentEvent;
          final playerNames = _buildPlayerNames(widget.game);

          return TableBackground(
            child: Column(
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
                  child: SectionCard(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Evento ${_controller.currentIndex}/${_controller.totalEvents}',
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                        const SizedBox(height: 8),
                        if (currentEvent != null)
                          _EventHighlight(
                              event: currentEvent, playerNames: playerNames)
                        else
                          const Text('Pressiona play para iniciar o replay.'),
                        const SizedBox(height: 10),
                        SliderTheme(
                          data: SliderTheme.of(context).copyWith(
                            activeTrackColor: const Color(0xFF155B42),
                            thumbColor: const Color(0xFFD7B46A),
                            inactiveTrackColor: const Color(0x336A4A2D),
                          ),
                          child: Slider(
                            value: _controller.currentIndex.toDouble(),
                            min: 0,
                            max: _controller.totalEvents.toDouble(),
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
                              tooltip: 'Inicio',
                            ),
                            IconButton(
                              onPressed: _controller.isAtStart
                                  ? null
                                  : _controller.stepBackward,
                              icon: const Icon(Icons.fast_rewind_rounded),
                              tooltip: 'Anterior',
                            ),
                            IconButton(
                              iconSize: 40,
                              onPressed: _controller.isPlaying
                                  ? _controller.pause
                                  : (_controller.isAtEnd
                                      ? null
                                      : _controller.play),
                              icon: Icon(
                                _controller.isPlaying
                                    ? Icons.pause_circle_filled_rounded
                                    : Icons.play_circle_filled_rounded,
                              ),
                              tooltip:
                                  _controller.isPlaying ? 'Pausar' : 'Play',
                            ),
                            IconButton(
                              onPressed: _controller.isAtEnd
                                  ? null
                                  : _controller.stepForward,
                              icon: const Icon(Icons.fast_forward_rounded),
                              tooltip: 'Proximo',
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 10),
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.fromLTRB(16, 0, 16, 24),
                    reverse: true,
                    itemCount: visibleEvents.length,
                    itemBuilder: (context, index) {
                      final reversedIndex =
                          visibleEvents.length - 1 - index;
                      final event = visibleEvents[reversedIndex];
                      final isCurrent = reversedIndex ==
                          _controller.currentIndex - 1;
                      return Padding(
                        padding: const EdgeInsets.only(bottom: 6),
                        child: _EventTile(
                          event: event,
                          isHighlighted: isCurrent,
                          playerNames: playerNames,
                        ),
                      );
                    },
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

Map<String, String> _buildPlayerNames(GameSummary game) {
  final names = <String, String>{};
  for (final p in game.players) {
    names[p.id] = p.username;
  }
  // Also extract from GAME_STARTED event teams payload
  for (final e in game.events) {
    if (e.type == 'GAME_STARTED' || e.type == 'GAME_ENDED') {
      final teams = e.payload?['teams'];
      if (teams is List) {
        for (final t in teams) {
          if (t is! Map) continue;
          final players = t['players'];
          if (players is List) {
            for (final p in players) {
              if (p is Map) {
                final id = p['id']?.toString();
                final name = p['name']?.toString();
                if (id != null && name != null) {
                  names[id] = name;
                }
              }
            }
          }
        }
      }
    }
  }
  return names;
}

class _EventHighlight extends StatelessWidget {
  const _EventHighlight({required this.event, required this.playerNames});

  final GameEvent event;
  final Map<String, String> playerNames;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        color: const Color(0xFF155B42),
        border: Border.all(color: const Color(0xFFD7B46A)),
      ),
      child: Row(
        children: [
          Icon(_iconForType(event.type),
              color: const Color(0xFFF8F0DB), size: 24),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  _labelForType(event.type),
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        color: const Color(0xFFF8F0DB),
                        fontWeight: FontWeight.bold,
                      ),
                ),
                if (_descriptionForEvent(event, playerNames).isNotEmpty)
                  Text(
                    _descriptionForEvent(event, playerNames),
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: const Color(0xCCF8F0DB),
                        ),
                  ),
              ],
            ),
          ),
        ],
      ),
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
    final desc = _descriptionForEvent(event, playerNames);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: isHighlighted
            ? const Color(0x40D7B46A)
            : const Color(0x1A6A4A2D),
        border: isHighlighted
            ? Border.all(color: const Color(0xFFD7B46A), width: 1.5)
            : null,
      ),
      child: Row(
        children: [
          Icon(_iconForType(event.type), size: 18,
              color: const Color(0xFF6A4A2D)),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  _labelForType(event.type),
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                ),
                if (desc.isNotEmpty)
                  Text(
                    desc,
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: const Color.fromARGB(255, 255, 255, 255),
                        ),
                    maxLines: 3,
                    overflow: TextOverflow.ellipsis,
                  ),
              ],
            ),
          ),
          Text(
            '#${event.sequence}',
            style: Theme.of(context).textTheme.labelSmall,
          ),
        ],
      ),
    );
  }
}

String _labelForType(String type) {
  switch (type) {
    case 'GAME_STARTED':
      return 'Jogo iniciado';
    case 'ROUND_STARTED':
      return 'Ronda iniciada';
    case 'TRICK_STARTED':
      return 'Vaza iniciada';
    case 'TRUMP_REVEALED':
      return 'Trunfo revelado';
    case 'CARD_DEALT':
      return 'Carta distribuída';
    case 'CARD_PLAYED':
      return 'Carta jogada';
    case 'TURN_CHANGED':
      return 'Vez mudou';
    case 'TRICK_ENDED':
      return 'Vaza terminada';
    case 'ROUND_ENDED':
      return 'Ronda terminada';
    case 'GAME_SCORE_UPDATED':
      return 'Pontuação atualizada';
    case 'GAME_ENDED':
      return 'Jogo terminado';
    default:
      return type;
  }
}

IconData _iconForType(String type) {
  switch (type) {
    case 'GAME_STARTED':
      return Icons.play_circle_outline_rounded;
    case 'ROUND_STARTED':
      return Icons.loop_rounded;
    case 'TRICK_STARTED':
      return Icons.flag_outlined;
    case 'TRUMP_REVEALED':
      return Icons.style_outlined;
    case 'CARD_DEALT':
      return Icons.back_hand_outlined;
    case 'CARD_PLAYED':
      return Icons.play_arrow_rounded;
    case 'TURN_CHANGED':
      return Icons.swap_horiz_rounded;
    case 'TRICK_ENDED':
      return Icons.check_circle_outline_rounded;
    case 'ROUND_ENDED':
      return Icons.sports_score_rounded;
    case 'GAME_SCORE_UPDATED':
      return Icons.scoreboard_outlined;
    case 'GAME_ENDED':
      return Icons.emoji_events_outlined;
    default:
      return Icons.info_outline_rounded;
  }
}

String _descriptionForEvent(GameEvent event, Map<String, String> names) {
  final payload = event.payload;
  if (payload == null) return '';

  String _name(String? id) {
    if (id == null || id.isEmpty) return '?';
    return names[id] ?? (id.length > 8 ? id.substring(0, 8) : id);
  }

  switch (event.type) {
    case 'GAME_STARTED':
      final teams = payload['teams'];
      if (teams is List) {
        final parts = <String>[];
        for (final t in teams) {
          if (t is! Map) continue;
          final teamId = t['id']?.toString() ?? '?';
          final players = t['players'];
          if (players is List) {
            final playerNames = players
                .whereType<Map>()
                .map((p) => p['name']?.toString() ?? '?')
                .join(', ');
            parts.add('$teamId: $playerNames');
          }
        }
        if (parts.isNotEmpty) return parts.join(' | ');
      }
      return '';
    case 'ROUND_STARTED':
      final roundNum = payload['roundNumber']?.toString() ?? '';
      final dealer = _name(payload['dealerId']?.toString());
      return 'Ronda $roundNum — Dealer: $dealer';
    case 'TRICK_STARTED':
      final leader = _name(payload['leaderId']?.toString());
      return 'Lider da vaza: $leader';
    case 'TRUMP_REVEALED':
      final suit = payload['suit']?.toString() ?? '';
      final card = payload['card'];
      if (card is Map) {
        final rank = card['rank']?.toString() ?? '';
        return '${_rankLabel(rank)} de ${_suitLabel(suit)} ${_suitSymbol(suit)}';
      }
      return 'Naipe: ${_suitLabel(suit)} ${_suitSymbol(suit)}';
    case 'CARD_DEALT':
      final player = _name(payload['playerId']?.toString());
      final card = payload['card'];
      if (card is Map) {
        final rank = card['rank']?.toString() ?? '';
        final suit = card['suit']?.toString() ?? '';
        return '$player recebeu ${_rankLabel(rank)} de ${_suitLabel(suit)} ${_suitSymbol(suit)}';
      }
      return 'Carta para $player';
    case 'CARD_PLAYED':
      final player = _name(payload['playerId']?.toString());
      final card = payload['card'];
      if (card is Map) {
        final rank = card['rank']?.toString() ?? '';
        final suit = card['suit']?.toString() ?? '';
        return '$player jogou ${_rankLabel(rank)} de ${_suitLabel(suit)} ${_suitSymbol(suit)}';
      }
      return player;
    case 'TURN_CHANGED':
      final player = _name(payload['playerId']?.toString());
      return 'Vez de $player';
    case 'TRICK_ENDED':
      final winner = _name(payload['winnerId']?.toString());
      final points = payload['points']?.toString() ?? '0';
      return '$winner ganhou a vaza (+$points pts)';
    case 'ROUND_ENDED':
      final winnerTeam = payload['winnerTeam']?.toString() ?? '';
      final scores = payload['score'];
      final scoreParts = <String>[];
      if (scores is List) {
        for (final s in scores) {
          if (s is Map) {
            scoreParts.add('${s['teamId']}: ${s['points']} pts');
          }
        }
      }
      final scoreStr = scoreParts.isNotEmpty ? ' (${scoreParts.join(' | ')})' : '';
      return 'Vencedor: $winnerTeam$scoreStr';
    case 'GAME_SCORE_UPDATED':
      final scores = payload['score'];
      if (scores is List) {
        return scores
            .whereType<Map>()
            .map((s) => '${s['teamId']}: ${s['points']} pts')
            .join(' | ');
      }
      return '';
    case 'GAME_ENDED':
      final winner = payload['winner']?.toString() ?? '';
      final finalScores = payload['finalScores'];
      final scoreParts = <String>[];
      if (finalScores is List) {
        for (final s in finalScores) {
          if (s is Map) {
            scoreParts.add('${s['teamId']}: ${s['points']} pts');
          }
        }
      }
      final scoreStr = scoreParts.isNotEmpty ? '\n${scoreParts.join(' | ')}' : '';
      return 'Vencedor: $winner$scoreStr';
    default:
      return '';
  }
}

String _suitSymbol(String suit) {
  switch (suit.toUpperCase()) {
    case 'HEARTS':
      return '♥';
    case 'DIAMONDS':
      return '♦';
    case 'CLUBS':
      return '♣';
    case 'SPADES':
      return '♠';
    default:
      return suit;
  }
}

String _suitLabel(String suit) {
  switch (suit.toUpperCase()) {
    case 'HEARTS':
      return 'Copas';
    case 'DIAMONDS':
      return 'Ouros';
    case 'CLUBS':
      return 'Paus';
    case 'SPADES':
      return 'Espadas';
    default:
      return suit;
  }
}

String _rankLabel(String rank) {
  switch (rank.toUpperCase()) {
    case 'A':
      return 'Ás';
    case 'K':
      return 'Rei';
    case 'Q':
      return 'Dama';
    case 'J':
      return 'Valete';
    case '7':
      return 'Visca';
    default:
      return rank;
  }
}
