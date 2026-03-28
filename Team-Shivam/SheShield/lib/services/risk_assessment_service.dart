import 'dart:math';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:http/http.dart' as http;

// ═══════════════════════════════════════════════════════════════════
//  Risk Assessment Service  –  OFFLINE-FIRST + ONLINE-ENHANCED
//
//  Offline mode (default):
//    • Time-based check only (night = risky)
//    • Uses hardcoded safe-zone list for nearest safe zone display
//
//  Online mode (when internet is available):
//    • Fetches nearby places via Overpass API
//    • If <= 2 places found nearby → isolated area
//    • Falls back to offline mode on API failure
// ═══════════════════════════════════════════════════════════════════

/// Risk severity levels — ordered from safest to most dangerous.
enum RiskLevel { safe, moderate, high, critical }

/// Whether online-enhanced detection was used.
enum DetectionMode { offline, online }

/// Info about the closest safe zone to the user.
class SafeZoneInfo {
  final String name;
  final String type;
  final double distanceKm;
  final double lat;
  final double lng;

  const SafeZoneInfo({
    required this.name,
    required this.type,
    required this.distanceKm,
    required this.lat,
    required this.lng,
  });

  String get displayText => '$name (${distanceKm.toStringAsFixed(1)} km)';
}

/// Holds the result of a risk assessment.
class RiskResult {
  final RiskLevel level;
  final String message;
  final String subtitle;
  final List<String> reasons;
  final SafeZoneInfo? nearestSafeZone;
  final DetectionMode mode;
  final int? nearbyPlacesCount;
  final int riskScore; // 0=safe, 1=moderate, 2=high, 3=critical

  const RiskResult({
    required this.level,
    required this.message,
    required this.subtitle,
    required this.reasons,
    this.nearestSafeZone,
    required this.mode,
    this.nearbyPlacesCount,
    required this.riskScore,
  });
}

/// Pure-logic service — offline-first with online enhancement.
class RiskAssessmentService {
  // ────────────────────────────────────────────────────────────────
  //  Demo toggle — can still be flipped from the UI switch.
  // ────────────────────────────────────────────────────────────────
  static bool isIsolatedArea = false;

  // ────────────────────────────────────────────────────────────────
  //  Named safe zones (offline fallback for nearest safe zone)
  // ────────────────────────────────────────────────────────────────
  static const List<_SafePlace> _safePlaces = [
    _SafePlace('Parliament St Police Station', 'Police Station', 28.6218, 77.2120),
    _SafePlace('AIIMS Hospital', 'Hospital', 28.5672, 77.2100),
    _SafePlace('Connaught Place', 'Public Area', 28.6315, 77.2167),
    _SafePlace('Safdarjung Hospital', 'Hospital', 28.5685, 77.2066),
    _SafePlace('Noida Sec 18 Market', 'Public Area', 28.5706, 77.3218),
    _SafePlace('CST Railway Station', 'Public Area', 19.0826, 72.8853),
    _SafePlace('JJ Hospital Mumbai', 'Hospital', 18.9634, 72.8330),
    _SafePlace('Colaba Police Station', 'Police Station', 18.9220, 72.8347),
    _SafePlace('MG Road Metro', 'Public Area', 12.9758, 77.6065),
    _SafePlace('Victoria Hospital', 'Hospital', 12.9563, 77.5737),
    _SafePlace('Cubbon Park Police', 'Police Station', 12.9763, 77.5929),
    _SafePlace('Chennai Central', 'Public Area', 13.0836, 80.2750),
    _SafePlace('Rajiv Gandhi Hospital', 'Hospital', 13.0700, 80.2800),
    _SafePlace('Park Street', 'Public Area', 22.5534, 88.3535),
    _SafePlace('SSKM Hospital Kolkata', 'Hospital', 22.5383, 88.3436),
    _SafePlace('HITEC City', 'Public Area', 17.4435, 78.3772),
    _SafePlace('Nizam\'s Hospital', 'Hospital', 17.3887, 78.4741),
    _SafePlace('Hawa Mahal Area', 'Public Area', 26.9239, 75.8267),
    _SafePlace('SMS Hospital Jaipur', 'Hospital', 26.9038, 75.8013),
    _SafePlace('SG Highway', 'Public Area', 23.0300, 72.5300),
    _SafePlace('Civil Hospital Ahmedabad', 'Hospital', 23.0500, 72.6000),
    _SafePlace('Cyber Hub Gurgaon', 'Public Area', 28.4949, 77.0887),
    _SafePlace('Medanta Hospital', 'Hospital', 28.4395, 77.0425),
  ];

  /// Minimum places needed to consider an area "populated".
  static const int _minPlacesForSafe = 5;

  // ────────────────────────────────────────────────────────────────
  //  Public helpers
  // ────────────────────────────────────────────────────────────────

  static bool isNightTime() =>
      DateTime.now().hour >= 21 || DateTime.now().hour < 5;

