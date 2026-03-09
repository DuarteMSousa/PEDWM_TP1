import 'dart:async';

import 'package:flutter/material.dart';

class RevealSlideFade extends StatefulWidget {
  const RevealSlideFade({
    super.key,
    required this.child,
    this.delay = Duration.zero,
    this.duration = const Duration(milliseconds: 380),
    this.beginOffset = const Offset(0, 0.08),
    this.curve = Curves.easeOutCubic,
  });

  final Widget child;
  final Duration delay;
  final Duration duration;
  final Offset beginOffset;
  final Curve curve;

  @override
  State<RevealSlideFade> createState() => _RevealSlideFadeState();
}

class _RevealSlideFadeState extends State<RevealSlideFade>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final Animation<double> _opacityAnimation;
  late final Animation<Offset> _offsetAnimation;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(vsync: this, duration: widget.duration);
    _opacityAnimation = CurvedAnimation(
      parent: _controller,
      curve: widget.curve,
    );
    _offsetAnimation = Tween<Offset>(
      begin: widget.beginOffset,
      end: Offset.zero,
    ).animate(CurvedAnimation(parent: _controller, curve: widget.curve));

    _timer = Timer(widget.delay, () {
      if (mounted) {
        _controller.forward();
      }
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return FadeTransition(
      opacity: _opacityAnimation,
      child: SlideTransition(position: _offsetAnimation, child: widget.child),
    );
  }
}

class HoverLift extends StatefulWidget {
  const HoverLift({
    super.key,
    required this.child,
    this.lift = 7,
    this.scale = 1.015,
    this.duration = const Duration(milliseconds: 180),
  });

  final Widget child;
  final double lift;
  final double scale;
  final Duration duration;

  @override
  State<HoverLift> createState() => _HoverLiftState();
}

class _HoverLiftState extends State<HoverLift> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedContainer(
        duration: widget.duration,
        curve: Curves.easeOutCubic,
        transform: Matrix4.translationValues(
          0,
          _isHovered ? -widget.lift : 0,
          0,
        ),
        child: AnimatedScale(
          duration: widget.duration,
          curve: Curves.easeOutCubic,
          scale: _isHovered ? widget.scale : 1,
          child: widget.child,
        ),
      ),
    );
  }
}
