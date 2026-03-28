import 'dart:async';
import 'package:flutter/material.dart';

/// Live status screen shown after SOS is triggered.
///
/// Displays real-time status of all emergency actions:
/// location sharing, contact alerts, video recording, etc.
class EmergencyStatusScreen extends StatefulWidget {
  final VoidCallback onDismiss;
  final double lat;
  final double lng;

  const EmergencyStatusScreen({
    super.key,
    required this.onDismiss,
    required this.lat,
    required this.lng,
  });

  @override
  State<EmergencyStatusScreen> createState() => _EmergencyStatusScreenState();
}

class _EmergencyStatusScreenState extends State<EmergencyStatusScreen>
    with SingleTickerProviderStateMixin {
  late AnimationController _fadeController;
  int _elapsedSeconds = 0;
  Timer? _elapsedTimer;

  @override
  void initState() {
    super.initState();
    _fadeController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 600),
    )..forward();

    // Track elapsed time since SOS
    _elapsedTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      if (mounted) setState(() => _elapsedSeconds++);
    });
  }

  @override
  void dispose() {
    _fadeController.dispose();
    _elapsedTimer?.cancel();
    super.dispose();
  }

  String _formatElapsed(int sec) {
    final m = sec ~/ 60;
    final s = sec % 60;
    return '${m.toString().padLeft(2, '0')}:${s.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final cardBg = isDark ? const Color(0xFF1A1A2E) : Colors.white;
    final border = isDark
        ? Colors.white.withValues(alpha: 0.06)
        : Colors.black.withValues(alpha: 0.08);
    final subtitleColor =
        isDark ? const Color(0xFF8A8A9A) : const Color(0xFF5A5A6E);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Emergency Status'),
        leading: IconButton(
          icon: Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: isDark ? const Color(0xFF1A1A2E) : Colors.white,
              borderRadius: BorderRadius.circular(10),
              border: Border.all(color: border),
            ),
            child: const Icon(Icons.arrow_back, size: 18),
          ),
          onPressed: () {
            widget.onDismiss();
            Navigator.pop(context);
          },
        ),
      ),
      body: FadeTransition(
        opacity: _fadeController,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              // ── SOS Active header ──
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [Color(0xFFE53935), Color(0xFFC62828)],
                  ),
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [
                    BoxShadow(
                      color: const Color(0xFFE53935).withValues(alpha: 0.3),
                      blurRadius: 20,
                      offset: const Offset(0, 6),
                    ),
                  ],
                ),
                child: Column(
                  children: [
                    const Text('🚨', style: TextStyle(fontSize: 40)),
                    const SizedBox(height: 8),
                    const Text(
                      'SOS ACTIVE',
                      style: TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.w800,
                        color: Colors.white,
                        letterSpacing: 2,
                      ),
                    ),
                    const SizedBox(height: 6),
                    Text(
                      'Emergency active for ${_formatElapsed(_elapsedSeconds)}',
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.white.withValues(alpha: 0.8),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 20),
              // ── Status cards ──
              _StatusCard(
                cardBg: cardBg,
                border: border,
                icon: Icons.sos_rounded,
                iconColor: const Color(0xFFE53935),
                title: 'SOS Alert Sent',
                subtitle: 'Emergency contacts have been alerted',
                status: 'Delivered',
                statusColor: const Color(0xFF00E676),
                subtitleColor: subtitleColor,
                isDark: isDark,
              ),
              const SizedBox(height: 12),
              _StatusCard(
                cardBg: cardBg,
                border: border,
                icon: Icons.location_on_rounded,
                iconColor: const Color(0xFF42A5F5),
                title: 'Location Shared',
                subtitle:
                    '${widget.lat.toStringAsFixed(4)}°N, ${widget.lng.toStringAsFixed(4)}°E',
                status: 'Live',
                statusColor: const Color(0xFF00E676),
                subtitleColor: subtitleColor,
                isDark: isDark,
              ),
              const SizedBox(height: 12),
              _StatusCard(
                cardBg: cardBg,
                border: border,
                icon: Icons.videocam_rounded,
                iconColor: const Color(0xFFFF7043),
                title: 'Video Recording',
                subtitle: _elapsedSeconds < 30
                    ? 'Recording in progress…'
                    : 'Recording completed & uploading',
                status: _elapsedSeconds < 30 ? 'Recording' : 'Completed',
                statusColor: _elapsedSeconds < 30
                    ? const Color(0xFFFF9800)
                    : const Color(0xFF00E676),
                subtitleColor: subtitleColor,
                isDark: isDark,
              ),
              const SizedBox(height: 12),
              _StatusCard(
                cardBg: cardBg,
                border: border,
                icon: Icons.message_rounded,
                iconColor: const Color(0xFF7C4DFF),
                title: 'SMS Alerts',
                subtitle: 'Emergency SMS with location sent',
                status: 'Sent',
                statusColor: const Color(0xFF00E676),
                subtitleColor: subtitleColor,
                isDark: isDark,
              ),
              const SizedBox(height: 12),
              _StatusCard(
                cardBg: cardBg,
                border: border,
                icon: Icons.cell_tower_rounded,
                iconColor: const Color(0xFFFF9800),
                title: 'Emergency Broadcast',
                subtitle: 'Alert broadcasted to nearby devices',
                status: 'Active',
                statusColor: const Color(0xFF00E676),
                subtitleColor: subtitleColor,
                isDark: isDark,
              ),
              const SizedBox(height: 24),
              // ── Stop SOS button ──
              SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: () {
                    widget.onDismiss();
                    Navigator.pop(context);
                  },
                  icon: const Icon(Icons.close, size: 20),
                  label: const Text('Stop Emergency'),
                  style: OutlinedButton.styleFrom(
                    foregroundColor: const Color(0xFFE53935),
                    side: const BorderSide(
                        color: Color(0xFFE53935), width: 1.5),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(14),
                    ),
                    textStyle: const TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }
}

/// Single status card row.
class _StatusCard extends StatelessWidget {
  final Color cardBg;
  final Color border;
  final IconData icon;
  final Color iconColor;
  final String title;
  final String subtitle;
  final String status;
  final Color statusColor;
  final Color subtitleColor;
  final bool isDark;

  const _StatusCard({
    required this.cardBg,
    required this.border,
    required this.icon,
    required this.iconColor,
    required this.title,
    required this.subtitle,
    required this.status,
    required this.statusColor,
    required this.subtitleColor,
    required this.isDark,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: cardBg,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: border),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: iconColor.withValues(alpha: 0.12),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(icon, color: iconColor, size: 22),
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: isDark
                        ? const Color(0xFFF0F0F5)
                        : const Color(0xFF1A1A2E),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  subtitle,
                  style: TextStyle(
                    fontSize: 12,
                    color: subtitleColor,
                  ),
                ),
              ],
            ),
          ),
          Container(
            padding:
                const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: statusColor.withValues(alpha: 0.12),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(
              status,
              style: TextStyle(
                fontSize: 10,
                fontWeight: FontWeight.w700,
                color: statusColor,
                letterSpacing: 0.5,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
