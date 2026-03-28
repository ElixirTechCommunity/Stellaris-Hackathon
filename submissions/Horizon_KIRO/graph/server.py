"""
Jarvis Graph Visualization Server
===================================
FastAPI server that:
  - Serves the live graph dashboard (index.html)
  - Accepts WebSocket connections from browsers
  - Receives agent events from the Jarvis backend
  - Broadcasts node-activation events to all connected browsers in real-time

Usage (started automatically by run_jarvis.py):
    from graph.server import start_viz_server
    start_viz_server()    # starts in a background thread on port 7860
"""

import asyncio
import json
import logging
import threading
from pathlib import Path
from typing import Set

import uvicorn
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles

logger = logging.getLogger(__name__)

# ──────────────────────────────────────────────
# FastAPI app
# ──────────────────────────────────────────────

app = FastAPI(title="Jarvis Live Graph", docs_url=None, redoc_url=None)

# ──────────────────────────────────────────────
# Connection manager — tracks open browser WS connections
# ──────────────────────────────────────────────

class ConnectionManager:
    def __init__(self):
        self.active: Set[WebSocket] = set()
        self._lock = asyncio.Lock()

    async def connect(self, ws: WebSocket):
        await ws.accept()
        async with self._lock:
            self.active.add(ws)
        logger.info(f"[viz] Browser connected — total: {len(self.active)}")

    async def disconnect(self, ws: WebSocket):
        async with self._lock:
            self.active.discard(ws)
        logger.info(f"[viz] Browser disconnected — total: {len(self.active)}")

    async def broadcast(self, message: dict):
        """Send JSON message to every connected browser."""
        dead = set()
        payload = json.dumps(message)
        async with self._lock:
            targets = set(self.active)
        for ws in targets:
            try:
                await ws.send_text(payload)
            except Exception:
                dead.add(ws)
        if dead:
            async with self._lock:
                self.active -= dead


manager = ConnectionManager()

# ──────────────────────────────────────────────
# Event loop reference — set when server starts
# ──────────────────────────────────────────────

_server_loop: asyncio.AbstractEventLoop | None = None


def broadcast_sync(message: dict):
    """
    Thread-safe broadcast called from background (non-async) threads.
    Schedules the coroutine onto the server's event loop.
    """
    global _server_loop
    if _server_loop and not _server_loop.is_closed():
        asyncio.run_coroutine_threadsafe(manager.broadcast(message), _server_loop)


# ──────────────────────────────────────────────
# HTTP routes
# ──────────────────────────────────────────────

STATIC_DIR = Path(__file__).parent / "static"


@app.get("/", response_class=HTMLResponse)
async def serve_dashboard():
    html_path = STATIC_DIR / "index.html"
    return HTMLResponse(content=html_path.read_text(encoding="utf-8"))


# ──────────────────────────────────────────────
# WebSocket endpoint
# ──────────────────────────────────────────────

@app.websocket("/ws")
async def websocket_endpoint(ws: WebSocket):
    await manager.connect(ws)
    # Send initial "ready" event so browser knows it's connected
    await ws.send_text(json.dumps({"type": "ready"}))
    try:
        while True:
            # Keep connection alive — wait for any client messages (ping/pong)
            data = await ws.receive_text()
            # Echo back as ack (not currently used by client)
            await ws.send_text(json.dumps({"type": "ack", "data": data}))
    except WebSocketDisconnect:
        await manager.disconnect(ws)
    except Exception as e:
        logger.warning(f"[viz] WebSocket error: {e}")
        await manager.disconnect(ws)


# ──────────────────────────────────────────────
# Server startup (call once, runs in background thread)
# ──────────────────────────────────────────────

_server_thread: threading.Thread | None = None
_server_started = threading.Event()


def start_viz_server(host: str = "127.0.0.1", port: int = 7860):
    """
    Start the FastAPI/uvicorn server in a daemon background thread.
    Safe to call multiple times — only starts once.
    """
    global _server_thread, _server_loop

    if _server_thread is not None and _server_thread.is_alive():
        logger.info("[viz] Server already running")
        return

    def _run():
        global _server_loop
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        _server_loop = loop

        config = uvicorn.Config(
            app,
            host=host,
            port=port,
            loop="asyncio",
            log_level="warning",   # keep terminal clean
            access_log=False,
        )
        server = uvicorn.Server(config)
        _server_started.set()
        loop.run_until_complete(server.serve())

    _server_thread = threading.Thread(target=_run, daemon=True, name="JarvisVizServer")
    _server_thread.start()
    _server_started.wait(timeout=5)
    logger.info(f"[viz] Dashboard live → http://{host}:{port}")
    print(f"\n🌐 Live graph dashboard → http://{host}:{port}\n")
