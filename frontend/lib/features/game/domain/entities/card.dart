import 'suit.dart';

class SuecaCard {
  const SuecaCard({required this.suit, required this.rank});

  final Suit suit;
  final int rank;

  int get points {
    switch (rank) {
      case 1:
        return 11;
      case 7:
        return 10;
      case 13:
        return 4;
      case 11:
        return 3;
      case 12:
        return 2;
      default:
        return 0;
    }
  }

  String get rankLabel {
    switch (rank) {
      case 1:
        return 'A';
      case 11:
        return 'J';
      case 12:
        return 'Q';
      case 13:
        return 'K';
      default:
        return rank.toString();
    }
  }

  String get compactLabel => '$rankLabel${suit.shortLabel}';

  String get backendId =>
      '${_backendRankToken(rank)}_${suit.name.toUpperCase()}';

  static String _backendRankToken(int rank) {
    switch (rank) {
      case 1:
        return 'A';
      case 11:
        return 'J';
      case 12:
        return 'Q';
      case 13:
        return 'K';
      default:
        return rank.toString();
    }
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) {
      return true;
    }
    return other is SuecaCard && other.suit == suit && other.rank == rank;
  }

  @override
  int get hashCode => Object.hash(suit, rank);
}
