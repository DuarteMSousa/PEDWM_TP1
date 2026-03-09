Card assets for the Sueca table UI.

Folder layout:
- `front/` -> card front images
- `back/` -> card back images
- `svg-cards/` -> SVG fronts used by the current game page

Current expected naming for SVG fronts (inside `svg-cards/`):
- `ace_of_clubs.svg`, `ace_of_diamonds.svg`, `ace_of_hearts.svg`, `ace_of_spades.svg`
- `king_of_*.svg`, `queen_of_*.svg`, `jack_of_*.svg`
- numeric cards as `2_of_*.svg` ... `10_of_*.svg`

Suit names must match code enum names:
- `clubs`
- `diamonds`
- `hearts`
- `spades`

Back image examples (the app auto-detects the first that exists):
- `back/back.svg`
- `back/back_blue.svg`
- `back/back_red.svg`
- `back/back.png`
- `back/back_blue.png`
- `back/back_red.png`
