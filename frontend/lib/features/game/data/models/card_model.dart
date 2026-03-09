import '../../domain/entities/card.dart';
import '../../domain/entities/suit.dart';

class CardModel {
  CardModel({required this.suit, required this.rank});

  final String suit;
  final int rank;

  factory CardModel.fromMap(Map<String, dynamic> map) {
    return CardModel(suit: map['suit'] as String, rank: map['rank'] as int);
  }

  SuecaCard toEntity() {
    final parsedSuit = Suit.values.firstWhere(
      (value) => value.name == suit,
      orElse: () => Suit.clubs,
    );
    return SuecaCard(suit: parsedSuit, rank: rank);
  }
}
