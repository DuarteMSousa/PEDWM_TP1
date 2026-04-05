import 'dart:async';

import 'package:flutter/foundation.dart';

import '../../../game/domain/entities/card.dart';
import '../../../game/domain/entities/suit.dart';
import '../../domain/entities/game_summary.dart';
import '../../domain/repositories/replay_repository.dart';

class ReplayPlayerController extends ChangeNotifier {
  ReplayPlayerController({
    required ReplayRepository replayRepository,
    required this.userId,
    required GameSummary initialGame,
  }) : _replayRepository = replayRepository,
       _game = initialGame;

  final ReplayRepository _replayRepository;
  final String userId;

  GameSummary _game;
  int _currentIndex = 0;
  bool _isPlaying = false;
  bool _isLoading = false;
  String? _errorMessage;
  double _playbackSpeed = 1;
  Timer? _timer;
  List<ReplayFrame> _frames = const [ReplayFrame.empty()];

  GameSummary get game => _game;
  int get currentIndex => _currentIndex;
  bool get isPlaying => _isPlaying;
  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage;
  double get playbackSpeed => _playbackSpeed;
  int get totalEvents => _game.events.length;
  bool get isAtEnd => _currentIndex >= totalEvents;
  bool get isAtStart => _currentIndex <= 0;
  ReplayFrame get currentFrame => _frames[_currentIndex.clamp(0, totalEvents)];

  List<GameEvent> get visibleEvents =>
      _game.events.sublist(0, _currentIndex.clamp(0, totalEvents));

  List<GameEvent> get recentEvents {
    final end = _currentIndex.clamp(0, totalEvents);
    final start = (end - 5).clamp(0, end);
    return _game.events.sublist(start, end).reversed.toList(growable: false);
  }

  GameEvent? get currentEvent =>
      _currentIndex > 0 && _currentIndex <= totalEvents
          ? _game.events[_currentIndex - 1]
          : null;