  /// Check if the device has internet connectivity.
  static Future<bool> _hasInternet() async {
    try {
      final result = await Connectivity().checkConnectivity();
      return !result.contains(ConnectivityResult.none);
    } catch (_) {
      return false;
    }
  }

  // ────────────────────────────────────────────────────────────────
  //  Main entry-point (async — checks internet first)
  // ────────────────────────────────────────────────────────────────

  /// Assess risk with online enhancement when internet is available.
  /// Falls back to offline-only when no internet or API fails.
  static Future<RiskResult> assessRiskAsync(double lat, double lng) async {
    final online = await _hasInternet();

    if (online) {
      try {
        return await _assessOnline(lat, lng);
      } catch (e) {
        debugPrint('RiskAssessment: Online check failed, falling back — $e');
        return _assessOffline(lat, lng);
      }
    } else {
      return _assessOffline(lat, lng);
    }
  }

  /// Synchronous offline-only assessment (instant, no network).
  static RiskResult assessRisk(double lat, double lng) {
    return _assessOffline(lat, lng);
  }

  // ────────────────────────────────────────────────────────────────
  //  Online assessment (Overpass API)
  // ────────────────────────────────────────────────────────────────

  static Future<RiskResult> _assessOnline(double lat, double lng) async {
    final reasons = <String>[];
    final hour = DateTime.now().hour;
    final nightTime = isNightTime();

    if (nightTime) {
      reasons.add('It\'s late night (${_formatHour(hour)})');
    }

    // Fetch nearby places (amenities) within 1000m radius via Overpass
    final placesCount = await _fetchNearbyPlacesCount(lat, lng);
    final isolated = isIsolatedArea || placesCount < _minPlacesForSafe;

    // Debug logs
    debugPrint('🔍 [ONLINE] Places count: $placesCount');
    debugPrint('🔍 [ONLINE] isIsolatedArea (demo toggle): $isIsolatedArea');
    debugPrint('🔍 [ONLINE] isolated (final): $isolated');

    if (isolated) {
      reasons.add('Only $placesCount places found nearby');
    }

    // Calculate risk score
    int riskScore = 0;
    if (nightTime) riskScore += 1;
    if (isolated) riskScore += 2;

    debugPrint('🔍 [ONLINE] riskScore: $riskScore (night=$nightTime, isolated=$isolated)');

    // Find nearest safe zone from our hardcoded list
    final nearest = _findNearestSafePlace(lat, lng);

    return _buildResult(
      riskScore: riskScore,
      reasons: reasons,
      nearest: nearest,
      mode: DetectionMode.online,
      placesCount: placesCount,
    );
  }

