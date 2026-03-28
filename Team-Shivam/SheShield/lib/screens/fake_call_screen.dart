import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/// Full-screen fake incoming call UI.
///
/// Simulates a phone call for escape situations.
/// No real call is made — purely a visual/audio simulation.
class FakeCallScreen extends StatefulWidget {
  const FakeCallScreen({super.key});

  @override
  State<FakeCallScreen> createState() => _FakeCallScreenState();
}

class _FakeCallScreenState extends State<FakeCallScreen>
    with TickerProviderStateMixin {
  late AnimationController _pulseController;
  late AnimationController _slideController;
  late Animation<double> _pulseAnimation;
  late Animation<Offset> _slideAnimation;
  Timer? _vibrateTimer;

  // Fake caller info — configurable
  final String _callerName = 'Mom';
  final String _callerNumber = '+91 98765 43210';

  @override
  void initState() {
    super.initState();

    // Pulse animation for the call icon
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1200),
    )..repeat(reverse: true);
    _pulseAnimation = Tween<double>(begin: 1.0, end: 1.15).animate(
      CurvedAnimation(parent: _pulseController, curve: Curves.easeInOut),
    );

    // Slide animation for accept button
    _slideController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1500),
    )..repeat(reverse: true);
    _slideAnimation = Tween<Offset>(
      begin: const Offset(0, 0.15),
      end: const Offset(0, -0.15),
    ).animate(CurvedAnimation(
        parent: _slideController, curve: Curves.easeInOut));

    // Simulate vibration pattern (like a real call)
    _vibrateTimer = Timer.periodic(
      const Duration(milliseconds: 1000),
      (_) => HapticFeedback.mediumImpact(),
    );
    HapticFeedback.mediumImpact();
  }

  @override
  void dispose() {
    _pulseController.dispose();
    _slideController.dispose();
    _vibrateTimer?.cancel();
    super.dispose();
  }

  void _acceptCall() {
    _vibrateTimer?.cancel();
    Navigator.pop(context);
    // Could navigate to a "In Call" screen — for now just close
  }

  void _declineCall() {
    _vibrateTimer?.cancel();
    Navigator.pop(context);
  }

  @override
  Widget build(BuildContext context) {
    return PopScope(
      canPop: false,
      child: Scaffold(
        body: Container(
          width: double.infinity,
          height: double.infinity,
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [
                Color(0xFF1A1A2E),
                Color(0xFF0A0A0F),
              ],
            ),
          ),
          child: SafeArea(
            child: Column(
              children: [
                const Spacer(flex: 2),
                // ── Caller avatar ──
                AnimatedBuilder(
                  animation: _pulseAnimation,
                  builder: (_, __) => Transform.scale(
                    scale: _pulseAnimation.value,
                    child: Container(
                      width: 120,
                      height: 120,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: const Color(0xFF42A5F5).withValues(alpha: 0.15),
                        boxShadow: [
                          BoxShadow(
                            color: const Color(0xFF42A5F5).withValues(alpha: 0.2),
                            blurRadius: 40,
                            spreadRadius: 10,
                          ),
                        ],
                      ),
                      child: const Center(
                        child: Icon(
                          Icons.person,
                          color: Color(0xFF42A5F5),
                          size: 56,
                        ),
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 32),
                // ── Caller name ──
                Text(
                  _callerName,
                  style: const TextStyle(
                    fontSize: 32,
                    fontWeight: FontWeight.w700,
                    color: Color(0xFFF0F0F5),
                    letterSpacing: 0.5,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  _callerNumber,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w400,
                    color: const Color(0xFFF0F0F5).withValues(alpha: 0.5),
                    letterSpacing: 1,
                  ),
                ),
                const SizedBox(height: 16),
                // ── "Incoming call" label ──
                Container(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 16, vertical: 6),
                  decoration: BoxDecoration(
                    color: const Color(0xFF00E676).withValues(alpha: 0.12),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Container(
                        width: 8,
                        height: 8,
                        decoration: const BoxDecoration(
                          shape: BoxShape.circle,
                          color: Color(0xFF00E676),
                        ),
                      ),
                      const SizedBox(width: 8),
                      const Text(
                        'Incoming Call',
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w600,
                          color: Color(0xFF00E676),
                          letterSpacing: 0.5,
                        ),
                      ),
                    ],
                  ),
                ),
                const Spacer(flex: 3),
                // ── Action buttons ──
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 48),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      // Decline
                      _CallButton(
                        icon: Icons.call_end,
                        color: const Color(0xFFE53935),
                        label: 'Decline',
                        onTap: _declineCall,
                      ),
                      // Accept
                      SlideTransition(
                        position: _slideAnimation,
                        child: _CallButton(
                          icon: Icons.call,
                          color: const Color(0xFF00E676),
                          label: 'Accept',
                          onTap: _acceptCall,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 48),
                // ── Swipe hint ──
                Text(
                  'Tap to answer or decline',
                  style: TextStyle(
                    fontSize: 12,
                    color: const Color(0xFFF0F0F5).withValues(alpha: 0.3),
                    letterSpacing: 0.5,
                  ),
                ),
                const SizedBox(height: 32),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

/// Circular call action button.
class _CallButton extends StatelessWidget {
  final IconData icon;
  final Color color;
  final String label;
  final VoidCallback onTap;

  const _CallButton({
    required this.icon,
    required this.color,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 72,
            height: 72,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: color,
              boxShadow: [
                BoxShadow(
                  color: color.withValues(alpha: 0.4),
                  blurRadius: 20,
                  spreadRadius: 2,
                ),
              ],
            ),
            child: Icon(icon, color: Colors.white, size: 32),
          ),
          const SizedBox(height: 10),
          Text(
            label,
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w500,
              color: const Color(0xFFF0F0F5).withValues(alpha: 0.7),
            ),
          ),
        ],
      ),
    );
  }
}
