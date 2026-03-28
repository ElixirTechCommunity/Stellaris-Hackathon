"""
Jarvis LangGraph StateGraph — Full Architecture Visualization
==============================================================
Emits fine-grained node events for EVERY real agent and sub-agent in the
Jarvis system.  The browser dashboard renders the complete agent map and
highlights the exact execution path in real time.

Runs in a background daemon thread alongside the real agent (zero interference).
"""

import logging
import time
from typing import Optional, Any
from typing_extensions import TypedDict

from langgraph.graph import StateGraph, END

logger = logging.getLogger(__name__)


# ─────────────────────────────────────────────────────────────────────────────
# State
# ─────────────────────────────────────────────────────────────────────────────

class JarvisState(TypedDict):
    user_input: str
    intent: str
    goal: str
    confidence: float
    requires_safety_check: bool
    response: str
    error: Optional[str]
    _node_cb: Optional[Any]   # callable(node_id: str, delay=True)


# ─────────────────────────────────────────────────────────────────────────────
# Emit helper
# ─────────────────────────────────────────────────────────────────────────────

def _emit(state: JarvisState, node_id: str, pause: float = 0.25):
    """Fire broadcast callback then pause so browser can animate."""
    cb = state.get("_node_cb")
    if cb:
        try:
            cb(node_id)
        except Exception:
            pass
    time.sleep(pause)


# ─────────────────────────────────────────────────────────────────────────────
# Intent constants
# ─────────────────────────────────────────────────────────────────────────────

CONVERSATIONAL_INTENTS = {"unknown", "help", "conversational", ""}
SAFETY_INTENTS         = {"system_control", "maintenance", "installer", "file_management"}

# intent → specialized agent node ID
INTENT_AGENT_MAP = {
    "file_management":   "file_manager",
    "system_control":    "system_control",
    "diagnostics":       "diagnostics_agent",
    "disk_analysis":     "disk_analysis_agent",
    "maintenance":       "maintenance_agent",
    "health_check":      "health_check_agent",
    "installer":         "installer_agent",
    "web_search":        "web_search_agent",
    "open_app":          "app_launcher_agent",
    "vision_analysis":   "vision_agent",
    "troubleshoot_screen":"vision_agent",
    "personalization":   "personalization_agent",
    "system_config":     "system_config_agent",
}


def _pick_sub_agent(intent: str, text: str) -> Optional[str]:
    """Keyword-based heuristic to choose the most likely sub-agent."""
    t = text.lower()
    if intent == "file_management":
        if any(w in t for w in ["duplicate", "dupli", "copy", "same"]):
            return "duplicate_detector"
        if any(w in t for w in ["organize", "sort", "arrange", "type", "extension"]):
            return "file_organizer"
        if any(w in t for w in ["large", "size", "big", "heavy", "huge", "space"]):
            return "large_file_scanner"
        if any(w in t for w in ["find", "locate", "search", "folder"]):
            return "folder_finder"
        return "file_organizer"

    if intent == "diagnostics":
        if any(w in t for w in ["cpu", "processor", "core", "speed"]):
            return "cpu_monitor"
        if any(w in t for w in ["ram", "memory", "gb", "mb"]):
            return "ram_monitor"
        if any(w in t for w in ["disk", "storage", "drive", "c:"]):
            return "disk_monitor"
        return "system_health_checker"

    if intent == "installer":
        if any(w in t for w in ["wallpaper", "image", "background", "picture"]):
            return "resource_downloader"
        return "winget_engine"

    if intent in ("vision_analysis", "troubleshoot_screen"):
        return "screenshot_tool"

    if intent == "maintenance":
        if any(w in t for w in ["old", "stale", "ancient", "unused"]):
            return "old_file_scanner"
        return "temp_file_cleaner"

    if intent == "disk_analysis":
        if any(w in t for w in ["folder", "directory", "which"]):
            return "large_folder_finder"
        return "disk_analyzer"

    return None


# ─────────────────────────────────────────────────────────────────────────────
# LangGraph nodes
# ─────────────────────────────────────────────────────────────────────────────

def node_receive_input(state: JarvisState) -> JarvisState:
    """Entry: simulate input arriving through UI."""
    # Show the input pipeline activating
    _emit(state, "voice_agent",   0.18)
    _emit(state, "widget",        0.18)
    _emit(state, "master_agent",  0.20)
    return {**state, "intent": "", "error": None}


def node_extract_intent(state: JarvisState) -> JarvisState:
    """Call Groq to extract intent + goal, animate LLM usage."""
    _emit(state, "groq_llm", 0.15)   # LLM activates first
    try:
        from utils.json_parser import extract_json
        from llm.groq_client import GroqClient

        prompt = (
            "Classify the following OS assistant command. "
            "Return ONLY JSON: "
            '{"intent": "<open_app|web_search|system_control|diagnostics|'
            "disk_analysis|maintenance|health_check|troubleshoot_screen|vision_analysis|"
            'personalization|file_management|system_config|installer|help|unknown>", '
            '"goal": "<short goal>", "confidence": 0.0-1.0, "requires_safety_check": true/false}\n\n'
            f"Command: {state['user_input']}"
        )
        raw = GroqClient().generate(prompt)
        info = extract_json(raw)

        if isinstance(info, dict):
            return {
                **state,
                "intent":               info.get("intent", "unknown"),
                "goal":                 info.get("goal", state["user_input"]),
                "confidence":           float(info.get("confidence", 0.5)),
                "requires_safety_check":bool(info.get("requires_safety_check", False)),
            }
    except Exception as e:
        logger.warning(f"[graph] extract_intent: {e}")

    return {**state, "intent": "unknown", "goal": state["user_input"], "confidence": 0.0}