  Future<void> load() async {
    pause();
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      _game = await _replayRepository.fetchGameReplay(
        userId: userId,
        gameId: _game.id,
      );
      _frames = _buildFrames(_game);
      _currentIndex = 0;
    } catch (error) {
      _errorMessage = error.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void play() {
    if (_isLoading || isAtEnd) {
      return;
    }

    _timer?.cancel();
    _isPlaying = true;
    notifyListeners();

    _timer = Timer.periodic(_playbackInterval, (_) {
      if (_currentIndex >= totalEvents) {
        pause();
        return;
      }
      _currentIndex++;
      notifyListeners();
    });
  }

  void pause() {
    _isPlaying = false;
    _timer?.cancel();
    _timer = null;
    notifyListeners();
  }

  void stepForward() {
    if (_currentIndex < totalEvents) {
      _currentIndex++;
      notifyListeners();
    }
  }

  void stepBackward() {
    if (_currentIndex > 0) {
      _currentIndex--;
      notifyListeners();
    }
  }

  void seekTo(int index) {
    _currentIndex = index.clamp(0, totalEvents);
    notifyListeners();
  }

  void reset() {
    pause();
    _currentIndex = 0;
    notifyListeners();
  }

  void setPlaybackSpeed(double value) {
    _playbackSpeed = value;
    if (_isPlaying) {
      play();
    } else {
      notifyListeners();
    }
  }

  Duration get _playbackInterval {
    final milliseconds = (900 / _playbackSpeed).round();
    return Duration(milliseconds: milliseconds.clamp(250, 1400));
  }

  List<ReplayFrame> _buildFrames(GameSummary game) {
    final players = <String, _MutableReplayPlayer>{};
    final teams = <String, _MutableReplayTeam>{};
    final teamOrder = <String>[];
    final tableCards = <String, SuecaCard>{};
    Suit? trumpSuit;
    String? currentPlayerId;
    String? lastWinnerId;
    var phaseLabel = 'Inicio do replay';

    for (var index = 0; index < game.players.length; index++) {
      final player = game.players[index];
      players[player.id] = _MutableReplayPlayer(
        id: player.id,
        name: player.username,
        sequence: index + 1,
      );
    }

    void ensureTeam(String teamId) {
      if (!teams.containsKey(teamId)) {
        teams[teamId] = _MutableReplayTeam(id: teamId);
        teamOrder.add(teamId);
      }
    }

    void assignPlayerToTeam(String playerId, String teamId) {
      ensureTeam(teamId);
      for (final team in teams.values) {
        team.playerIds.remove(playerId);
      }
      teams[teamId]!.playerIds.add(playerId);
      final player = players[playerId];
      if (player != null) {
        player.teamId = teamId;
        player.isActive = true;
        player.handCount = player.handCount == 0 ? 10 : player.handCount;
      }
    }

    void applyTeamsPayload(dynamic rawTeams) {
      if (rawTeams is! List) {
        return;
      }

      for (final team in teams.values) {
        team.playerIds.clear();
      }

      for (final rawTeam in rawTeams) {
        if (rawTeam is! Map) {
          continue;
        }
        final mappedTeam = Map<String, dynamic>.from(rawTeam);
        final teamId = mappedTeam['id']?.toString();
        if (teamId == null || teamId.isEmpty) {
          continue;
        }
        ensureTeam(teamId);

        final rawPlayers = mappedTeam['players'];
        if (rawPlayers is! List) {
          continue;
        }

        for (final rawPlayer in rawPlayers) {
          if (rawPlayer is! Map) {
            continue;
          }
          final mappedPlayer = Map<String, dynamic>.from(rawPlayer);
          final playerId = mappedPlayer['id']?.toString();
          if (playerId == null || playerId.isEmpty) {
            continue;
          }

          final existing = players[playerId];
          final player = existing ??
              _MutableReplayPlayer(
                id: playerId,
                name: mappedPlayer['name']?.toString() ?? 'Jogador',
                sequence: _toInt(mappedPlayer['sequence']) ?? players.length + 1,
              );

          player.name = mappedPlayer['name']?.toString() ?? player.name;
          player.sequence = _toInt(mappedPlayer['sequence']) ?? player.sequence;
          player.isActive = true;
          players[playerId] = player;
          assignPlayerToTeam(playerId, teamId);
        }
      }
    }

    void applyScorePayload(dynamic rawScores, {required bool round}) {
      if (rawScores is! List) {
        return;
      }
      for (final rawScore in rawScores) {
        if (rawScore is! Map) {
          continue;
        }
        final mappedScore = Map<String, dynamic>.from(rawScore);
        final teamId = mappedScore['teamId']?.toString();
        if (teamId == null || teamId.isEmpty) {
          continue;
        }
        ensureTeam(teamId);
        final points = _toInt(mappedScore['points']) ?? 0;
        if (round) {
          teams[teamId]!.roundScore = points;
        } else {
          teams[teamId]!.score = points;
        }
      }
    }

    ReplayFrame snapshot() {
      final activePlayers = players.values
          .where((player) => player.isActive)
          .toList(growable: false)
        ..sort((a, b) => a.sequence.compareTo(b.sequence));

      final orderedPlayers = _rotateForViewer(activePlayers);
      final seatPlayers = <ReplaySeat?>[
        orderedPlayers.isNotEmpty ? ReplaySeat.fromMutable(orderedPlayers[0]) : null,
        orderedPlayers.length > 1 ? ReplaySeat.fromMutable(orderedPlayers[1]) : null,
        orderedPlayers.length > 2 ? ReplaySeat.fromMutable(orderedPlayers[2]) : null,
        orderedPlayers.length > 3 ? ReplaySeat.fromMutable(orderedPlayers[3]) : null,
      ];

      final orderedTeams = teamOrder
          .map((teamId) => teams[teamId])
          .whereType<_MutableReplayTeam>()
          .map((team) => ReplayTeamSnapshot.fromMutable(team, players))
          .toList(growable: false);

      return ReplayFrame(
        seats: seatPlayers,
        teams: orderedTeams,
        tableCards: Map<String, SuecaCard>.from(tableCards),
        trumpSuit: trumpSuit,
        currentPlayerId: currentPlayerId,
        phaseLabel: phaseLabel,
        winnerId: lastWinnerId,
      );
    }

    final frames = <ReplayFrame>[snapshot()];

    for (final event in game.events) {
      final payload = event.payload ?? const <String, dynamic>{};

      switch (event.type) {
        case 'GAME_STARTED':
          applyTeamsPayload(payload['teams']);
          for (final player in players.values) {
            if (player.isActive) {
              player.handCount = 10;
            }
          }
          phaseLabel = 'Jogo iniciado';
          tableCards.clear();
          currentPlayerId = null;
          lastWinnerId = null;
          break;
        case 'PLAYER_LEFT':
          final playerId = payload['playerId']?.toString();
          if (playerId != null) {
            final leftPlayer = players[playerId];
            if (leftPlayer != null) {
              leftPlayer.isActive = false;
              leftPlayer.handCount = 0;
            }
            for (final team in teams.values) {
              team.playerIds.remove(playerId);
            }
            tableCards.remove(playerId);
            if (currentPlayerId == playerId) {
              currentPlayerId = null;
            }
          }
          phaseLabel = 'Jogador saiu';
          break;
        case 'PLAYER_JOINED':
          final playerId = payload['playerId']?.toString();
          if (playerId != null && playerId.isNotEmpty) {
            final player = players[playerId] ??
                _MutableReplayPlayer(
                  id: playerId,
                  name: payload['name']?.toString() ?? 'Jogador',
                  sequence: _toInt(payload['slot']) ?? players.length + 1,
                );
            player.name = payload['name']?.toString().trim().isNotEmpty == true
                ? payload['name']!.toString().trim()
                : player.name;
            player.sequence = _toInt(payload['slot']) ?? player.sequence;
            player.isActive = true;
            player.handCount = player.handCount == 0 ? 10 : player.handCount;
            players[playerId] = player;

            final teamWithSpace = teams.values.where((team) => team.playerIds.length < 2);
            if (teamWithSpace.isNotEmpty) {
              assignPlayerToTeam(playerId, teamWithSpace.first.id);
            }
          }
          phaseLabel = 'Jogador entrou';
          break;
        case 'ROUND_STARTED':
          phaseLabel = 'Ronda ${payload['roundNumber']?.toString() ?? ''}'.trim();
          tableCards.clear();
          lastWinnerId = null;
          currentPlayerId = payload['dealerId']?.toString();
          for (final team in teams.values) {
            team.roundScore = 0;
          }
          break;
        case 'TRICK_STARTED':
          phaseLabel = 'Nova vaza';
          currentPlayerId = payload['leaderId']?.toString();
          tableCards.clear();
          lastWinnerId = null;
          break;
        case 'TURN_CHANGED':
          currentPlayerId = payload['playerId']?.toString();
          phaseLabel = 'Mudanca de turno';
          break;
        case 'TRUMP_REVEALED':
          trumpSuit = _parseSuit(
            payload['suit']?.toString() ??
                (payload['card'] is Map
                    ? (payload['card'] as Map)['suit']?.toString()
                    : null),
          );
          phaseLabel = 'Trunfo revelado';
          break;
        case 'CARD_DEALT':
          final playerId = payload['playerId']?.toString();
          if (playerId != null && playerId.isNotEmpty) {
            final dealtPlayer = players[playerId];
            if (dealtPlayer != null) {
              dealtPlayer.handCount = (dealtPlayer.handCount + 1).clamp(0, 10);
            }
          }
          phaseLabel = 'Cartas distribuidas';
          break;
        case 'CARD_PLAYED':
          final playerId = payload['playerId']?.toString();
          final card = _parseCard(payload['card']);
          if (playerId != null && card != null) {
            tableCards[playerId] = card;
            final player = players[playerId];
            if (player != null && player.handCount > 0) {
              player.handCount--;
            }
          }
          phaseLabel = 'Carta jogada';
          break;
        case 'TRICK_ENDED':
          final winnerId = payload['winnerId']?.toString();
          final points = _toInt(payload['points']) ?? 0;
          lastWinnerId = winnerId;
          if (winnerId != null) {
            final team = teams.values.where((item) => item.playerIds.contains(winnerId));
            if (team.isNotEmpty) {
              team.first.roundScore += points;
            }
          }
          phaseLabel = 'Vaza concluida';
          tableCards.clear();
          currentPlayerId = winnerId;
          break;
        case 'ROUND_ENDED':
          applyScorePayload(payload['score'], round: true);
          tableCards.clear();
          phaseLabel = 'Ronda terminada';
          break;
        case 'GAME_SCORE_UPDATED':
          applyScorePayload(payload['score'], round: false);
          phaseLabel = 'Pontuacao atualizada';
          break;
        case 'GAME_ENDED':
          applyTeamsPayload(payload['teams']);
          applyScorePayload(payload['finalScores'], round: false);
          phaseLabel = 'Jogo terminado';
          currentPlayerId = null;
          tableCards.clear();
          break;
        default:
          break;
      }

      frames.add(snapshot());
    }

    return frames;
  }

  List<_MutableReplayPlayer> _rotateForViewer(
    List<_MutableReplayPlayer> orderedPlayers,
  ) {
    if (orderedPlayers.isEmpty) {
      return const <_MutableReplayPlayer>[];
    }

    final viewerIndex = orderedPlayers.indexWhere((player) => player.id == userId);
    if (viewerIndex <= 0) {
      return orderedPlayers;
    }

    return <_MutableReplayPlayer>[
      ...orderedPlayers.skip(viewerIndex),
      ...orderedPlayers.take(viewerIndex),
    ];
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }
}

class ReplayFrame {
  const ReplayFrame({
    required this.seats,
    required this.teams,
    required this.tableCards,
    required this.trumpSuit,
    required this.currentPlayerId,
    required this.phaseLabel,
    required this.winnerId,
  });

