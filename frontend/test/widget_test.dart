import 'package:flutter_test/flutter_test.dart';

import 'package:sueca_pedwm/app/app.dart';
import 'package:sueca_pedwm/app/app_dependencies.dart';

void main() {
  testWidgets('login page loads', (WidgetTester tester) async {
    await tester.pumpWidget(SuecaApp(dependencies: AppDependencies.create()));
    await tester.pumpAndSettle();

    expect(find.text('Sueca Online'), findsOneWidget);
    expect(find.text('Entrar no lobby'), findsOneWidget);
  });
}
