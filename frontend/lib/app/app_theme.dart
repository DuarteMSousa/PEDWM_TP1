import 'package:flutter/material.dart';

class AppTheme {
  AppTheme._();

  static const feltDark = Color(0xFF0E3A2B);
  static const feltMid = Color(0xFF155B42);
  static const feltLight = Color(0xFF2D8A63);
  static const woodBrown = Color(0xFF6A4A2D);
  static const gold = Color(0xFFD7B46A);
  static const cream = Color(0xFFF8F0DB);
  static const ivory = Color(0xFFFFFAEE);
  static const ink = Color(0xFF2A241A);
  static const danger = Color(0xFFAA2E1F);

  static ThemeData light() {
    final baseScheme = ColorScheme.fromSeed(
      seedColor: feltMid,
      brightness: Brightness.light,
    );
    final colorScheme = baseScheme.copyWith(
      primary: feltMid,
      onPrimary: ivory,
      secondary: gold,
      onSecondary: ink,
      surface: ivory,
      onSurface: ink,
      error: danger,
    );

    final baseTextTheme = ThemeData.light().textTheme;
    final textTheme = baseTextTheme
        .apply(bodyColor: ink, displayColor: ink, fontFamily: 'Georgia')
        .copyWith(
          headlineLarge: baseTextTheme.headlineLarge?.copyWith(
            fontWeight: FontWeight.w700,
            letterSpacing: 0.4,
          ),
          headlineMedium: baseTextTheme.headlineMedium?.copyWith(
            fontWeight: FontWeight.w700,
            letterSpacing: 0.3,
          ),
          headlineSmall: baseTextTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.w700,
          ),
          titleLarge: baseTextTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.w700,
          ),
          titleMedium: baseTextTheme.titleMedium?.copyWith(
            fontWeight: FontWeight.w600,
          ),
          bodyLarge: baseTextTheme.bodyLarge?.copyWith(height: 1.3),
          bodyMedium: baseTextTheme.bodyMedium?.copyWith(height: 1.35),
        );

    return ThemeData(
      useMaterial3: true,
      colorScheme: colorScheme,
      textTheme: textTheme,
      scaffoldBackgroundColor: feltDark,
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.transparent,
        foregroundColor: ivory,
        elevation: 0,
        centerTitle: true,
        scrolledUnderElevation: 0,
      ),
      cardTheme: CardThemeData(
        color: ivory,
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: const Color(0xFFF2E8CF),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 14,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: const BorderSide(color: Color(0xFFA1864C), width: 1.2),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: const BorderSide(color: Color(0xFFA1864C), width: 1.2),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: const BorderSide(color: feltMid, width: 1.6),
        ),
        labelStyle: const TextStyle(fontWeight: FontWeight.w600),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: feltMid,
          foregroundColor: ivory,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 13),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
          ),
          textStyle: const TextStyle(fontWeight: FontWeight.w700),
        ),
      ),
      floatingActionButtonTheme: const FloatingActionButtonThemeData(
        backgroundColor: gold,
        foregroundColor: ink,
      ),
      snackBarTheme: SnackBarThemeData(
        behavior: SnackBarBehavior.floating,
        backgroundColor: const Color(0xFF3B3022),
        contentTextStyle: textTheme.bodyMedium?.copyWith(color: ivory),
      ),
    );
  }
}