  const ReplayFrame.empty()
    : seats = const <ReplaySeat?>[null, null, null, null],
      teams = const <ReplayTeamSnapshot>[],
      tableCards = const <String, SuecaCard>{},
      trumpSuit = null,
      currentPlayerId = null,
      phaseLabel = 'A carregar replay',
      winnerId = null;

  final List<ReplaySeat?> seats;
  final List<ReplayTeamSnapshot> teams;
  final Map<String, SuecaCard> tableCards;
  final Suit? trumpSuit;
  final String? currentPlayerId;
  final String phaseLabel;
  final String? winnerId;
}

class ReplaySeat {
  const ReplaySeat({
    required this.id,
    required this.name,
    required this.sequence,
    required this.teamId,
    required this.handCount,
  });

  factory ReplaySeat.fromMutable(_MutableReplayPlayer player) {
    return ReplaySeat(
      id: player.id,
      name: player.name,
      sequence: player.sequence,
      teamId: player.teamId,
      handCount: player.handCount,
    );
  }

  final String id;
  final String name;
  final int sequence;
  final String? teamId;
  final int handCount;
}

class ReplayTeamSnapshot {
  const ReplayTeamSnapshot({
    required this.id,
    required this.playerNames,
    required this.score,
    required this.roundScore,
  });

