enum Suit { clubs, diamonds, hearts, spades }

extension SuitX on Suit {
  String get shortLabel {
    switch (this) {
      case Suit.clubs:
        return 'C';
      case Suit.diamonds:
        return 'D';
      case Suit.hearts:
        return 'H';
      case Suit.spades:
        return 'S';
    }
  }

  String get fullLabel {
    switch (this) {
      case Suit.clubs:
        return 'clubs';
      case Suit.diamonds:
        return 'diamonds';
      case Suit.hearts:
        return 'hearts';
      case Suit.spades:
        return 'spades';
    }
  }
}
