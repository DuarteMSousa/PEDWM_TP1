import 'package:flutter/material.dart';

import '../../../../app/app_routes.dart';
import '../../../../core/shared_widgets/motion.dart';
import '../../../../core/shared_widgets/section_card.dart';
import '../../../../core/shared_widgets/table_background.dart';
import '../../domain/repositories/auth_repository.dart';
import '../state/auth_controller.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key, required this.authRepository});

  final AuthRepository authRepository;

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  late final AuthController _controller;
  final TextEditingController _nicknameController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _controller = AuthController(authRepository: widget.authRepository);
  }

  @override
  void dispose() {
    _controller.dispose();
    _nicknameController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    await _controller.login(_nicknameController.text);
    if (!mounted) {
      return;
    }
    if (_controller.currentUser != null) {
      Navigator.of(context).pushReplacementNamed(
        AppRoutes.lobby,
        arguments: _controller.currentUser,
      );
      return;
    }
    if (_controller.errorMessage != null) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text(_controller.errorMessage!)));
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    const ivory = Color(0xFFF8F0DB);

    final loginCard = SectionCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        mainAxisSize: MainAxisSize.min,
        children: [
          Text('Entrar no lobby', style: theme.textTheme.headlineSmall),
          const SizedBox(height: 10),
          Text(
            'Define o teu nickname para entrares no lobby.',
            style: theme.textTheme.bodyMedium,
          ),
          const SizedBox(height: 18),
          TextField(
            controller: _nicknameController,
            decoration: const InputDecoration(
              labelText: 'Nickname',
              hintText: 'ex: parceiro_1',
              prefixIcon: Icon(Icons.person_outline),
            ),
            textInputAction: TextInputAction.done,
            onSubmitted: (_) => _submit(),
          ),
          const SizedBox(height: 16),
          SizedBox(
            height: 48,
            child: ElevatedButton.icon(
              onPressed: _controller.isLoading ? null : _submit,
              icon: const Icon(Icons.login_rounded),
              label: Text(_controller.isLoading ? 'A entrar...' : 'Entrar'),
            ),
          ),
          const SizedBox(height: 8),
          TextButton(
            onPressed: _controller.isLoading
                ? null
                : () {
                    _nicknameController.text = 'guest';
                    _submit();
                  },
            child: const Text('Entrar como convidado'),
          ),
        ],
      ),
    );

    return Scaffold(
      body: AnimatedBuilder(
        animation: _controller,
        builder: (context, _) {
          return TableBackground(
            child: SafeArea(
              child: Center(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.all(20),
                  child: ConstrainedBox(
                    constraints: const BoxConstraints(maxWidth: 980),
                    child: LayoutBuilder(
                      builder: (context, constraints) {
                        final isWide = constraints.maxWidth > 760;
                        final spacing = isWide ? 24.0 : 18.0;

                        if (isWide) {
                          return Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Expanded(
                                child: RevealSlideFade(
                                  delay: const Duration(milliseconds: 60),
                                  beginOffset: const Offset(-0.04, 0),
                                  child: _IntroPanel(
                                    textTheme: theme.textTheme,
                                    foreground: ivory,
                                  ),
                                ),
                              ),
                              SizedBox(width: spacing),
                              Expanded(
                                child: RevealSlideFade(
                                  delay: const Duration(milliseconds: 180),
                                  beginOffset: const Offset(0.04, 0),
                                  child: loginCard,
                                ),
                              ),
                            ],
                          );
                        }

                        return Column(
                          crossAxisAlignment: CrossAxisAlignment.stretch,
                          children: [
                            RevealSlideFade(
                              delay: const Duration(milliseconds: 60),
                              beginOffset: const Offset(0, 0.03),
                              child: _IntroPanel(
                                textTheme: theme.textTheme,
                                foreground: ivory,
                              ),
                            ),
                            SizedBox(height: spacing),
                            RevealSlideFade(
                              delay: const Duration(milliseconds: 180),
                              beginOffset: const Offset(0, 0.04),
                              child: loginCard,
                            ),
                          ],
                        );
                      },
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

class _IntroPanel extends StatelessWidget {
  const _IntroPanel({required this.textTheme, required this.foreground});

  final TextTheme textTheme;
  final Color foreground;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 7),
            decoration: BoxDecoration(
              color: const Color(0x1FD7B46A),
              borderRadius: BorderRadius.circular(999),
              border: Border.all(color: const Color(0x66D7B46A)),
            ),
            child: Text(
              'TRUQUE, VAZA, VITORIA',
              style: textTheme.labelSmall?.copyWith(
                color: foreground,
                letterSpacing: 1.2,
                fontWeight: FontWeight.w700,
              ),
            ),
          ),
          const SizedBox(height: 18),
          Text(
            'Sueca Online',
            style: textTheme.headlineLarge?.copyWith(
              color: foreground,
              fontWeight: FontWeight.w800,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            'Um lobby com cara de mesa real. Junta a equipa, escolhe a sala e joga com estilo tradicional.',
            style: textTheme.bodyLarge?.copyWith(
              color: foreground.withAlpha(220),
            ),
          ),
          const SizedBox(height: 18),
          const Wrap(
            spacing: 10,
            runSpacing: 10,
            children: [
              _InfoPill(icon: Icons.groups_2_outlined, label: '4 Jogadores'),
              _InfoPill(icon: Icons.style_outlined, label: '40 Cartas'),
              _InfoPill(
                icon: Icons.emoji_events_outlined,
                label: 'Mais pontos vence',
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _InfoPill extends StatelessWidget {
  const _InfoPill({required this.icon, required this.label});

  final IconData icon;
  final String label;

  @override
  Widget build(BuildContext context) {
    const foreground = Color(0xFFF8F0DB);
    return DecoratedBox(
      decoration: BoxDecoration(
        color: const Color(0x141D7B5A),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: const Color(0x47CBB06C)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 9),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 16, color: foreground),
            const SizedBox(width: 8),
            Text(
              label,
              style: Theme.of(context).textTheme.labelLarge?.copyWith(
                color: foreground,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
