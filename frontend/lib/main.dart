import 'package:flutter/widgets.dart';

import 'app/app.dart';
import 'app/app_dependencies.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  final dependencies = AppDependencies.create();
  runApp(SuecaApp(dependencies: dependencies));
}