  /// Query Overpass API for nearby amenities within 1000m.
  /// Returns the count of results found.
  static Future<int> _fetchNearbyPlacesCount(double lat, double lng) async {
    // Search for ANY amenity, shop, or leisure node within 1000m
    // No type filter — broad query returns more results in crowded areas
    final query = '[out:json][timeout:8];'
        '('
        'node["amenity"](around:1000,$lat,$lng);'
        'node["shop"](around:1000,$lat,$lng);'
        'node["leisure"](around:1000,$lat,$lng);'
        ')'
        ';out count;';

    try {
      final response = await http
          .post(
            Uri.parse('https://overpass-api.de/api/interpreter'),
            body: {'data': query},
          )
          .timeout(const Duration(seconds: 8));

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        final elements = data['elements'] as List<dynamic>? ?? [];
        // "out count" returns a single element with a "tags" map containing "total"
        if (elements.isNotEmpty && elements[0]['tags'] != null) {
          final total = int.tryParse(
              elements[0]['tags']['total']?.toString() ?? '0');
          debugPrint('🔍 Overpass API returned: $total places within 1000m');
          return total ?? 0;
        }
        debugPrint('🔍 Overpass API returned: ${elements.length} elements');
        return elements.length;
      }
      debugPrint('🔍 Overpass API status: ${response.statusCode}');
    } catch (e) {
      debugPrint('🔍 Overpass query failed — $e');
    }
    // FAIL-SAFE: On failure, assume safe (return high count to avoid false alerts)
    debugPrint('🔍 API failed → assuming safe (returning 99)');
    return 99;
  }

  // ────────────────────────────────────────────────────────────────
  //  Offline assessment (local data only, instant)
  // ────────────────────────────────────────────────────────────────

  static RiskResult _assessOffline(double lat, double lng) {
    final reasons = <String>[];
    final hour = DateTime.now().hour;
    final nightTime = isNightTime();

    if (nightTime) {
      reasons.add('It\'s late night (${_formatHour(hour)})');
    }

    // In offline mode, use demo toggle or GPS distance for isolation
    // Use 5km threshold (more lenient to avoid false positives)
    final nearestKm = _distanceToNearestSafeZone(lat, lng);
    final isolated = isIsolatedArea || nearestKm > 5.0;

    // Debug logs
    debugPrint('🔍 [OFFLINE] Nearest safe zone: ${nearestKm.toStringAsFixed(1)} km');
    debugPrint('🔍 [OFFLINE] isIsolatedArea (demo toggle): $isIsolatedArea');
    debugPrint('🔍 [OFFLINE] isolated (final): $isolated');

    if (isolated) {
      reasons.add('${nearestKm.toStringAsFixed(1)} km from nearest safe zone');
    }

    // Calculate risk score
    int riskScore = 0;
    if (nightTime) riskScore += 1;
    if (isolated) riskScore += 2;

    debugPrint('🔍 [OFFLINE] riskScore: $riskScore (night=$nightTime, isolated=$isolated)');

    final nearest = _findNearestSafePlace(lat, lng);

    return _buildResult(
      riskScore: riskScore,
      reasons: reasons,
      nearest: nearest,
      mode: DetectionMode.offline,
    );
  }

  // ────────────────────────────────────────────────────────────────
  //  Shared result builder
  // ────────────────────────────────────────────────────────────────

  static RiskResult _buildResult({
    required int riskScore,
    required List<String> reasons,
    SafeZoneInfo? nearest,
    required DetectionMode mode,
    int? placesCount,
  }) {
    // Score 0 = safe, 1 = moderate, 2 = high, 3 = critical
    if (riskScore >= 3) {
      return RiskResult(
        level: RiskLevel.critical,
        message: '🚨 Critical Risk: Immediate Action Needed',
        subtitle: 'Late night + low activity area — stay alert',
        reasons: reasons,
        nearestSafeZone: nearest,
        mode: mode,
        nearbyPlacesCount: placesCount,
        riskScore: riskScore,
      );
    } else if (riskScore == 2) {
      return RiskResult(
        level: RiskLevel.high,
        message: '🚨 High Risk: Low Activity Area',
        subtitle: 'You appear to be in an isolated location',
        reasons: reasons,
        nearestSafeZone: nearest,
        mode: mode,
        nearbyPlacesCount: placesCount,
        riskScore: riskScore,
      );
    } else if (riskScore == 1) {
      return RiskResult(
        level: RiskLevel.moderate,
        message: '⚠️ Be Alert: Late Night',
        subtitle: 'Late-night hours — stay in well-lit areas',
        reasons: reasons,
        nearestSafeZone: nearest,
        mode: mode,
        nearbyPlacesCount: placesCount,
        riskScore: riskScore,
      );
    } else {
      return RiskResult(
        level: RiskLevel.safe,
        message: '✅ You are in a safe area',
        subtitle: 'No risk factors detected right now',
        reasons: reasons,
        nearestSafeZone: nearest,
        mode: mode,
        nearbyPlacesCount: placesCount,
        riskScore: riskScore,
      );
    }
  }

  // ────────────────────────────────────────────────────────────────
  //  Helpers
  // ────────────────────────────────────────────────────────────────

  static SafeZoneInfo? _findNearestSafePlace(double lat, double lng) {
    if (_safePlaces.isEmpty) return null;
    _SafePlace? best;
    double bestDist = double.infinity;
    for (final p in _safePlaces) {
      final d = _haversineKm(lat, lng, p.lat, p.lng);
      if (d < bestDist) {
        bestDist = d;
        best = p;
      }
    }
    if (best == null) return null;
    return SafeZoneInfo(
      name: best.name,
      type: best.type,
      distanceKm: bestDist,
      lat: best.lat,
      lng: best.lng,
    );
  }

  static double _distanceToNearestSafeZone(double lat, double lng) {
    double nearest = double.infinity;
    for (final zone in _safePlaces) {
      final d = _haversineKm(lat, lng, zone.lat, zone.lng);
      if (d < nearest) nearest = d;
    }
    return nearest;
  }

  static double _haversineKm(
    double lat1, double lng1,
    double lat2, double lng2,
  ) {
    const earthRadiusKm = 6371.0;
    final dLat = _degToRad(lat2 - lat1);
    final dLng = _degToRad(lng2 - lng1);
    final a = sin(dLat / 2) * sin(dLat / 2) +
        cos(_degToRad(lat1)) *
            cos(_degToRad(lat2)) *
            sin(dLng / 2) *
            sin(dLng / 2);
    final c = 2 * atan2(sqrt(a), sqrt(1 - a));
    return earthRadiusKm * c;
  }

  static double _degToRad(double deg) => deg * (pi / 180);

  static String _formatHour(int hour) {
    final h = hour % 12 == 0 ? 12 : hour % 12;
    final period = hour >= 12 ? 'PM' : 'AM';
    return '$h $period';
  }
}

class _SafePlace {
  final String name;
  final String type;
  final double lat;
  final double lng;
  const _SafePlace(this.name, this.type, this.lat, this.lng);
}