  factory ReplayTeamSnapshot.fromMutable(
    _MutableReplayTeam team,
    Map<String, _MutableReplayPlayer> players,
  ) {
    final playerNames = team.playerIds
        .map((playerId) => players[playerId]?.name)
        .whereType<String>()
        .toList(growable: false);

    return ReplayTeamSnapshot(
      id: team.id,
      playerNames: playerNames,
      score: team.score,
      roundScore: team.roundScore,
    );
  }

  final String id;
  final List<String> playerNames;
  final int score;
  final int roundScore;
}

class _MutableReplayPlayer {
  _MutableReplayPlayer({
    required this.id,
    required this.name,
    required this.sequence,
    this.teamId,
    this.isActive = true,
    this.handCount = 0,
  });

  final String id;
  String name;
  int sequence;
  String? teamId;
  bool isActive;
  int handCount;
}

class _MutableReplayTeam {
  _MutableReplayTeam({required this.id});

  final String id;
  final List<String> playerIds = <String>[];
  int score = 0;
  int roundScore = 0;
}

SuecaCard? _parseCard(dynamic rawCard) {
  if (rawCard is! Map) {
    return null;
  }
  final map = Map<String, dynamic>.from(rawCard);
  final suit = _parseSuit(map['suit']?.toString());
  final rank = _parseRank(map['rank']);
  if (suit == null || rank == null) {
    return null;
  }
  return SuecaCard(suit: suit, rank: rank);
}

Suit? _parseSuit(String? rawSuit) {
  if (rawSuit == null || rawSuit.trim().isEmpty) {
    return null;
  }

  switch (rawSuit.trim().toUpperCase()) {
    case 'HEARTS':
      return Suit.hearts;
    case 'SPADES':
      return Suit.spades;
    case 'DIAMONDS':
      return Suit.diamonds;
    case 'CLUBS':
      return Suit.clubs;
    default:
      return null;
  }
}

int? _parseRank(dynamic rawRank) {
  if (rawRank is num) {
    return rawRank.toInt();
  }
  if (rawRank == null) {
    return null;
  }

  final token = rawRank.toString().trim().toUpperCase();
  switch (token) {
    case 'A':
      return 1;
    case 'K':
      return 13;
    case 'Q':
      return 12;
    case 'J':
      return 11;
    default:
      return int.tryParse(token);
  }
}

int? _toInt(dynamic value) {
  if (value is num) {
    return value.toInt();
  }
  if (value == null) {
    return null;
  }
  return int.tryParse(value.toString());
}