def node_safety_check(state: JarvisState) -> JarvisState:
    """Safety validation gate."""
    _emit(state, "safety_check", 0.30)
    return state


def node_task_plan(state: JarvisState) -> JarvisState:
    """Simulate task decomposition via Task Planner + Groq."""
    _emit(state, "task_planner", 0.20)
    _emit(state, "groq_llm",     0.20)
    _emit(state, "intent_router",0.20)
    return state


def node_execute_tasks(state: JarvisState) -> JarvisState:
    """
    Activate Execution Engine, then the specific specialized agent,
    then the relevant sub-agent(s) — all based on detected intent.
    """
    _emit(state, "execution_engine", 0.25)

    intent     = state.get("intent", "unknown")
    user_input = state.get("user_input", "")

    # Fire specialized agent
    agent_id = INTENT_AGENT_MAP.get(intent)
    if agent_id:
        _emit(state, agent_id, 0.25)

        # Fire sub-agent
        sub = _pick_sub_agent(intent, user_input)
        if sub:
            _emit(state, sub, 0.25)

        # Vision also needs Gemini + extra sub-agents
        if intent in ("vision_analysis", "troubleshoot_screen"):
            _emit(state, "vision_analyzer",  0.20)
            _emit(state, "gemini_llm",        0.20)
            _emit(state, "auto_fixer",        0.20)
            _emit(state, "solution_parser",   0.15)

    return state


def node_save_memory(state: JarvisState) -> JarvisState:
    """Persist execution in Memory Manager."""
    _emit(state, "memory_manager", 0.20)
    return state


def node_conversational_fallback(state: JarvisState) -> JarvisState:
    """Direct LLM reply for unknown/casual inputs."""
    _emit(state, "conversational_fallback", 0.20)
    _emit(state, "groq_llm",               0.25)
    return {**state, "response": "(conversational reply via main agent)"}


def node_format_response(state: JarvisState) -> JarvisState:
    """Final response formatting."""
    _emit(state, "format_response", 0.15)
    if not state.get("response"):
        return {**state, "response": f"Completed: {state.get('goal', '')}"}
    return state


# ─────────────────────────────────────────────────────────────────────────────
# Routing
# ─────────────────────────────────────────────────────────────────────────────

def route_after_intent(state: JarvisState) -> str:
    intent = state.get("intent", "unknown")
    if intent in CONVERSATIONAL_INTENTS:
        return "conversational_fallback"
    if intent in SAFETY_INTENTS or state.get("requires_safety_check"):
        return "safety_check"
    return "task_plan"


def route_after_safety(state: JarvisState) -> str:
    return "conversational_fallback" if state.get("error") else "task_plan"


# ─────────────────────────────────────────────────────────────────────────────
# Graph builder
# ─────────────────────────────────────────────────────────────────────────────

def build_graph():
    g = StateGraph(JarvisState)

    g.add_node("receive_input",           node_receive_input)
    g.add_node("extract_intent",          node_extract_intent)
    g.add_node("safety_check",            node_safety_check)
    g.add_node("task_plan",               node_task_plan)
    g.add_node("execute_tasks",           node_execute_tasks)
    g.add_node("save_memory",             node_save_memory)
    g.add_node("conversational_fallback", node_conversational_fallback)
    g.add_node("format_response",         node_format_response)

    g.set_entry_point("receive_input")
    g.add_edge("receive_input", "extract_intent")

    g.add_conditional_edges("extract_intent", route_after_intent, {
        "safety_check":            "safety_check",
        "task_plan":               "task_plan",
        "conversational_fallback": "conversational_fallback",
    })
    g.add_conditional_edges("safety_check", route_after_safety, {
        "task_plan":               "task_plan",
        "conversational_fallback": "conversational_fallback",
    })

    g.add_edge("task_plan",               "execute_tasks")
    g.add_edge("execute_tasks",           "save_memory")
    g.add_edge("save_memory",             "format_response")
    g.add_edge("conversational_fallback", "format_response")
    g.add_edge("format_response",         END)

    return g.compile()


_compiled_graph = None

def get_graph():
    global _compiled_graph
    if _compiled_graph is None:
        _compiled_graph = build_graph()
    return _compiled_graph


# ─────────────────────────────────────────────────────────────────────────────
# Public API
# ─────────────────────────────────────────────────────────────────────────────

def run_viz(user_input: str, node_callback=None) -> dict:
    """
    Run the lightweight visualization graph.
    Fires node_callback(node_id) for every agent that activates.
    """
    graph = get_graph()
    initial: JarvisState = {
        "user_input":            user_input,
        "intent":                "",
        "goal":                  "",
        "confidence":            0.0,
        "requires_safety_check": False,
        "response":              "",
        "error":                 None,
        "_node_cb":              node_callback,
    }
    return graph.invoke(initial)
