import 'player.dart';

class Team {
  Team({required this.id, required this.players, this.score = 0, this.roundScore = 0});

  final String id;
  final List<Player> players;
  int score = 0;
  int roundScore = 0;

  Team copyWith({String? id, List<Player>? players, int? score, int? roundScore}) {
    return Team(
      id: id ?? this.id,
      players: players ?? this.players,
      score: score ?? this.score,
      roundScore: roundScore ?? this.roundScore,
    );
  }

   factory Team.fromJson(Map<String, dynamic> json) {
    return Team(
      id: json['id'],
      score: json['score'] ?? 0,
      roundScore: json['roundScore'] ?? 0,
      players: (json['players'] as List? ?? [])
          .map((p) => Player.fromJson(p))
          .toList(),
    );
  }
  
}
