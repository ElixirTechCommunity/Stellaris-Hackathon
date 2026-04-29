import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:url_launcher/url_launcher.dart';
import 'home_screen.dart';

/// Blocking screen that requires the user to confirm they've enabled
/// power-off verification in device settings before accessing the app.
///
/// Uses SharedPreferences to persist the confirmation so it only shows once.
class SecuritySetupScreen extends StatefulWidget {
  const SecuritySetupScreen({super.key});

  /// Check if security setup has been completed previously.
  static Future<bool> isCompleted() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool('security_setup_done') ?? false;
  }

  @override
  State<SecuritySetupScreen> createState() => _SecuritySetupScreenState();
}

class _SecuritySetupScreenState extends State<SecuritySetupScreen> {
  bool _confirmed = false;

  /// Open Android security settings.
  Future<void> _openSecuritySettings() async {
    // Try opening a helpful guide for the user
    try {
      // Use the Android Settings activity action
      await launchUrl(
        Uri.parse('https://www.google.com/search?q=how+to+enable+power+off+verification'),
        mode: LaunchMode.externalApplication,
      );
    } catch (_) {}

    // Also try opening the actual settings page
    try {
      await launchUrl(
        Uri.parse('content://settings/system'),
        mode: LaunchMode.externalApplication,
      );
    } catch (_) {
      // Fallback: show instruction
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Go to Settings → Security → Power off verification'),
            duration: Duration(seconds: 4),
          ),
        );
      }
    }
  }

  /// Mark setup as done and proceed to home.
  Future<void> _completeSetup() async {
    if (!_confirmed) return;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('security_setup_done', true);
    if (mounted) {
      // If this is the root screen (after signup), go to HomeScreen
      // Otherwise pop back to SplashGate which will show HomeScreen
      final canPop = Navigator.of(context).canPop();
      if (canPop) {
        Navigator.of(context).pop(true);
      } else {
        Navigator.of(context).pushAndRemoveUntil(
          MaterialPageRoute(builder: (_) => const HomeScreen()),
          (route) => false,
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final cardBg = isDark ? const Color(0xFF1A1A2E) : Colors.white;
    final border = isDark
        ? Colors.white.withValues(alpha: 0.06)
        : Colors.black.withValues(alpha: 0.08);

    return PopScope(
      canPop: false, // Block back button
      child: Scaffold(
        body: Container(
          width: double.infinity,
          height: double.infinity,
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: isDark
                  ? [const Color(0xFF0A0A0F), const Color(0xFF1A1A2E)]
                  : [const Color(0xFFF5F5F8), Colors.white],
            ),
          ),
          child: SafeArea(
            child: Center(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(horizontal: 24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // ── Warning icon ──
                    Container(
                      width: 80,
                      height: 80,
                      decoration: BoxDecoration(
                        color: const Color(0xFFFF9800).withValues(alpha: 0.12),
                        shape: BoxShape.circle,
                      ),
                      child: const Center(
                        child: Text('⚠️', style: TextStyle(fontSize: 40)),
                      ),
                    ),
                    const SizedBox(height: 24),
                    // ── Title ──
                    Text(
                      'Security Setup Required',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.w800,
                        color: isDark
                            ? const Color(0xFFF0F0F5)
                            : const Color(0xFF1A1A2E),
                        letterSpacing: -0.3,
                      ),
                    ),
                    const SizedBox(height: 20),
                    // ── Info card ──
                    Container(
                      padding: const EdgeInsets.all(20),
                      decoration: BoxDecoration(
                        color: cardBg,
                        borderRadius: BorderRadius.circular(18),
                        border: Border.all(color: border),
                      ),
                      child: Column(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(12),
                            decoration: BoxDecoration(
                              color: const Color(0xFFFF9800).withValues(alpha: 0.08),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Text(
                              'For your safety, please enable power-off verification in your device settings. This prevents the phone from being switched off during emergencies.',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w500,
                                height: 1.5,
                                color: isDark
                                    ? const Color(0xFFC0C0D0)
                                    : const Color(0xFF3A3A4E),
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),
                          // ── Steps ──
                          _buildStep(isDark, '1', 'Open device Settings'),
                          const SizedBox(height: 8),
                          _buildStep(isDark, '2', 'Go to Security / Lock Screen'),
                          const SizedBox(height: 8),
                          _buildStep(isDark, '3', 'Enable "Power off verification"'),
                        ],
                      ),
                    ),
                    const SizedBox(height: 20),
                    // ── Go to Settings button ──
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        onPressed: _openSecuritySettings,
                        icon: const Icon(Icons.settings, size: 20),
                        label: const Text(
                          'Go to Settings',
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFFFF9800),
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(14),
                          ),
                          elevation: 0,
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    // ── Confirmation checkbox ──
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      decoration: BoxDecoration(
                        color: isDark
                            ? Colors.white.withValues(alpha: 0.04)
                            : Colors.black.withValues(alpha: 0.03),
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(
                          color: _confirmed
                              ? const Color(0xFF00E676).withValues(alpha: 0.3)
                              : border,
                        ),
                      ),
                      child: Row(
                        children: [
                          SizedBox(
                            width: 24,
                            height: 24,
                            child: Checkbox(
                              value: _confirmed,
                              onChanged: (val) =>
                                  setState(() => _confirmed = val ?? false),
                              activeColor: const Color(0xFF00E676),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(6),
                              ),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Text(
                              'I confirm I have enabled this setting',
                              style: TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: isDark
                                    ? const Color(0xFFC0C0D0)
                                    : const Color(0xFF3A3A4E),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 16),
                    // ── I Have Enabled It button ──
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        onPressed: _confirmed ? _completeSetup : null,
                        icon: Icon(
                          _confirmed
                              ? Icons.check_circle
                              : Icons.lock_outline,
                          size: 20,
                        ),
                        label: const Text(
                          'I Have Enabled It',
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w700,
                          ),
                        ),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: _confirmed
                              ? const Color(0xFF00E676)
                              : const Color(0xFF5A5A6E),
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(14),
                          ),
                          elevation: 0,
                          disabledBackgroundColor:
                              isDark ? const Color(0xFF2A2A3E) : const Color(0xFFE0E0E0),
                          disabledForegroundColor:
                              isDark ? const Color(0xFF5A5A6E) : const Color(0xFF9A9A9A),
                        ),
                      ),
                    ),
                    const SizedBox(height: 32),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildStep(bool isDark, String num, String text) {
    return Row(
      children: [
        Container(
          width: 28,
          height: 28,
          decoration: BoxDecoration(
            color: const Color(0xFFFF9800).withValues(alpha: 0.12),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Center(
            child: Text(
              num,
              style: const TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.w800,
                color: Color(0xFFFF9800),
              ),
            ),
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Text(
            text,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: isDark
                  ? const Color(0xFFC0C0D0)
                  : const Color(0xFF3A3A4E),
            ),
          ),
        ),
      ],
    );
  }
}
