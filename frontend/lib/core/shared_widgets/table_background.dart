import 'package:flutter/material.dart';

class TableBackground extends StatelessWidget {
  const TableBackground({super.key, required this.child});

  final Widget child;

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0xFF0A2A1F), Color(0xFF0F3A2A), Color(0xFF14553D)],
        ),
      ),
      child: Stack(
        fit: StackFit.expand,
        children: [
          const IgnorePointer(
            child: Stack(
              children: [
                _BackgroundGlow(
                  alignment: Alignment(-1.15, -0.95),
                  size: 320,
                  color: Color(0x1FE2C57F),
                ),
                _BackgroundGlow(
                  alignment: Alignment(1.1, -0.8),
                  size: 260,
                  color: Color(0x161F7D5A),
                ),
                _BackgroundGlow(
                  alignment: Alignment(1.25, 0.95),
                  size: 360,
                  color: Color(0x1597F0B8),
                ),
                _BackgroundGlow(
                  alignment: Alignment(-1.2, 0.9),
                  size: 290,
                  color: Color(0x149DD3B8),
                ),
              ],
            ),
          ),
          child,
        ],
      ),
    );
  }
}

class _BackgroundGlow extends StatelessWidget {
  const _BackgroundGlow({
    required this.alignment,
    required this.size,
    required this.color,
  });

  final Alignment alignment;
  final double size;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: alignment,
      child: DecoratedBox(
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          boxShadow: [
            BoxShadow(color: color, blurRadius: 90, spreadRadius: 26),
          ],
        ),
        child: SizedBox.square(dimension: size),
      ),
    );
  }
}
