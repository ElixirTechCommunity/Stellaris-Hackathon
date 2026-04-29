import 'dart:io';
import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

/// Simple full-screen video player for local SOS recordings.
class VideoPlayerScreen extends StatefulWidget {
  final String filePath;
  final String title;

  const VideoPlayerScreen({
    super.key,
    required this.filePath,
    this.title = 'SOS Recording',
  });

  @override
  State<VideoPlayerScreen> createState() => _VideoPlayerScreenState();
}

class _VideoPlayerScreenState extends State<VideoPlayerScreen> {
  late VideoPlayerController _controller;
  bool _initialized = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _initPlayer();
  }

  Future<void> _initPlayer() async {
    try {
      final file = File(widget.filePath);
      if (!await file.exists()) {
        setState(() => _error = 'Recording file not found');
        return;
      }

      _controller = VideoPlayerController.file(file);
      await _controller.initialize();
      await _controller.play();
      if (mounted) setState(() => _initialized = true);
    } catch (e) {
      if (mounted) setState(() => _error = 'Could not play video: $e');
    }
  }

  @override
  void dispose() {
    if (_initialized) _controller.dispose();
    super.dispose();
  }

  String _formatDuration(Duration d) {
    final mins = d.inMinutes.remainder(60).toString().padLeft(2, '0');
    final secs = d.inSeconds.remainder(60).toString().padLeft(2, '0');
    return '$mins:$secs';
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        title: Text(widget.title,
            style: const TextStyle(color: Colors.white, fontSize: 16)),
        iconTheme: const IconThemeData(color: Colors.white),
      ),
      body: _error != null
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.error_outline, color: Colors.red, size: 48),
                  const SizedBox(height: 12),
                  Text(_error!,
                      style: const TextStyle(color: Colors.white70, fontSize: 14)),
                ],
              ),
            )
          : !_initialized
              ? const Center(
                  child: CircularProgressIndicator(color: Color(0xFFE53935)))
              : Column(
                  children: [
                    Expanded(
                      child: Center(
                        child: AspectRatio(
                          aspectRatio: _controller.value.aspectRatio,
                          child: VideoPlayer(_controller),
                        ),
                      ),
                    ),
                    // Controls
                    Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 20, vertical: 16),
                      color: Colors.black,
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          // Progress bar
                          ValueListenableBuilder<VideoPlayerValue>(
                            valueListenable: _controller,
                            builder: (_, value, __) {
                              return Column(
                                children: [
                                  SliderTheme(
                                    data: SliderThemeData(
                                      thumbShape: const RoundSliderThumbShape(
                                          enabledThumbRadius: 6),
                                      overlayShape:
                                          const RoundSliderOverlayShape(
                                              overlayRadius: 14),
                                      activeTrackColor:
                                          const Color(0xFFE53935),
                                      inactiveTrackColor: Colors.white24,
                                      thumbColor: const Color(0xFFE53935),
                                    ),
                                    child: Slider(
                                      value: value.position.inMilliseconds
                                          .toDouble(),
                                      max: value.duration.inMilliseconds
                                          .toDouble()
                                          .clamp(1, double.infinity),
                                      onChanged: (v) {
                                        _controller.seekTo(Duration(
                                            milliseconds: v.toInt()));
                                      },
                                    ),
                                  ),
                                  Padding(
                                    padding: const EdgeInsets.symmetric(
                                        horizontal: 16),
                                    child: Row(
                                      mainAxisAlignment:
                                          MainAxisAlignment.spaceBetween,
                                      children: [
                                        Text(
                                          _formatDuration(value.position),
                                          style: const TextStyle(
                                              color: Colors.white60,
                                              fontSize: 12),
                                        ),
                                        Text(
                                          _formatDuration(value.duration),
                                          style: const TextStyle(
                                              color: Colors.white60,
                                              fontSize: 12),
                                        ),
                                      ],
                                    ),
                                  ),
                                ],
                              );
                            },
                          ),
                          const SizedBox(height: 8),
                          // Play/Pause button
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              IconButton(
                                onPressed: () {
                                  _controller.seekTo(
                                      _controller.value.position -
                                          const Duration(seconds: 5));
                                },
                                icon: const Icon(Icons.replay_5,
                                    color: Colors.white, size: 32),
                              ),
                              const SizedBox(width: 16),
                              ValueListenableBuilder<VideoPlayerValue>(
                                valueListenable: _controller,
                                builder: (_, value, __) {
                                  return GestureDetector(
                                    onTap: () {
                                      value.isPlaying
                                          ? _controller.pause()
                                          : _controller.play();
                                    },
                                    child: Container(
                                      width: 56,
                                      height: 56,
                                      decoration: const BoxDecoration(
                                        color: Color(0xFFE53935),
                                        shape: BoxShape.circle,
                                      ),
                                      child: Icon(
                                        value.isPlaying
                                            ? Icons.pause
                                            : Icons.play_arrow,
                                        color: Colors.white,
                                        size: 32,
                                      ),
                                    ),
                                  );
                                },
                              ),
                              const SizedBox(width: 16),
                              IconButton(
                                onPressed: () {
                                  _controller.seekTo(
                                      _controller.value.position +
                                          const Duration(seconds: 5));
                                },
                                icon: const Icon(Icons.forward_5,
                                    color: Colors.white, size: 32),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
    );
  }
}
