# 🤖 JARVIS - Complete Workflow & Architecture Guide

> **Intelligent OS Automation Agent - From Input to Output**

This document explains how JARVIS processes user inputs, plans tasks, executes them, and returns results. It covers the complete workflow with concrete examples.

---

## 📑 Table of Contents

1. [Quick Overview](#quick-overview)
2. [Architecture Diagram](#architecture-diagram)
3. [Input Entry Points](#input-entry-points)
4. [Complete Workflow: Step-by-Step](#complete-workflow-step-by-step)
5. [Core Components](#core-components)
6. [Concrete Examples](#concrete-examples)
7. [File Structure Reference](#file-structure-reference)
8. [How to Use JARVIS](#how-to-use-jarvis)
9. [Key Concepts](#key-concepts)

---

## 🚀 Quick Overview

JARVIS is an **AI-powered task automation system** that converts natural language commands into executable system tasks. 

**The Flow:**
```
USER INPUT → BRAIN (MasterAgent) → PLANNING (TaskPlanner) → EXECUTION (ExecutionEngine) 
→ ROUTING (IntentRouter) → TOOLS → OUTPUT → MEMORY (Storage)
```

---

## 📊 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                            INPUT LAYER                                   │
├─────────────────────────────────────────────────────────────────────────┤
│  Text Input (CLI) │  Voice Input (Mic) │  GUI Input (Chat Card)         │
│  main.py          │  voice_listener.py  │  ui/app.py                    │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                      BRAIN LAYER (Orchestration)                         │
├─────────────────────────────────────────────────────────────────────────┤
│              agent/master_agent.py (Central Hub)                         │
│  • Receives user input                                                   │
│  • Extracts intent using LLM                                             │
│  • Coordinates all components                                            │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                    REASONING LAYER (AI & Planning)                       │
├─────────────────────────────────────────────────────────────────────────┤
│  llm/groq_client.py         │  agent/task_planner.py                    │
│  (AI Reasoning)             │  (Task Decomposition)                      │
│  • Understand intent        │  • Break down goals                        │
│  • Extract goals            │  • Create task dependencies                │
│  • Generate solutions       │  • Assign priorities                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                    EXECUTION LAYER (Task Running)                        │
├─────────────────────────────────────────────────────────────────────────┤
│  agent/execution_engine.py                                              │
│  • Validates tasks                                                       │
│  • Executes sequentially                                                 │
│  • Handles errors & retries                                              │
│  • Reports results                                                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                     ROUTING LAYER (Intent Router)                        │
├─────────────────────────────────────────────────────────────────────────┤
│              router/intent_router.py (Dispatcher)                        │
│  • Validates intent & action                                             │
│  • Routes to correct handler                                             │
│  • Maps parameters                                                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
     ┌──────────────────────────────┼──────────────────────────────┐
     ↓                              ↓                              ↓
┌─────────────────────────┐  ┌──────────────────────┐  ┌──────────────────────┐
│   SYSTEM TOOLS LAYER    │  │  FILE MANAGER LAYER  │  │  VISION/DIAGNOSTICS  │
├─────────────────────────┤  ├──────────────────────┤  ├──────────────────────┤
│ tools/                  │  │ file_manager/        │  │ troubleshooter/      │
│ • system_tools.py       │  │ • file_organizer.py  │  │ • vision_analyzer.py │
│ • app_tools.py          │  │ • duplicate_detector │  │ • screenshot_tool.py │
│ • diagnostics_tools.py  │  │ • large_file_scanner │  │ vision/              │
│ • system_config.py      │  │ • llm_organizer.py   │  │ • vision_engine.py   │
│ • web_tools.py          │  │ • file_organizer.py  │  │ personalisation/     │
│                         │  │                      │  │ • personalisation_..  │
│ Functions:              │  │ Functions:           │  │ installer/           │
│ • sleep_pc()            │  │ • organize_files()   │  │ • installer.py       │
│ • lock_pc()             │  │ • find_duplicates()  │  │                      │
│ • kill_process()        │  │ • remove_duplicates()│  │ Functions:           │
│ • clean_temp()          │  │ • analyze_folder()   │  │ • capture_screen()   │
│ • get_cpu_usage()       │  │ • get_report()       │  │ • analyze_image()    │
│ • get_ram_usage()       │  │                      │  │ • install_app()      │
│ • get_disk_usage()      │  │                      │  │                      │
└─────────────────────────┘  └──────────────────────┘  └──────────────────────┘
     ↓                              ↓                              ↓
     └──────────────────────────────┼──────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                         OUTPUT LAYER                                     │
├─────────────────────────────────────────────────────────────────────────┤
│  Terminal Output │  GUI Display │  Chat Card │  Log Files               │
└─────────────────────────────────────────────────────────────────────────┘
                                    ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                       MEMORY LAYER (Persistence)                         │
├─────────────────────────────────────────────────────────────────────────┤
│  agent/memory_manager.py                                                │
│  agent_memory/                                                           │
│  • execution_history.json                                                │
│  • context.json                                                          │
│  • plans.json                                                            │
│  • preferences.json                                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 📥 Input Entry Points

JARVIS supports **3 modes of input**:

### 1️⃣ **CLI Text Mode** (Terminal-based)
```bash
python main.py
```
- **File:** `main.py`
- **Input Method:** Keyboard (interactive terminal)
- **Best For:** Server/headless environments, scripting
- **Flow:**
  ```
  User types → main.py reads input() → passes to MasterAgent
  ```

### 2️⃣ **GUI Chat Mode** (Desktop UI)
```bash
python run_jarvis.py
```
- **Files:** `ui/app.py`, `ui/chat_card.py`, `widget/widget.py`
- **Input Method:** Chat card UI with text input
- **Best For:** Desktop users, comfortable UI
- **Flow:**
  ```
  User types in chat → ChatCard captures → MasterAgent processes
  ```

### 3️⃣ **GUI + Voice Mode** (Desktop + Microphone)
```bash
python run_jarvis.py --voice
```
- **Files:** `voice/voice_listener.py`, `ui/app.py`
- **Input Method:** Microphone (Groq Whisper API for transcription)
- **Best For:** Hands-free operation, accessibility
- **Flow:**
  ```
  User speaks → VoiceListener records audio → Groq Whisper transcribes 
  → Text sent to MasterAgent
  ```

---

## 🔄 Complete Workflow: Step-by-Step

### **STEP 1: INPUT CAPTURE**
```
┌─ User provides input (text/voice/GUI)
└─ Input is captured by entry point
   • CLI: input() function
   • GUI: Chat card captures text
   • Voice: VoiceListener records & transcribes
```

### **STEP 2: MASTER AGENT - GOAL EXTRACTION**
**File:** `agent/master_agent.py` → `process_command()`

```python
# Input: "Find duplicate files in my Downloads folder"

# Step 2a: Extract Goal using LLM
goal_info = _extract_goal(user_input)

# Output:
{
    "intent": "file_management",
    "goal": "find duplicate files in Downloads",
    "confidence": 0.95,
    "is_complex": True,
    "requires_safety_check": False
}

# Step 2b: Display to user
🎯 Goal: find duplicate files in Downloads
🔍 Confidence: 95%
```

**Key Details:**
- Uses `llm/groq_client.py` to call LLM
- Extracts structured JSON from response
- Confidence score helps validate understanding
- Marks if goal is complex or needs safety check

---

### **STEP 3: TASK PLANNER - DECOMPOSITION**
**File:** `agent/task_planner.py` → `plan_goal()`

```python
# Input: Goal string and user input
# Output: TaskPlan with array of tasks

goal = "find duplicate files in Downloads"
plan = planner.plan_goal(user_input, goal)

# LLM Decomposition generates:
[
  {
    "order": 1,
    "intent": "file_management",
    "action": "find_duplicates",
    "parameters": {"folder_name": "Downloads"},
    "description": "Scan for duplicate files in Downloads",
    "priority": "high",
    "depends_on": [],
    "requires_confirmation": False,
    "dry_run": False
  },
  {
    "order": 2,
    "intent": "file_management",
    "action": "remove_duplicates",
    "parameters": {"folder_name": "Downloads"},
    "description": "Remove detected duplicate files",
    "priority": "high",
    "depends_on": ["task_1"],
    "requires_confirmation": True,  # ⚠️ Safety check needed
    "dry_run": False
  }
]

# Display:
📄 Created plan with 2 tasks:
  1. Scan for duplicates
  2. Remove duplicates (needs confirmation)
```

**Key Features:**
- Breaks complex goals into atomic tasks
- Tracks dependencies (task2 depends on task1)
- Marks tasks needing user confirmation
- Assigns priorities
- Uses LLM for intelligent reasoning

---

### **STEP 4: EXECUTION ENGINE - RUN TASKS**
**File:** `agent/execution_engine.py` → `execute_plan()`

```python
# Input: TaskPlan
# Process: Execute each task sequentially

execution_summary = executor.execute_plan(plan, auto_confirm=False)

# For each task:
for task in plan.tasks:
    # 1. Check dependencies
    if task dependencies_failed:
        task.mark_skipped()
        continue
    
    # 2. Ask for confirmation if needed
    if task.requires_confirmation:
        user_input = input("Confirm? (y/n): ")
        if user_input != 'y':
            task.mark_skipped()
            continue
    
    # 3. Route to correct executor
    success = _execute_single_task(task)
    
    # 4. Record result
    task.mark_completed() if success else task.mark_failed()

# Output:
📋 TASK PLAN: find duplicate files in Downloads
 ──────────────────────────────
Total Tasks: 2

  💼 TASK 1: Scan for duplicates
     └─ Status: Running...

  💼 TASK 2: Remove duplicates  
     └─ Status: Waiting (depends on Task 1)

============================================================
✅ EXECUTION COMPLETE
============================================================
Total: 2 | ✓ Completed: 2 | ✗ Failed: 0 | ⏭️  Skipped: 0
============================================================
```

---

### **STEP 5: INTENT ROUTER - DISPATCH**
**File:** `router/intent_router.py` → `route_intent()`

```python
# Input: Task(intent, action, parameters)
command = {
    "intent": "file_management",
    "action": "find_duplicates",
    "parameters": {"folder_name": "Downloads"}
}

# Validation
is_valid, msg = _validate_intent(intent, action)
if not is_valid:
    print(f"❌ Validation Error: {msg}")
    return

# Routing (pattern matching)
if intent == "file_management":
    if action == "find_duplicates":
        handler = _handle_file_management(action, params)
        # Routes to: file_manager/duplicate_detector.py
    elif action == "organize_by_type":
        # Routes to: file_manager/file_organizer.py
    elif action == "find_large_files":
        # Routes to: file_manager/large_file_scanner.py
```

**Mapping Table:**

| Intent | Action | Routes To |
|--------|--------|-----------|
| `system_control` | `sleep` | `tools/system_tools.py::sleep_pc()` |
| `system_control` | `lock` | `tools/system_tools.py::lock_pc()` |
| `diagnostics` | `check_cpu` | `tools/diagnostics_tools.py::get_cpu_usage()` |
| `diagnostics` | `check_ram` | `tools/diagnostics_tools.py::get_ram_usage()` |
| `file_management` | `find_duplicates` | `file_manager/duplicate_detector.py` |
| `file_management` | `organize_by_type` | `file_manager/file_organizer.py` |
| `installer` | `install_software` | `installer/installer.py` |
| `vision_analysis` | `analyze_screen` | `troubleshooter/vision_analyzer.py` |
| `personalization` | `toggle_dark_mode` | `personalisation/personalisation_tools.py` |

---

### **STEP 6: TOOLS - ACTUAL EXECUTION**

The router calls the actual tool functions that perform the work:

#### **System Tools** (`tools/system_tools.py`)
```python
def sleep_pc():
    """Put system to sleep"""
    os.system("rundll32.exe powrprof.dll,SetSuspendState 0,1,0")

def lock_pc():
    """Lock the workstation"""
    os.system("rundll32.exe user32.dll,LockWorkStation")

def kill_process(process_name):
    """Kill a running process"""
    os.system(f"taskkill /IM {process_name} /F")
```

#### **Diagnostics Tools** (`tools/diagnostics_tools.py`)
```python
def get_cpu_usage():
    """Get CPU metrics using psutil"""
    cpu_percent = psutil.cpu_percent(interval=1)
    cpu_count = psutil.cpu_count()
    per_cpu = psutil.cpu_percent(interval=1, percpu=True)
    return {
        "cpu_percent": cpu_percent,
        "cpu_count": cpu_count,
        "per_cpu": per_cpu
    }

def get_ram_usage():
    """Get RAM metrics"""
    ram = psutil.virtual_memory()
    return {
        "total_gb": round(ram.total / (1024**3), 2),
        "used_gb": round(ram.used / (1024**3), 2),
        "available_gb": round(ram.available / (1024**3), 2),
        "percent": ram.percent
    }
```

#### **File Manager Tools** (`file_manager/`)
```python
# duplicate_detector.py
def find_duplicates(folder_path):
    """Find duplicate files by hash comparison"""
    # Returns: {duplicates_found, groups, total_size}

# file_organizer.py
def organize_files(folder_path, strategy="by_type"):
    """Organize files into subfolders"""
    # Returns: {files_moved, new_folders_created}

# large_file_scanner.py
def find_large_files(folder_path, min_size_mb=100):
    """Find large files"""
    # Returns: {large_files, total_space}
```

#### **Vision Tools** (`troubleshooter/`, `vision/`)
```python
def capture_screen_base64():
    """Take screenshot and convert to base64"""
    # Returns: base64_encoded_image

def analyze_screen_with_gemini(image_base64):
    """Send screenshot to Google Gemini for analysis"""
    # Returns: {analysis, problems_detected, solutions}

def install_app(app_name):
    """Install app using winget"""
    # Returns: {status, version, message}
```

---

### **STEP 7: OUTPUT DISPLAY**

Results are displayed in multiple formats:

```
📊 RESULT DISPLAY:

🖥️ CPU USAGE
├─ Usage: 45.2%
├─ Cores: 8
└─ Per-Core: [10.5, 45.2, 32.1, 28.0, 15.3, 38.9, 42.1, 22.5]

Memory Usage: 12.5 GB / 32 GB (39%)

💾 Disk Space (C:)
├─ Used: 256 GB
├─ Free: 244 GB
└─ Usage: 51.2%

✅ Tasks Completed Successfully!
```

---

### **STEP 8: MEMORY STORAGE**

All executions are persisted to `agent_memory/`:

```json
// execution_history.json
[
  {
    "timestamp": "2026-03-28T10:30:45.123456",
    "plan_id": "abc12345",
    "user_input": "find duplicate files in Downloads",
    "goal": "find duplicate files in Downloads",
    "summary": {
      "total_tasks": 2,
      "completed": 2,
      "failed": 0,
      "skipped": 0,
      "execution_time_seconds": 45.2
    }
  }
]

// context.json
{
  "last_command": {
    "value": "find duplicate files in Downloads",
    "timestamp": "2026-03-28T10:30:45.123456"
  },
  "user_theme": {
    "value": "dark",
    "timestamp": "2026-03-28T09:15:00.000000"
  }
}

// preferences.json
{
  "auto_confirm": {
    "value": false,
    "timestamp": "2026-03-28T08:00:00.000000"
  }
}
```

**Benefits:**
- Learn from past executions
- Provide context for future commands
- Track user preferences
- Enable statistics and reporting

---

## 💻 Core Components

### **1. MasterAgent** (`agent/master_agent.py`)
**Role:** Central orchestrator - ties everything together

**Main Methods:**
- `process_command(user_input)` - Entry point for all user input
- `_extract_goal(user_input)` - Uses LLM to understand intent
- `_get_execution_report(plan)` - Generates detailed report

**Responsibilities:**
- Receive user input
- Extract goal using LLM
- Create task plans
- Execute plans
- Store results in memory
- Provide feedback to user

**Example Usage:**
```python
agent = MasterAgent()
result = agent.process_command("Check my disk space", auto_confirm=False)
print(result)
```

---

### **2. TaskPlanner** (`agent/task_planner.py`)
**Role:** Breaks down complex goals into executable tasks

**Main Methods:**
- `plan_goal(user_input, goal)` - Create task plan
- `_decompose_goal(goal)` - Use LLM to decompose
- `_validate_task(task)` - Verify task structure

**Responsibilities:**
- Understand decomposition of goals
- Create TaskPlan with ordered tasks
- Assign dependencies
- Mark confirmation requirements
- Set priorities

**Example Task Structure:**
```python
{
    "order": 1,
    "intent": "diagnostics",
    "action": "check_cpu",
    "parameters": {},
    "description": "Check CPU usage",
    "priority": "high",
    "depends_on": [],
    "requires_confirmation": False
}
```

---

### **3. ExecutionEngine** (`agent/execution_engine.py`)
**Role:** Executes task plans with error handling

**Main Methods:**
- `execute_plan(plan, auto_confirm)` - Run all tasks
- `_execute_single_task(task, auto_confirm)` - Run one task
- `get_execution_report(plan)` - Generate report

**Responsibilities:**
- Execute tasks sequentially
- Validate dependencies
- Handle confirmations
- Catch and log errors
- Retry failed tasks
- Generate detailed reports

**Features:**
- Dependency tracking
- Timeout protection
- Retry logic
- Dry-run mode
- Comprehensive logging

---

### **4. IntentRouter** (`router/intent_router.py`)
**Role:** Dispatcher that routes tasks to correct tools

**Main Methods:**
- `route_intent(command)` - Route to handler
- `_validate_intent(intent, action)` - Validate command
- `_handle_[intent_name](action, params)` - Handler for each intent

**Responsibilities:**
- Validate intent and action
- Map to correct handler
- Pass parameters
- Error handling

**Supported Intents:**
```
• help
• open_app
• web_search
• system_control
• diagnostics
• disk_analysis
• maintenance
• health_check
• troubleshoot_screen
• vision_analysis
• personalization
• file_management
• system_config
• installer
```

---

### **5. MemoryManager** (`agent/memory_manager.py`)
**Role:** Persistent storage and context management

**Main Methods:**
- `save_execution(plan)` - Log execution
- `save_plan(plan)` - Save task plan
- `add_to_context_chain(text, role)` - Store context
- `get_execution_history(limit)` - Retrieve history

**Stored Data:**
- `execution_history.json` - All past executions
- `context.json` - Conversation context
- `plans.json` - Task plans
- `preferences.json` - User preferences

---

### **6. VoiceListener** (`voice/voice_listener.py`)
**Role:** Microphone input and speech-to-text

**Main Methods:**
- `listen()` - Record and transcribe
- `_record_until_silence()` - Audio capture
- `_transcribe(audio_data)` - Send to Groq Whisper

**Features:**
- Auto-stop on silence
- Noise detection
- Groq Whisper large-v3 API
- Configurable thresholds

---

### **7. LLM Clients** (`llm/`)

#### **GroqClient** (`llm/groq_client.py`)
- Uses Groq API for fast LLM inference
- For task planning and goal extraction
- Also provides speech-to-text via Whisper

#### **GeminiClient** (`llm/gemini_client.py`)
- Uses Google Gemini for vision tasks
- Screenshot analysis
- Problem detection

---

## 📚 Concrete Examples

### **Example 1: Simple Diagnostic Command**

**User Input:**
```
"Check my CPU usage"
```

**Processing:**

```
Step 1: MasterAgent receives input
  ↓
Step 2: Extract Goal (LLM)
  Input: "Check my CPU usage"
  LLM Response: {
    "intent": "diagnostics",
    "action": "check_cpu",
    "confidence": 0.99
  }
  Display: 🎯 Goal: Check CPU usage
          🔍 Confidence: 99%
  ↓
Step 3: TaskPlanner creates plan
  Output: 1 task
  {
    "intent": "diagnostics",
    "action": "check_cpu",
    "parameters": {},
    "description": "Check current CPU usage"
  }
  Display: 📄 Created plan with 1 task
  ↓
Step 4: ExecutionEngine executes
  Task 1: Check CPU
    └─ Running...
  ↓
Step 5: IntentRouter routes to diagnostics handler
  intent=diagnostics, action=check_cpu
  └─ Call: tools/diagnostics_tools.py::get_cpu_usage()
  ↓
Step 6: Tool executes
  psutil.cpu_percent(interval=1) → 45.2%
  psutil.cpu_count() → 8
  psutil.cpu_percent(percpu=True) → [10.5, 45.2, ...]
  
  Return: {
    "cpu_percent": 45.2,
    "cpu_count": 8,
    "per_cpu": [10.5, 45.2, 32.1, 28.0, 15.3, 38.9, 42.1, 22.5]
  }
  ↓
Step 7: Display output
  🖥️ CPU USAGE
  ├─ Usage: 45.2%
  ├─ Cores: 8
  └─ Per-Core: [10.5, 45.2, 32.1, 28.0, 15.3, 38.9, 42.1, 22.5]
  ↓
Step 8: Store in memory
  execution_history.json updated
  ✅ EXECUTION COMPLETE
```

---

### **Example 2: Complex Multi-Step Task**

**User Input:**
```
"Organize my Downloads, find duplicate files, and tell me how much space they use"
```

**Processing:**

```
Step 1-2: MasterAgent extracts goal
  LLM Response: {
    "intent": "mixed",
    "goal": "organize downloads and manage duplicates",
    "is_complex": true
  }
  Display: 🎯 Goal: organize downloads and manage duplicates
          🔍 Confidence: 92%
  ↓
Step 3: TaskPlanner decomposes into 3 tasks
  
  Task 1: Organize Files
  {
    "order": 1,
    "intent": "file_management",
    "action": "organize_by_type",
    "parameters": {"folder_name": "Downloads"},
    "priority": "high"
  }
  
  Task 2: Find Duplicates
  {
    "order": 2,
    "intent": "file_management",
    "action": "find_duplicates",
    "parameters": {"folder_name": "Downloads"},
    "depends_on": [],
    "priority": "high"
  }
  
  Task 3: Analyze Disk Usage
  {
    "order": 3,
    "intent": "disk_analysis",
    "action": "analyze_usage",
    "parameters": {"folder_name": "Downloads"},
    "depends_on": [],
    "priority": "medium"
  }
  
  Display: 📄 Created plan with 3 tasks
  ↓
Step 4: ExecutionEngine executes sequentially
  
  ════════════════════════════════════════════
  📋 TASK PLAN: organize downloads and manage duplicates
  ════════════════════════════════════════════
  Total Tasks: 3
  
  ─── TASK 1/3: Organize files by type ──────
  Status: RUNNING...
    └─ Scanning Downloads folder...
    └─ Found: 342 files
    └─ Creating folders: Documents, Images, Videos, Archives, Code
    └─ Moving files...
  Status: ✓ COMPLETED
    └─ Files organized: 342
    └─ New folders created: 5
    └─ Time taken: 8.2 seconds
  
  ─── TASK 2/3: Find duplicate files ────────
  Status: RUNNING...
    └─ Scanning Downloads folder...
    └─ Computing file hashes...
    └─ Comparing hashes...
  Status: ✓ COMPLETED
    └─ Duplicates found: 23
    └─ Total duplicate size: 845 MB
    └─ Time taken: 15.3 seconds
  
  ─── TASK 3/3: Analyze disk usage ─────────
  Status: RUNNING...
    └─ Calculating folder size...
  Status: ✓ COMPLETED
    └─ Used: 12.5 GB
    └─ Total: 100 GB
    └─ Usage: 12.5%
    └─ Time taken: 3.1 seconds
  
  ════════════════════════════════════════════
  ✅ EXECUTION COMPLETE
  ════════════════════════════════════════════
  Total: 3 | ✓ Completed: 3 | ✗ Failed: 0 | ⏭️  Skipped: 0
  ════════════════════════════════════════════
  ↓
Step 5-8: Display & Storage
  📊 SUMMARY REPORT
  ──────────────────────────────────────────
  ✅ All tasks completed successfully
  
  Task Results:
    1. Organized 342 files into 5 categories
    2. Found 23 duplicate files (845 MB)
    3. Downloads folder: 12.5 GB of 100 GB (12.5%)
  
  Execution Time: 26.6 seconds
  Status: SUCCESS
  ──────────────────────────────────────────
  
  execution_history.json updated with all results
```

---

### **Example 3: Voice Input with Confirmation**

**User Speaks:**
```
"Install Discord"
```

**Processing:**

```
Step 1: VoiceListener captures audio
  🎤 Listening... (speak now, auto-stops on silence)
  [User speaks: "install discord"]
  📝 You said: "install discord"
  ↓
Step 2: MasterAgent receives text
  Input: "install discord"
  ↓
Step 3: Extract Goal (LLM)
  LLM Response: {
    "intent": "installer",
    "action": "install_software",
    "parameters": {"app_name": "Discord"},
    "requires_safety_check": true
  }
  Display: 🎯 Goal: Install Discord
          🔍 Confidence: 98%
          ⚠️ This operation involves system changes
  ↓
Step 4: TaskPlanner creates plan
  {
    "intent": "installer",
    "action": "install_software",
    "parameters": {"app_name": "Discord"},
    "requires_confirmation": true
  }
  Display: 📄 Created plan with 1 task
  ↓
Step 5: ExecutionEngine executes
  Before executing, asks for confirmation!
  
  ──────────────────────────────────────────
  📦 TASK: Install Discord
  ──────────────────────────────────────────
  
  ⚠️ This task requires confirmation
  
  About to execute:
    • Download Discord installer
    • Execute installation
    • Verify installation
  
  Continue? (y/n): y
  ↓
Step 6: IntentRouter routes
  intent=installer, action=install_software
  └─ Call: installer/installer.py::install_app("Discord")
  ↓
Step 7: Tool executes
  Using Windows Package Manager (winget):
  $ winget install Discord
  
  [==================================================] 100%
  
  Discord 0.0.301 installed successfully
  
  Return: {
    "status": "installed",
    "version": "0.0.301",
    "path": "C:\\Users\\..\\Discord"
  }
  ↓
Step 8: Display output
  ✅ Successfully installed Discord v0.0.301
  
  Installation Details:
    • Package: Discord
    • Version: 0.0.301
    • Location: C:\Users\ASUS\AppData\Local\Discord
    • Time taken: 45.2 seconds
  
  ✅ EXECUTION COMPLETE
  ↓
Step 9: Store in memory
  execution_history.json updated
```

---

### **Example 4: Troubleshooting Workflow**

**User Input:**
```
"My screen looks broken, can you fix it?"
```

**Processing:**

```
Step 1-3: MasterAgent extracts & plans
  LLM Response: {
    "intent": "troubleshoot_screen",
    "goal": "analyze screen issues and provide solutions"
  }
  
  TaskPlan: 4 tasks
    1. Capture screenshot
    2. Analyze with vision AI
    3. Identify problems
    4. Suggest solutions (+ auto-fix if user approves)
  ↓
Step 4: ExecutionEngine executes
  
  TASK 1: Capture Screenshot
    Status: RUNNING...
      └─ Using mss to capture screen
      └─ Converting to base64
    Status: ✓ COMPLETED (2.1 seconds)
  
  TASK 2: Vision Analysis
    Status: RUNNING...
      └─ Sending screenshot to Google Gemini
      └─ Analyzing image...
    Status: ✓ COMPLETED
    Result: {
      "issues": [
        "Color banding on right side",
        "Flickering at bottom",
        "Oversaturation in blues"
      ]
    }
  
  TASK 3: Identify Root Causes
    Status: ✓ COMPLETED
    LLM Analysis: {
      "root_causes": [
        {
          "issue": "Color banding",
          "cause": "Outdated GPU driver",
          "severity": "high"
        },
        {
          "issue": "Flickering",
          "cause": "Refresh rate mismatch",
          "severity": "medium"
        }
      ]
    }
  
  TASK 4: Generate Solutions
    Status: ✓ COMPLETED
    Solutions:
      1. Update GPU drivers (Recommended - High Impact)
      2. Change refresh rate to 60Hz (Medium Impact)
      3. Disable V-Sync (Low Impact)
  ↓
Step 5-7: Display & get user approval
  
  🔍 SCREEN ANALYSIS RESULTS
  ════════════════════════════════════════════
  
  Issues Detected:
    ⚠️ Color banding on right side
    ⚠️ Flickering at bottom
    ⚠️ Oversaturation in blues
  
  Root Causes:
    🔧 Outdated GPU driver (HIGH PRIORITY)
    🔧 Refresh rate mismatch (MEDIUM)
    🔧 V-Sync settings (LOW)
  
  Solutions (in priority order):
    1️⃣ Update GPU Drivers ⭐ RECOMMENDED
       • Download latest driver from NVIDIA/AMD
       • Install & restart system
       • Expected result: Fix 80% of issues
    
    2️⃣ Change Refresh Rate to 60Hz
       • Settings → Display → Refresh Rate
       • Expected result: Fix flickering
    
    3️⃣ Disable V-Sync
       • Graphics settings → Advanced
       • Expected result: Minor improvement
  
  ════════════════════════════════════════════
  
  Auto-fix recommended solutions? (y/n): y
  ↓
Step 8: Apply fixes
  
  Executing auto-fixes...
  
  FIX 1: Update GPU Driver
    └─ Detecting GPU: NVIDIA GeForce RTX 3070
    └─ Current driver: 551.44
    └─ Latest driver: 563.41
    └─ Downloading...
    └─ Installing...
    └─ ✓ Driver updated successfully!
    └─ ⚠️ Restart required
  
  FIX 2: Change Refresh Rate
    └─ Current: 50 Hz
    └─ Setting to: 60 Hz
    └─ ✓ Changes applied
  
  FIX 3: V-Sync Disabled
    └─ ✓ V-Sync disabled in GPU control panel
  
  ════════════════════════════════════════════
  ✅ AUTO-FIXES COMPLETED
  ════════════════════════════════════════════
  
  Summary:
    ✓ GPU driver updated
    ✓ Refresh rate changed to 60Hz
    ✓ V-Sync disabled
    
  Please restart your PC for changes to take effect.
  ↓
Step 9: Store results
  execution_history.json updated
```

---

## �️ Complete Tools & Layers Reference

This section provides exhaustive details on all tool modules and their functions. No detail is left out.

---

### **LAYER 1: SYSTEM CONTROL TOOLS** (`tools/system_tools.py`)

**Functions:**
| Function | Parameters | What It Does | Example |
|----------|-----------|--------------|---------|
| `sleep_pc()` | None | Puts PC to sleep mode | `sleep_pc()` |
| `lock_pc()` | None | Locks the workstation (shows login) | `lock_pc()` |
| `kill_process()` | `process_name: str` | Terminates a running process | `kill_process("chrome.exe")` |
| `clean_temp()` | None | Deletes temporary files from system TEMP folder | `clean_temp()` |
| `empty_recycle_bin()` | None | Empties the Windows Recycle Bin permanently | `empty_recycle_bin()` |

**Example Usage Flow:**
```
User: "Put my PC to sleep"
    ↓
TaskPlanner: {intent: "system_control", action: "sleep"}
    ↓
ExecutionEngine: Requires confirmation
    ↓
IntentRouter: Routes to system_tools.sleep_pc()
    ↓
Result: PC enters sleep mode
```

---

### **LAYER 2: DIAGNOSTICS TOOLS** (`tools/diagnostics_tools.py`)

**Functions:**

#### **CPU Diagnostics**
```python
def get_cpu_usage() -> dict
```
Returns real-time CPU metrics using psutil.

**Output:**
```json
{
  "cpu_percent": 45.2,
  "cpu_count": 8,
  "per_cpu": [10.5, 45.2, 32.1, 28.0, 15.3, 38.9, 42.1, 22.5]
}
```

#### **RAM Diagnostics**
```python
def get_ram_usage() -> dict
```
Returns memory usage statistics.

**Output:**
```json
{
  "total_gb": 32.0,
  "used_gb": 12.5,
  "available_gb": 19.5,
  "percent": 39.1
}
```

#### **Disk Diagnostics**
```python
def get_disk_usage(drive: str = "C:") -> dict
```
Returns storage information for a specific drive.

**Output:**
```json
{
  "total_gb": 500.0,
  "used_gb": 256.0,
  "free_gb": 244.0,
  "percent": 51.2
}
```

#### **Complete System Health**
```python
def get_complete_system_health() -> dict
```
Comprehensive health check combining CPU, RAM, Disk, and alerts.

**Output:**
```json
{
  "cpu": {...},
  "ram": {...},
  "disk": {...},
  "alerts": ["Disk usage above 70%", "Temperature normal"],
  "overall_status": "good"
}
```

#### **Disk Analysis**
```python
def analyze_disk_usage(drive: str = "C:") -> dict
```
Detailed disk usage analysis with breakdown by folder type.

#### **Large Files Finding**
```python
def find_large_folders(drive: str = "C:", limit: int = 10) -> dict
```
Finds the largest folders consuming disk space.

#### **File Scanning**
```python
def scan_cleanup_files(folder: str = "C:\\") -> dict
```
Scans for old and temporary files suitable for cleanup.

---

### **LAYER 3: APPLICATION TOOLS** (`tools/app_tools.py`)

**Functions:**

```python
def open_app(app_name: str)
```

**Supported Apps:**
```python
APP_PATHS = {
    "chrome": r"C:\Program Files\Google\Chrome\Application\chrome.exe",
    "notepad": r"C:\Windows\System32\notepad.exe",
    "calculator": r"C:\Windows\System32\calc.exe",
    "cmd": r"C:\Windows\System32\cmd.exe",
    ... (more can be added)
}
```

**Usage:**
```
User: "Open Chrome"
    ↓
route_intent(intent="open_app", parameters={"app_name": "Chrome"})
    ↓
open_app("chrome")
    ↓
Chrome browser launches
```

---

### **LAYER 4: WEB TOOLS** (`tools/web_tools.py`)

**Functions:**

```python
def search_web(query: str)
```

Opens Google search in default browser for the given query.

**Usage:**
```
User: "Search for Python tutorial"
    ↓
search_web("Python tutorial")
    ↓
Browser opens: https://www.google.com/search?q=Python+tutorial
```

---

### **LAYER 5: SYSTEM CONFIG TOOLS** (`tools/system_config.py`)

**Functions:**
| Function | Purpose | Output |
|----------|---------|--------|
| `get_system_config()` | Get complete system info | `{os, processor, ram, disk, hostname, ...}` |
| `print_system_config()` | Display system info in formatted table | Prints to console |

**Typical Output:**
```
Operating System: Windows 11 Pro
Processor: Intel Core i7-10700K
RAM: 32 GB
Disk Space: 500 GB (256 GB used)
Hostname: DESKTOP-ABC123
Username: ASUS
Python Version: 3.11.0
...
```

---

### **LAYER 6: FILE MANAGER TOOLS** (`file_manager/`)

The file manager is modular with multiple sub-components:

#### **6a. Manager Core** (`file_manager/manager.py`)

**Main Class:** `FileManager`

**Public Methods:**
```python
class FileManager:
    def load_folder(self, folder_name: str) -> Dict
    def analyze_folder(self) -> Dict
    def scan_large_files(self, min_size_mb: float = 100, limit: int = 20) -> Dict
    def scan_duplicates(self, by_name: bool = False) -> Dict
    def remove_duplicates(self, dry_run: bool = True) -> Dict
    def organize_files(self, strategy: str = "by_type", dry_run: bool = False) -> Dict
    def get_comprehensive_report(self) -> Dict
```

#### **6b. File Organizer** (`file_manager/file_organizer.py`)

**Main Function:**
```python
def organize_files_by_type(folder_path: Path, dry_run: bool = False) -> Dict
```

**File Categories Created:**
| Category | Extensions |
|----------|-----------|
| **Images** | .jpg, .jpeg, .png, .gif, .bmp, .svg, .webp, .ico |
| **Documents** | .pdf, .doc, .docx, .txt, .xlsx, .xls, .ppt, .pptx |
| **Videos** | .mp4, .avi, .mkv, .mov, .wmv, .flv, .webm |
| **Audio** | .mp3, .wav, .flac, .aac, .m4a, .ogg, .wma |
| **Archives** | .zip, .rar, .7z, .tar, .gz, .bz2 |
| **Code** | .py, .js, .java, .cpp, .c, .html, .css, .php, .go, .rs |
| **Data** | .json, .xml, .csv, .sql, .db |
| **Executables** | .exe, .bat, .sh, .msi, .app |
| **Others** | Everything else |

**Output Example:**
```json
{
  "total_files_processed": 342,
  "files_moved": 338,
  "organization_map": {
    "Images": 125,
    "Documents": 87,
    "Videos": 65,
    "Audio": 32,
    "Archives": 18,
    "Code": 10,
    "Others": 5
  },
  "errors": []
}
```

#### **6c. Duplicate Detector** (`file_manager/duplicate_detector.py`)

**Main Function:**
```python
def find_duplicates(folder_path: Path, by_name: bool = False) -> Dict
```

**Detection Methods:**
- **By Hash (default):** Compares SHA256 hash of files (most accurate)
- **By Name:** Groups files with identical names

**Output Example:**
```json
{
  "total_files": 542,
  "duplicates_found": true,
  "duplicate_groups": [
    {
      "hash_or_name": "abc123...",
      "files": [
        "/path/to/photo1.jpg",
        "/path/to/photo1_copy.jpg",
        "/path/to/photo1_backup.jpg"
      ],
      "count": 3,
      "size_per_file": 2500000,
      "total_duplicate_size": 5000000
    }
  ],
  "total_duplicate_files": 23,
  "space_savings": 845000000
}
```

**Duplicate Removal:**
```python
def remove_duplicates(folder_path: Path, dry_run: bool = True) -> Dict
```
Identifies largest duplicate set and removes, keeping one master copy.

#### **6d. Large File Scanner** (`file_manager/large_file_scanner.py`)

**Main Function:**
```python
def find_large_files(folder_path: Path, min_size_mb: float = 100, limit: int = 20) -> Dict
```

**Output Example:**
```json
{
  "total_files": 10234,
  "large_files_found": 45,
  "large_files": [
    {
      "path": "C:/Users/.../Downloads/movie.mp4",
      "name": "movie.mp4",
      "size_bytes": 2147483648,
      "size_mb": 2048,
      "size_gb": 2.0
    },
    {
      "path": "C:/Users/.../Documents/backup.iso",
      "name": "backup.iso",
      "size_bytes": 1073741824,
      "size_mb": 1024,
      "size_gb": 1.0
    }
  ],
  "total_scanned_size_mb": 125000.00
}
```

#### **6e. Folder Finder** (`file_manager/folder_finder.py`)

**Main Function:**
```python
def find_folder(folder_name: str) -> Optional[Path]
```

**Search Locations (in order):**
1. Current working directory
2. Parent of current directory
3. Desktop
4. Downloads
5. Documents
6. User home directory

**Usage:**
```python
path = find_folder("Downloads")  # Returns: Path object or None
```

#### **6f. Folder Information** (`file_manager/folder_finder.py`)

**Function:**
```python
def get_folder_info(folder_path: Path) -> dict
```

**Output:**
```json
{
  "folder_name": "Downloads",
  "folder_path": "C:/Users/ASUS/Downloads",
  "file_count": 342,
  "subfolder_count": 12,
  "total_size_mb": 12500,
  "total_size_gb": 12.5,
  "file_size_breakdown": {
    "< 1MB": 150,
    "1-10MB": 100,
    "10-100MB": 65,
    "> 100MB": 27
  }
}
```

---

### **LAYER 7: PERSONALIZATION TOOLS** (`personalisation/personalisation_tools.py`)

**Functions:**

#### **Dark Mode / Light Mode**
```python
def toggle_dark_mode(enable: bool = True) -> dict
```
**Output:**
```json
{
  "success": true,
  "message": "✅ Switched to 🌙 DARK MODE",
  "mode": "dark"
}
```

```python
def get_current_theme() -> dict
```
**Output:**
```json
{
  "success": true,
  "current_theme": "dark"
}
```

#### **Wallpaper Management**
```python
def set_wallpaper(image_path: str) -> dict
```
**Usage:**
```python
result = set_wallpaper("C:/Users/ASUS/Pictures/sunset.jpg")
```
**Output:**
```json
{
  "success": true,
  "message": "✅ Wallpaper changed to: sunset.jpg",
  "wallpaper_path": "C:/Users/ASUS/Pictures/sunset.jpg"
}
```

#### **Display Brightness**
```python
def set_display_brightness(brightness: int) -> dict
```
**Parameters:** `brightness` (0-100)

**Usage:**
```python
set_display_brightness(75)  # Set to 75%
```
**Output:**
```json
{
  "success": true,
  "brightness": 75,
  "message": "✅ Brightness set to 75%"
}
```

#### **Theme/Color Management**
```python
def set_accent_color(color: str) -> dict
def apply_theme_preset(theme_name: str) -> dict
```

#### **Profile Management**
```python
def save_personalization_profile(profile_name: str) -> dict
def load_personalization_profile(profile_name: str) -> dict
```

#### **System Settings**
```python
def manage_startup_apps(action: str, app_name: str) -> dict
def set_default_app(app_type: str, app_name: str) -> dict
```

---

### **LAYER 8: TROUBLESHOOTER TOOLS** (`troubleshooter/`)

#### **8a. Screenshot Tool** (`troubleshooter/screenshot_tool.py`)

**Function:**
```python
def capture_screen_base64() -> str
```

**Output:** Base64-encoded PNG image string (for sending to vision API)

**Usage:**
```python
screenshot_b64 = capture_screen_base64()
# Send to Gemini Vision API
```

#### **8b. Vision Analyzer** (`troubleshooter/vision_analyzer.py`)

**Class:** `VisionAnalyzer`

**Main Method:**
```python
def analyze_screenshot_for_issues(screenshot_path: str) -> dict
```

**Output Example:**
```json
{
  "issues_detected": [
    {
      "type": "color_banding",
      "severity": "high",
      "description": "Color banding visible on right side of screen",
      "region": "right_half"
    },
    {
      "type": "flickering",
      "severity": "medium",
      "description": "Screen flickering detected at bottom",
      "region": "bottom_quarter"
    }
  ],
  "recommendations": [
    "Update GPU drivers",
    "Check refresh rate settings",
    "Test monitor cable"
  ],
  "root_causes": [
    "Outdated GPU driver",
    "Refresh rate mismatch",
    "Cable interference"
  ]
}
```

#### **8c. Solution Parser** (`troubleshooter/solution_parser.py`)

**Function:**
```python
def parse_solution(raw_response: str) -> dict
```

Converts LLM response to structured JSON with solutions.

**Output:**
```json
{
  "solutions": [
    {
      "title": "Update GPU Drivers",
      "steps": [
        "Go to Device Manager",
        "Find Display adapters",
        "Right-click GPU → Update driver",
        "Select 'Search automatically'"
      ],
      "priority": "high",
      "estimated_time": "10 minutes"
    }
  ]
}
```

#### **8d. Auto-Fix Engine** (`troubleshooter/auto_fix_engine.py`)

**Function:**
```python
def execute_fixes(commands: list)
```

Executes PowerShell commands to automatically apply fixes.

**Command Format:**
```python
commands = [
    {
        "command": "Update-WmiObject ...",
        "requires_admin": True,
        "description": "Update GPU driver"
    },
    {
        "command": "Set-ItemProperty ...",
        "requires_admin": True,
        "description": "Change refresh rate"
    }
]
```

**Admin Detection:**
```python
def is_admin() -> bool
```
Returns whether current user has admin privileges.

---

### **LAYER 9: VISION TOOLS** (`vision/vision_engine.py`)

**Functions:**

#### **Screen Capture**
```python
def capture_screen() -> tuple[Path, str]
```
**Returns:** `(path_to_screenshot.png, timestamp_iso)`

#### **Screen Analysis with Gemini**
```python
def analyze_screen(
    image_path: Path,
    user_query: str = "What is on this screen?",
    api_key: str = None,
    model: str = "gemini-2.5-flash"
) -> dict
```

**Output Example:**
```json
{
  "timestamp_iso": "2026-03-28T10:30:45Z",
  "screen_summary": "Windows desktop with Chrome browser open showing Google homepage",
  "window_title_or_app": "Google - Google Chrome",
  "detected_buttons": [
    {
      "label": "Search",
      "context": "Main search bar",
      "action_hint": "Perform search"
    },
    {
      "label": "Sign in",
      "context": "Top right corner",
      "action_hint": "Open Google account login"
    }
  ],
  "detected_options": [
    {
      "label": "Settings",
      "type": "menu",
      "context": "Top navigation"
    }
  ],
  "visible_text_regions": [
    {
      "text": "Google Search",
      "role": "heading"
    }
  ],
  "possible_errors_or_alerts": [],
  "accessibility_hints": "Standard web page layout, no modals or dialogs detected"
}
```

**Output Saving:**
```
vision_output/
├── analysis_20260228T100030Z.json     # Gemini analysis
├── ocr_analysis_20260228T100030Z.json # OCR results
└── screenshots/
    └── screen_20260228_100030.png     # Actual screenshot
```

---

### **LAYER 10: INSTALLER TOOLS** (`installer/`)

The installer system is comprehensive with multiple components:

#### **10a. Main Installer Agent** (`installer/installer.py`)

**Class:** `InstallerAgent`

**Methods:**

```python
def install_wallpaper(query: str, folder: str = None, auto_set: bool = False) -> Dict
```
**Example:**
```python
result = installer.install_wallpaper("sunset landscape", auto_set=True)
```
**Output:**
```json
{
  "status": "success",
  "wallpaper_path": "/path/to/sunset.jpg",
  "from_cache": false,
  "auto_set": true,
  "message": "Downloaded sunset.jpg and set as background"
}
```

```python
def install_software(app_name: str, version: str = "latest") -> Dict
```
**Example:**
```python
result = installer.install_software("Discord")
```
**Uses:** Windows Package Manager (winget)

**Output:**
```json
{
  "status": "success",
  "app_name": "Discord",
  "version": "0.0.301",
  "install_path": "C:\\Users\\ASUS\\AppData\\Local\\Discord",
  "message": "Successfully installed Discord v0.0.301"
}
```

#### **10b. Resource Finder** (`installer/resource_finder.py`)

**Class:** `ResourceFinder`

**Methods:**

```python
def find_resource(query: str, resource_type: str) -> Optional[Dict]
```

**Supported Resource Types:**
- `"wallpaper"` - Desktop wallpapers
- `"software"` - Software packages
- `"image"` - General images
- `"document"` - Documents

**Sources (Priority Order):**
1. Unsplash API (highest quality images)
2. Pexels API
3. Pixabay
4. Google Images

**Output Example:**
```json
{
  "name": "Beautiful Sunset",
  "url": "https://images.unsplash.com/...",
  "size": "2.5MB",
  "quality": "high",
  "format": "jpg",
  "source": "unsplash",
  "filename": "beautiful_sunset.jpg"
}
```

#### **10c. Downloader** (`installer/downloader.py`)

**Class:** `Downloader`

**Methods:**

```python
def download_file(
    url: str,
    filename: str = None,
    folder: str = None,
    on_progress: Callable = None,
    show_progress: bool = True
) -> Dict
```

**Features:**
- Progress tracking with callbacks
- Resume capability for interrupted downloads
- Automatic retry (3 attempts)
- Timeout protection (30 seconds)
- File verification

**Output:**
```json
{
  "status": "success",
  "path": "/Downloads/file.jpg",
  "size": 2500000,
  "message": "Downloaded 2.5MB successfully"
}
```

#### **10d. Cache Manager** (`installer/cache_manager.py`)

**Class:** `CacheManager`

**Methods:**

```python
def get_cached_resource(resource_id: str) -> Optional[Dict]
```
**Returns:** Cached resource if available and not expired (30 days default)

**Output:**
```json
{
  "path": "/cache/wallpaper_sunset.jpg",
  "cached_at": "2026-02-28T14:30:00",
  "size": 2500000,
  "metadata": {
    "query": "sunset landscape",
    "source": "unsplash"
  }
}
```

```python
def cache_resource(resource_id: str, filepath: str, metadata: Dict) -> bool
```
Stores resource in cache with metadata.

```python
def clear_cache() -> bool
def get_cache_info() -> Dict
```

**Cache Storage:**
```
~/.jarvis_cache/installer/
├── cache_metadata.json       # Metadata index
└── [cached files]
```

---

### **LAYER 11: UTILITIES** (`utils/`)

#### **JSON Parser** (`utils/json_parser.py`)

**Function:**
```python
def extract_json(response: str) -> dict
```
Extracts JSON from LLM responses (handles markdown code fences, etc.)

#### **Prompt Loader** (`utils/prompt_loader.py`)

**Function:**
```python
def load_prompt(filepath: str) -> str
```
Loads prompt templates from `prompts/` directory.

**Stored Prompts:**
- `prompts/system_prompt.txt` - System behavior guide
- `prompts/command_prompt.txt` - Command template

#### **System Monitor** (`utils/system_monitor.py`)

**Functions:**
- Real-time performance monitoring
- Alert generation for threshold violations
- Resource tracking

---

### **LAYER 12: HELPER TOOLS** (`tools/help_commands.py`)

**Functions:**
```python
def print_help()
def show_command_examples()
def get_intent_description(intent: str) -> str
```

Provides user-facing help documentation.

---

## �📁 File Structure Reference

## 📁 Complete File Structure Reference

### **Entry Points & Launchers**
```
main.py                          # CLI entry point (text mode)
run_jarvis.py                    # GUI launcher (selector: GUI/CLI/Voice)
QUICK_REFERENCE.py               # Quick command reference
todo.md                          # Project tasks
README.md                        # Project overview
requirements.txt                 # Python dependencies
```

### **Agent Core Architecture** (`agent/`)
```
agent/
 ├─ master_agent.py             # MAIN: Central orchestrator (58-258 lines)
 │   ├── process_command()       # Entry point for all input
 │   ├── _extract_goal()         # LLM: Understand intent
 │   ├── get_execution_history() # Retrieve past executions
 │   ├── get_execution_stats()   # Statistics
 │   └── get_context_summary()   # Current context
 │
 ├─ task_planner.py             # PLANNING: Decomposes goals (268 lines)
 │   ├── plan_goal()             # Create TaskPlan
 │   ├── _decompose_goal()       # LLM: Break into tasks
 │   ├── _create_fallback_task() # Fallback if decomposition fails
 │   └── _validate_task()        # Verify task structure
 │
 ├─ execution_engine.py          # EXECUTION: Runs tasks (298 lines)
 │   ├── execute_plan()          # Execute all tasks
 │   ├── _execute_single_task()  # Run one task with error handling
 │   ├── _display_task_list()    # Show tasks to user
 │   ├── get_execution_report()  # Generate final report
 │   └── TaskExecutionPhase      # Tracking: VALIDATION, CONFIRMATION, EXECUTION, RESULT_PROCESSING
 │
 ├─ memory_manager.py            # PERSISTENCE: Storage (231 lines)
 │   ├── save_execution()        # Log execution to history
 │   ├── save_plan()             # Save task plan
 │   ├── add_to_context_chain()  # Store conversation context
 │   ├── get_execution_history() # Retrieve past executions
 │   ├── get_execution_stats()   # Get statistics
 │   └── get_context()           # Retrieval of stored context
 │
 ├─ task.py                      # DATA STRUCTURES
 │   ├── Task                    # Single task representation
 │   ├── TaskPlan                # Collection of tasks with dependencies
 │   ├── TaskStatus              # PENDING, RUNNING, COMPLETED, FAILED, SKIPPED
 │   └── TaskPriority            # HIGH, MEDIUM, LOW
 │
 └─ __init__.py                  # Module initialization
```

### **UI Layer** (`ui/`)
```
ui/
 ├─ app.py                       # MAIN APP: GUI application (100 lines)
 │   ├── JarvisGUIApp            # Main GUI class
 │   ├── __init__()              # Initialize UI components
 │   ├── show_chat_card()        # Display chat interface
 │   └── run()                   # Launch GUI
 │
 ├─ chat_card.py                 # CHAT INTERFACE: Message display
 │   ├── ChatCard                # Qt Dialog for chat
 │   ├── add_message()           # Display message
 │   ├── process_input()         # Handle user input
 │   └── update_voice_button()   # Voice input UI
 │
 ├─ __init__.py                  # Module init
 └─ README.md                    # UI documentation
```

### **Voice Input** (`voice/`)
```
voice/
 ├─ voice_listener.py            # SPEECH-TO-TEXT (223 lines)
 │   ├── VoiceListener           # Main voice class
 │   ├── listen()                # Record and transcribe
 │   ├── _record_until_silence() # Audio capture
 │   ├── _transcribe()           # Send to Groq Whisper
 │   └── calibrate_silence_threshold()  # Auto-calibration
 │
 ├─ __init__.py                  # Module init
 └─ __pycache__/
```

### **Routing Layer** (`router/`)
```
router/
 ├─ intent_router.py             # DISPATCHER (1132 lines) ⭐ LARGEST FILE
 │   ├── VALID_INTENTS           # All 14+ intents defined
 │   ├── route_intent()          # Main router
 │   ├── _validate_intent()      # Validate intent/action
 │   ├── _safe_execute()         # Error wrapper
 │   │
 │   ├── _handle_system_control() # System operations
 │   ├── _handle_diagnostics()   # CPU, RAM, Disk checks
 │   ├── _handle_disk_analysis() # Disk usage analysis
 │   ├── _handle_maintenance()   # Cleanup operations
 │   ├── _handle_health_check()  # System health
 │   ├── _handle_troubleshoot_screen()  # Vision troubleshooting
 │   ├── _handle_vision_analysis()      # Screen analysis
 │   ├── _handle_personalization()      # Theme, wallpaper, etc.
 │   ├── _handle_file_management()      # File operations
 │   ├── _handle_file_organization()    # Organize files
 │   ├── _handle_system_config()        # System info
 │   └── _handle_installer()            # Install apps/resources
 │
 └─ __pycache__/
```

### **LLM Clients** (`llm/`)
```
llm/
 ├─ groq_client.py               # GROQ API (Fast inference + Whisper)
 │   ├── GroqClient              # Main class
 │   ├── generate()              # LLM inference
 │   └── transcribe()            # Audio transcription
 │
 ├─ gemini_client.py             # GOOGLE GEMINI API (Vision)
 │   ├── GeminiClient            # Main class
 │   ├── analyze_image()         # Image analysis
 │   └── extract_text()          # OCR
 │
 ├─ __init__.py                  # Module init
 └─ __pycache__/
```

### **System Tools** (`tools/`)
```
tools/
 ├─ system_tools.py              # POWER & PROCESS CONTROL
 │   ├── sleep_pc()              # Sleep/suspend
 │   ├── lock_pc()               # Lock workstation
 │   ├── kill_process()          # Terminate process
 │   ├── clean_temp()            # Delete temp files
 │   └── empty_recycle_bin()     # Empty trash
 │
 ├─ app_tools.py                 # APPLICATION LAUNCHING
 │   ├── open_app()              # Launch application
 │   └── APP_PATHS               # Known app locations
 │
 ├─ web_tools.py                 # WEB OPERATIONS
 │   └── search_web()            # Google search
 │
 ├─ diagnostics_tools.py          # SYSTEM DIAGNOSTICS (810 lines)
 │   ├── get_cpu_usage()          # CPU metrics
 │   ├── get_ram_usage()          # Memory metrics
 │   ├── get_disk_usage()         # Storage metrics
 │   ├── get_complete_system_health()  # Full health check
 │   ├── analyze_disk_usage()     # Detailed disk analysis
 │   ├── find_large_folders()     # Large file finder
 │   ├── scan_cleanup_files()     # Old/temp file scanner
 │   ├── check_health_alerts()    # Alert system
 │   └── execute_cleanup()        # Cleanup executor
 │
 ├─ system_config.py             # SYSTEM INFORMATION
 │   ├── get_system_config()      # Full system info
 │   └── print_system_config()    # Formatted display
 │
 ├─ help_commands.py             # HELP SYSTEM
 │   ├── print_help()            # User-facing help
 │   └── show_command_examples() # Command examples
 │
 ├─ __init__.py                  # Module init
 └─ __pycache__/
```

### **File Manager** (`file_manager/`)
```
file_manager/
 ├─ manager.py                   # MAIN MANAGER (208 lines)
 │   ├── FileManager             # Main class
 │   ├── load_folder()           # Load folder
 │   ├── analyze_folder()        # Comprehensive analysis
 │   ├── scan_large_files()      # Large file detection
 │   ├── scan_duplicates()       # Duplicate detection
 │   ├── remove_duplicates()     # Duplicate removal
 │   ├── organize_files()        # File organization
 │   └── get_comprehensive_report()  # Full report
 │
 ├─ file_organizer.py            # FILE ORGANIZATION (143 lines)
 │   ├── FILE_CATEGORIES         # 9 file categories
 │   ├── get_file_category()     # Classify file
 │   ├── organize_files_by_type()  # Main organizer
 │   └── get_organization_stats() # Statistics
 │
 ├─ duplicate_detector.py         # DUPLICATE FINDER (165 lines)
 │   ├── calculate_file_hash()   # SHA256 hash
 │   ├── find_duplicates()       # Detect duplicates
 │   ├── remove_duplicates()     # Delete duplicates
 │   └── get_space_savings()     # Storage saved
 │
 ├─ large_file_scanner.py         # LARGE FILE FINDER (163 lines)
 │   ├── find_large_files()      # Find files > N MB
 │   ├── get_folder_size_breakdown()  # Size distribution
 │   ├── FileInfo               # File data class
 │   └── format_file_size()     # Human-readable size
 │
 ├─ folder_finder.py              # FOLDER LOCATOR (101 lines)
 │   ├── find_folder()           # Search for folder
 │   ├── get_folder_info()       # Folder metadata
 │   └── search_locations        # Desktop, Downloads, Docs, etc.
 │
 ├─ llm_organizer.py             # AI-POWERED ORGANIZATION
 │   ├── prepare_files_for_llm()  # Prepare data
 │   └── generate_llm_prompt()    # Create prompt
 │
 ├─ demo_file_manager.py          # DEMO/TESTING
 ├─ __init__.py                   # Module init
 ├─ __pycache__/
 └─ README.md                     # Documentation
```

### **Troubleshooter** (`troubleshooter/`)
```
troubleshooter/
 ├─ screenshot_tool.py           # SCREENSHOT CAPTURE
 │   └── capture_screen_base64() # PNG to base64
 │
 ├─ vision_analyzer.py           # SCREEN ANALYSIS
 │   ├── VisionAnalyzer          # Main class
 │   ├── analyze_screenshot_for_issues()  # Analyze issues
 │   └── detect_ui_elements()    # UI detection
 │
 ├─ solution_parser.py           # SOLUTION PARSING
 │   └── parse_solution()        # LLM response to JSON
 │
 ├─ auto_fix_engine.py           # AUTO-FIX EXECUTION
 │   ├── execute_fixes()         # Run commands
 │   ├── run_command()           # Execute PowerShell
 │   ├── run_command_as_admin()  # Admin execution
 │   └── is_admin()              # Check privileges
 │
 └─ __pycache__/
```

### **Vision Engine** (`vision/`)
```
vision/
 ├─ vision_engine.py             # VISION AI (181 lines)
 │   ├── capture_screen()        # Take screenshot
 │   ├── analyze_screen()        # Send to Gemini
 │   ├── run_vision_pipeline()   # Full pipeline
 │   ├── ensure_output_dirs()    # Create folders
 │   ├── VISION_SYSTEM_PROMPT    # Analysis prompt
 │   └── VISION_OUTPUT_DIR       # Output location
 │
 ├─ __init__.py                  # Module init
 ├─ __pycache__/
 └─ README.md                    # Documentation
```

### **Personalization** (`personalisation/`)
```
personalisation/
 ├─ personalisation_tools.py     # PERSONALIZATION (641 lines)
 │   ├── toggle_dark_mode()      # Dark/light mode
 │   ├── get_current_theme()     # Get current theme
 │   ├── set_wallpaper()         # Set background
 │   ├── set_display_brightness()  # Brightness (0-100)
 │   ├── set_accent_color()      # Accent color
 │   ├── apply_theme_preset()    # Apply theme
 │   ├── save_personalization_profile()  # Save settings
 │   ├── load_personalization_profile()  # Load settings
 │   ├── manage_startup_apps()   # Startup apps
 │   └── set_default_app()       # Default app
 │
 ├─ demo_personalisation.py      # DEMO/TESTING
 ├─ __init__.py                  # Module init
 ├─ __pycache__/
 └─ README.md                    # Documentation
```

### **Installer** (`installer/`)
```
installer/
 ├─ installer.py                 # MAIN INSTALLER (760 lines)
 │   ├── InstallerAgent          # Main class
 │   ├── install_wallpaper()     # Download & set wallpaper
 │   ├── install_software()      # Install via winget
 │   ├── download_resource()     # Generic download
 │   └── auto_setup()            # Post-install setup
 │
 ├─ resource_finder.py           # RESOURCE SEARCH (307 lines)
 │   ├── ResourceFinder          # Main class
 │   ├── find_resource()         # Search online
 │   ├── search_wallpapers()     # Wallpaper search
 │   ├── search_software()       # Software search
 │   ├── search_images()         # Image search
 │   ├── _search_unsplash()      # Unsplash API
 │   ├── _search_pexels()        # Pexels API
 │   └── _search_pixabay()       # Pixabay API
 │
 ├─ downloader.py                # DOWNLOAD ENGINE (263 lines)
 │   ├── Downloader              # Main class
 │   ├── download_file()         # Download with progress
 │   ├── _download_with_resume() # Resume capability
 │   ├── _extract_filename()     # Get filename from URL
 │   ├── verify_file()           # Checksum verification
 │   └── max_retries = 3         # Retry logic
 │
 ├─ cache_manager.py             # CACHE SYSTEM (299 lines)
 │   ├── CacheManager            # Main class
 │   ├── get_cached_resource()   # Retrieve from cache
 │   ├── cache_resource()        # Store in cache
 │   ├── remove_cached_resource()  # Delete from cache
 │   ├── clear_cache()           # Full cache clear
 │   ├── get_cache_info()        # Cache statistics
 │   └── default_expiry_days = 30  # Expiration time
 │
 ├─ __init__.py                  # Module init
 └─ __pycache__/
```

### **Utilities** (`utils/`)
```
utils/
 ├─ json_parser.py               # JSON EXTRACTION
 │   └── extract_json()          # Parse LLM responses
 │
 ├─ prompt_loader.py             # PROMPT MANAGEMENT
 │   └── load_prompt()           # Load prompt templates
 │
 ├─ system_monitor.py            # SYSTEM MONITORING
 │   ├── get_resource_usage()    # Current usage
 │   └── generate_alerts()       # Threshold alerts
 │
 ├─ __init__.py                  # Module init
 └─ __pycache__/
```

### **Prompts** (`prompts/`)
```
prompts/
 ├─ system_prompt.txt            # System behavior guide
 └─ command_prompt.txt           # Command template
```

### **Memory & State** (`agent_memory/`)
```
agent_memory/
 ├─ execution_history.json       # All past executions [JSON Array]
 │   └── Contains: timestamp, plan_id, user_input, goal, summary
 │
 ├─ context.json                 # Conversation context [JSON Object]
 │   └── Contains: user messages, system responses, context
 │
 ├─ plans.json                   # Saved task plans [JSON Object]
 │   └── Contains: plan_id → TaskPlan mapping
 │
 └─ preferences.json             # User preferences [JSON Object]
    └── Contains: theme, auto_confirm, language, etc.
```

### **Vision Output** (`vision_output/`)
```
vision_output/
 ├─ analysis_*.json              # Vision analysis results
 ├─ ocr_analysis_*.json          # OCR text extraction
 └─ screenshots/
    └─ screen_*.png              # Captured screenshots
```

### **Logs** (`logs/`)
```
logs/
 ├─ jarvis.log                   # Main application log
 └─ [other service logs]
```

### **Testing** (`tests/`)
```
tests/
 ├─ test_master_agent.py         # Agent testing
 ├─ test_file_manager.py         # File manager testing
 ├─ test_file_management_integration.py  # Integration tests
 ├─ test_diagnostics.py          # Diagnostics testing
 ├─ test_gui_startup.py          # GUI testing
 ├─ TEST_COMMANDS.py             # Command examples
 └─ [other test files]
```

### **Archives** (`Archives/`)
```
Archives/
 ├─ run_ocr_demo.py              # Old OCR demo
 └─ vision_ocr.py                # Old vision code
```

### **Cleanup Test Directory** (`test_cleanup/`)
```
test_cleanup/
 ├─ large_files/                 # Test large files
 ├─ old_files/                   # Test old files
 │   ├─ old_doc1.txt
 │   └─ old_doc2.txt
 └─ temp_files/                  # Test temp files
```

---

## 📊 Intent Routing Reference Table

This comprehensive table shows every intent, its valid actions, handler function, and what tools it calls:

| **Intent** | **Actions** | **Handler Function** | **Tools Called** | **Requires Confirmation** |
|-----------|-----------|----------------------|------------------|----------------------|
| `help` | None | `_show_help()` | help_commands | No |
| `open_app` | None | App opener | app_tools | No |
| `web_search` | None | Web search | web_tools | No |
| `system_control` | sleep | `_handle_system_control()` | system_tools.sleep_pc() | Yes |
| `system_control` | lock | `_handle_system_control()` | system_tools.lock_pc() | Yes |
| `system_control` | kill_process | `_handle_system_control()` | system_tools.kill_process() | Yes |
| `system_control` | clean_temp | `_handle_system_control()` | system_tools.clean_temp() | No |
| `system_control` | empty_recycle_bin | `_handle_system_control()` | system_tools.empty_recycle_bin() | Yes |
| `diagnostics` | check_cpu | `_handle_diagnostics()` | diagnostics_tools.get_cpu_usage() | No |
| `diagnostics` | check_ram | `_handle_diagnostics()` | diagnostics_tools.get_ram_usage() | No |
| `diagnostics` | check_disk | `_handle_diagnostics()` | diagnostics_tools.get_disk_usage() | No |
| `diagnostics` | full_health_check | `_handle_diagnostics()` | diagnostics_tools.get_complete_system_health() | No |
| `disk_analysis` | analyze_usage | `_handle_disk_analysis()` | diagnostics_tools.analyze_disk_usage() | No |
| `disk_analysis` | find_large_folders | `_handle_disk_analysis()` | diagnostics_tools.find_large_folders() | No |
| `disk_analysis` | check_alerts | `_handle_disk_analysis()` | diagnostics_tools.check_health_alerts() | No |
| `maintenance` | scan_temp | `_handle_maintenance()` | diagnostics_tools.scan_cleanup_files() | No |
| `maintenance` | scan_old_files | `_handle_maintenance()` | diagnostics_tools.scan_cleanup_files() | No |
| `maintenance` | clean_temp_files | `_handle_maintenance()` | diagnostics_tools.clean_temp_files() | No |
| `maintenance` | clean_old_files | `_handle_maintenance()` | diagnostics_tools.clean_old_files() | No |
| `health_check` | None | `_handle_health_check()` | diagnostics_tools.get_complete_system_health() | No |
| `troubleshoot_screen` | None | `_handle_troubleshoot_screen()` | vision_analyzer + screenshot_tool | No |
| `vision_analysis` | analyze_screen | `_handle_vision_analysis()` | vision_engine.analyze_screen() | No |
| `vision_analysis` | detect_ui | `_handle_vision_analysis()` | vision_engine + gemini_client | No |
| `vision_analysis` | read_text | `_handle_vision_analysis()` | gemini_client.extract_text() | No |
| `personalization` | toggle_dark_mode | `_handle_personalization()` | personalisation_tools.toggle_dark_mode() | No |
| `personalization` | set_accent_color | `_handle_personalization()` | personalisation_tools.set_accent_color() | No |
| `personalization` | set_brightness | `_handle_personalization()` | personalisation_tools.set_display_brightness() | No |
| `personalization` | set_wallpaper | `_handle_personalization()` | personalisation_tools.set_wallpaper() | No |
| `personalization` | apply_preset | `_handle_personalization()` | personalisation_tools.apply_theme_preset() | No |
| `personalization` | save_profile | `_handle_personalization()` | personalisation_tools.save_personalization_profile() | No |
| `personalization` | load_profile | `_handle_personalization()` | personalisation_tools.load_personalization_profile() | No |
| `personalization` | manage_startup | `_handle_personalization()` | personalisation_tools.manage_startup_apps() | No |
| `personalization` | set_default_app | `_handle_personalization()` | personalisation_tools.set_default_app() | No |
| `file_management` | organize_by_type | `_handle_file_management()` | file_manager.organize_files() | No |
| `file_management` | find_duplicates | `_handle_file_management()` | file_manager.scan_duplicates() | No |
| `file_management` | remove_duplicates | `_handle_file_management()` | file_manager.remove_duplicates() | Yes |
| `file_management` | find_large_files | `_handle_file_management()` | file_manager.scan_large_files() | No |
| `file_management` | analyze_folder | `_handle_file_management()` | file_manager.analyze_folder() | No |
| `file_management` | get_report | `_handle_file_management()` | file_manager.get_comprehensive_report() | No |
| `system_config` | None | `_handle_system_config()` | system_config.get_system_config() | No |
| `installer` | install_software | `_handle_installer()` | installer.install_software() (winget) | Yes |
| `installer` | download_wallpaper | `_handle_installer()` | installer.install_wallpaper() | No |
| `installer` | download_software | `_handle_installer()` | installer.download_software() | No |
| `installer` | execute_installer | `_handle_installer()` | auto_fix_engine.execute_fixes() | Yes |
| `installer` | cache_info | `_handle_installer()` | cache_manager.get_cache_info() | No |
| `installer` | clear_cache | `_handle_installer()` | cache_manager.clear_cache() | Yes |

---

## 🔗 Dependency Map
```
main.py                          Entry point for CLI mode
run_jarvis.py                    Launcher script (CLI/GUI/Voice selector)
```

### **Agent Core** (`agent/`)
```
agent/
 ├─ master_agent.py            Central orchestrator (Main brain)
 ├─ task_planner.py            Decomposes goals into tasks
 ├─ execution_engine.py        Executes task plans
 ├─ memory_manager.py          Persistent storage & context
 ├─ task.py                    Task data structure
 └─ __init__.py
```

### **UI Layer** (`ui/`)
```
ui/
 ├─ app.py                     Main GUI application
 ├─ chat_card.py               Chat interface dialog
 ├─ __init__.py
 └─ README.md
```

### **Voice Input** (`voice/`)
```
voice/
 ├─ voice_listener.py          Microphone & speech-to-text
 ├─ __init__.py
 └─ __pycache__/
```

### **Routing** (`router/`)
```
router/
 ├─ intent_router.py           Intent dispatcher/executor
 └─ __pycache__/
```

### **LLM Clients** (`llm/`)
```
llm/
 ├─ groq_client.py             Groq API client (fast inference + Whisper)
 ├─ gemini_client.py           Google Gemini API client (vision)
 └─ __init__.py
```

### **Tools** (`tools/`)
```
tools/
 ├─ system_tools.py            System control (sleep, lock, kill process)
 ├─ app_tools.py               App launching
 ├─ web_tools.py               Web search
 ├─ diagnostics_tools.py       CPU, RAM, Disk monitoring
 ├─ help_commands.py           Help system
 ├─ system_config.py           System info
 └─ __init__.py
```

### **File Manager** (`file_manager/`)
```
file_manager/
 ├─ manager.py                 Main file manager class
 ├─ file_organizer.py          File organization by type
 ├─ duplicate_detector.py      Find duplicate files
 ├─ large_file_scanner.py      Find large files
 ├─ folder_finder.py           Recursive folder search
 ├─ llm_organizer.py           AI-powered organization
 ├─ demo_file_manager.py       Demo/testing
 ├─ __init__.py
 └─ README.md
```

### **Troubleshooter** (`troubleshooter/`)
```
troubleshooter/
 ├─ vision_analyzer.py         Screenshot analysis
 ├─ screenshot_tool.py         Screenshot capture
 ├─ solution_parser.py         Parse LLM solutions
 ├─ auto_fix_engine.py         Auto-apply fixes
 └─ __pycache__/
```

### **Vision** (`vision/`)
```
vision/
 ├─ vision_engine.py           Vision analysis pipeline
 ├─ __init__.py
 ├─ README.md
 └─ __pycache__/
```

### **Personalization** (`personalisation/`)
```
personalisation/
 ├─ personalisation_tools.py   Dark mode, themes, brightness
 ├─ demo_personalisation.py    Demo/testing
 ├─ __init__.py
 ├─ README.md
 └─ __pycache__/
```

### **Installer** (`installer/`)
```
installer/
 ├─ installer.py               App installation via winget
 ├─ resource_finder.py         Resource locator
 ├─ cache_manager.py           Package caching
 ├─ downloader.py              File downloader
 ├─ __init__.py
 └─ __pycache__/
```

### **Memory** (`agent_memory/`)
```
agent_memory/
 ├─ execution_history.json    All past executions
 ├─ context.json              Conversation context
 ├─ plans.json                Task plans
 └─ preferences.json          User preferences
```

### **Utilities** (`utils/`)
```
utils/
 ├─ json_parser.py            JSON extraction from LLM responses
 ├─ prompt_loader.py          Load prompt templates
 ├─ system_monitor.py         System monitoring
 └─ __pycache__/
```

### **Prompts** (`prompts/`)
```
prompts/
 ├─ system_prompt.txt         System behavior prompt
 └─ command_prompt.txt        Command template
```

---

## 🎮 How to Use JARVIS

### **1. Installation & Setup**

```bash
# Clone/navigate to project
cd c:\Users\ASUS\Desktop\Hackathon\Jarvis

# Activate virtual environment
.\.venv\Scripts\Activate.ps1

# Install dependencies
pip install -r requirements.txt
```

### **2. CLI Mode (Text)**

```bash
python main.py
```

**Usage:**
```
🤖 JARVIS - Intelligent OS Automation Agent
   ⌨️  TEXT MODE
============================================================
Commands: Type your request naturally, or try:
  - 'help' - Show help and commands
  - 'history' - Show execution history
  - 'stats' - Show execution statistics
  - 'status' - Show agent status
  - 'clear history' - Clear all memory
  - 'exit' - Exit JARVIS
============================================================

🤖 JARVIS> 
(Enter command)
```

**Examples:**
```
🤖 JARVIS> Check my CPU usage
🤖 JARVIS> Find duplicate files in Downloads
🤖 JARVIS> Organize my Downloads folder
🤖 JARVIS> Install Discord
🤖 JARVIS> Show system health
🤖 JARVIS> Clean temp files
```

### **3. GUI Mode (Chat Interface)**

```bash
python run_jarvis.py
```

**Features:**
- Floating bot widget
- Chat card interface
- Message history
- Copy/paste output

### **4. GUI + Voice Mode**

```bash
python run_jarvis.py --voice
```

**Usage:**
- Click the floating widget
- Open chat card
- Speak your command (or type)
- Auto-transcribed via Groq Whisper

### **5. Special Commands**

| Command | Purpose |
|---------|---------|
| `help` / `?` | Show available commands |
| `history` | Show recent executions |
| `stats` | Show execution statistics |
| `status` | Show agent status |
| `clear history` | Clear all memory |
| `exit` | Exit JARVIS |
| `calibrate` | Recalibrate voice silence threshold |

---

## 🔑 Key Concepts

### **Intent**
The **high-level goal** the user wants to accomplish.

Examples:
- `diagnostics` - Check system health
- `file_management` - Organize files
- `system_control` - Control PC
- `installer` - Install apps

### **Action**
The **specific operation** within an intent.

Examples:
- Intent: `diagnostics` → Action: `check_cpu`
- Intent: `file_management` → Action: `find_duplicates`
- Intent: `system_control` → Action: `sleep`

### **Parameters**
The **specific details** needed to execute an action.

Examples:
```python
{
  "intent": "file_management",
  "action": "organize_by_type",
  "parameters": {"folder_name": "Downloads"}    # Parameters
}
```

### **TaskPlan**
A list of **ordered, interdependent tasks** that execute a user goal.

Features:
- Sequential execution
- Dependency tracking
- Confirmation requirements
- Priority levels
- Rollback on failure

### **Execution History**
A log of **all past executions** stored in `agent_memory/execution_history.json`.

Used for:
- Learning from past patterns
- Statistics & reporting
- Debugging issues
- User analytics

### **Safety Checks**
Confirmation prompts for **potentially destructive actions**.

Examples:
- Deleting files
- Killing processes
- Empty recycle bin
- System sleep/lock

---

## 🔄 Data Flow Summary

```
USER INPUT
    ↓
MasterAgent (process_command)
    ├─ LLM: Extract goal
    ├─ Create context
    ↓
TaskPlanner (plan_goal)
    ├─ LLM: Decompose into tasks
    ├─ Create dependencies
    ↓
ExecutionEngine (execute_plan)
    ├─ Validate each task
    ├─ Get confirmations
    ├─ Execute sequentially
    ↓
IntentRouter (route_intent)
    ├─ Validate intent/action
    ├─ Route to handler
    ↓
Tools
    ├─ system_tools
    ├─ diagnostics_tools
    ├─ file_manager
    ├─ vision_engine
    ├─ personalisation_tools
    ├─ installer
    ↓
Results
    ↓
Output Display
    ├─ Terminal
    ├─ GUI Chat Card
    ├─ Log files
    ↓
MemoryManager
    └─ Save to agent_memory/
```

---

## 📊 Execution Statistics Example

```python
agent.get_execution_stats()

Output:
{
    "total_executions": 42,
    "successful": 38,
    "failed": 4,
    "success_rate": 90.48,
    "most_used_intents": {
        "diagnostics": 15,
        "file_management": 12,
        "system_control": 10
    },
    "average_execution_time": 12.5  # seconds
}
```

---

## 🛡️ Error Handling

JARVIS implements multi-level error handling:

### **Level 1: Input Validation**
```python
if not user_input or len(user_input.strip()) == 0:
    print("❌ Please provide a command")
    return
```

### **Level 2: Intent Validation**
```python
is_valid, msg = _validate_intent(intent, action)
if not is_valid:
    print(f"❌ {msg}")
    return
```

### **Level 3: Task Execution**
```python
try:
    result = tool_function(params)
except Exception as e:
    task.mark_failed()
    print(f"❌ Error: {str(e)}")
    logger.error(str(e))
```

### **Level 4: Recovery**
```python
if task.failed and task.retryable:
    retry_result = retry_task(task)  # Retry logic
```

---

---

## 🔐 Safety & Security Features

### **1. Permission System**

| Feature | Description |
|---------|-------------|
| **Admin Elevation** | Dangerous operations trigger UAC/admin prompt |
| **Dry-Run Mode** | Preview changes before executing |
| **Undo Capability** | Some operations (duplicates) keep originals |
| **Confirmation Gate** | User must approve destructive operations |

### **2. Protected Operations** (Require Confirmation)

```python
DANGEROUS_OPERATIONS = [
    "kill_process",           # Terminate applications
    "sleep_pc",               # Put PC to sleep
    "lock_pc",                # Lock workstation
    "empty_recycle_bin",      # Permanent deletion
    "remove_duplicates",      # Delete duplicate files
    "install_software",       # System modification
    "auto_fix_enable",        # Apply automatic fixes
    "clear_cache",            # Delete cache
    "format_drive" (future),  # Format entire drive
    "registry_modify" (future) # Registry changes
]
```

### **3. Audit Trail**

All operations logged to:
- `agent_memory/execution_history.json` - Detailed logs
- `logs/jarvis.log` - Application logs
- `vision_output/` - Screenshot analysis logs

### **4. Data Privacy**

- Screenshots stored locally only
- No cloud uploads without permission
- Passwords never logged
- Sensitive data masked in logs

### **5. Validation Layers**

```
Input → Sanitization → Intent Validation → 
Safety Check → Permissions → Execution → 
Logging → Error Recovery
```

---

## ⚡ Performance Metrics

### **Typical Execution Times**

| Operation | Time | Notes |
|-----------|------|-------|
| **Check CPU/Ram** | 0.5-1 sec | Real-time monitoring |
| **Check Disk** | 1-2 sec | Drive scanning |
| **Full Health Check** | 3-5 sec | All metrics |
| **Find Duplicates** | 10-30 sec | Depends on folder size |
| **Organize Files** | 5-15 sec | Depends on file count |
| **Find Large Files** | 5-10 sec | Fast scan |
| **Vision Analysis** | 3-8 sec | Gemini API latency |
| **Install Software** | 30-120 sec | Download + setup |
| **Download Wallpaper** | 2-5 sec | API + download |
| **Voice Transcription** | 2-4 sec | Groq Whisper |

### **Memory Usage** (Typical)

| Component | Memory |
|-----------|--------|
| **Base Python Process** | 50-80 MB |
| **Agent + LLM Clients** | 100-150 MB |
| **GUI (PyQt5)** | 150-200 MB |
| **Full Application** | 300-400 MB |

### **Network Usage**

| Operation | Data |
|-----------|------|
| **Voice Transcription** | ~500 KB per command |
| **Vision Analysis** | ~1-2 MB per screenshot |
| **Wallpaper Download** | 2-5 MB typical |
| **Software Install** | 50-500 MB (varies) |

---

## 🎓 Quick Reference by Use Case

### **Diagnostic Check-Up** 
```
User: "I want to check my PC health"
↓
Intent: diagnostics → full_health_check
↓
Shows: CPU %, RAM %, Disk %, Alerts
↓
Time: ~5 seconds
```

### **Free Up Disk Space**
```
User: "Clean up my disk"
↓
Intent: disk_analysis → find_large_folders
↓
Shows: Large folders + sizes
↓
Time: ~10 seconds
```

### **Organize Downloads**
```
User: "Organize my Downloads"
↓
Intent: file_management → organize_by_type
↓
Action: Creates folders (Images, Documents, Videos, etc.), moves files
↓
Time: 5-15 seconds depending on file count
```

### **Find & Remove Duplicates**
```
User: "Find duplicate files in Downloads and remove them"
↓
Step 1: Intent: file_management → find_duplicates
   └─ Shows duplicate groups + space that can be saved
↓
Step 2: User approves removal
↓
Step 3: Intent: file_management → remove_duplicates
   └─ Deletes duplicates, keeps master copy
↓
Time: 10-30 seconds total
```

### **Install an App**
```
User: "Install Discord"
↓
Intent: installer → install_software
↓
Step 1: Validates app exists in winget
Step 2: Shows confirmation (requires admin)
Step 3: Downloads & installs
Step 4: Verifies success
↓
Time: 30-120 seconds
```

### **Troubleshoot Screen Issues**
```
User: "My screen looks broken, can you fix it?"
↓
Step 1: Capture screenshot
Step 2: Send to Gemini for vision analysis
Step 3: Detect issues (flickering, color banding, etc.)
Step 4: Generate solutions
Step 5: Offer auto-fix
↓
Time: 5-10 seconds analysis + execution time
```

### **Change PC Theme**
```
User: "Switch to dark mode"
↓
Intent: personalization → toggle_dark_mode
↓
Action: Updates Windows Registry + UI theme
↓
Time: 2-3 seconds
```

---

## 🧪 Testing & Validation

### **Unit Tests** (`tests/`)

```python
# test_file_manager.py
def test_find_duplicates():
    """Verify duplicate detection works"""
    result = file_manager.find_duplicates(test_folder)
    assert result['duplicates_found'] == True
    assert len(result['duplicate_groups']) > 0

# test_diagnostics.py
def test_cpu_usage():
    """Verify CPU metrics are valid"""
    result = get_cpu_usage()
    assert 0 <= result['cpu_percent'] <= 100
    assert result['cpu_count'] > 0

# test_master_agent.py
def test_goal_extraction():
    """Verify LLM goal extraction"""
    result = agent.process_command("Check my CPU")
    assert result['goal'] is not None
    assert 'cpu' in result['goal'].lower()
```

### **Integration Tests** (`test_file_management_integration.py`)

```python
def test_complete_file_workflow():
    """Test full file management workflow"""
    # 1. Organize files
    org_result = file_manager.organize_files()
    assert org_result['files_moved'] > 0
    
    # 2. Find duplicates
    dup_result = file_manager.find_duplicates()
    
    # 3. Find large files
    large_result = file_manager.find_large_files()
    
    assert org_result['success'] == True
```

### **GUI Tests** (`test_gui_startup.py`)

```python
def test_gui_initialization():
    """Verify GUI starts without errors"""
    app = JarvisGUIApp()
    assert app is not None
    assert app.widget is not None
```

### **Running Tests**

```bash
# Run all tests
python -m pytest tests/

# Run specific test
python -m pytest tests/test_file_manager.py -v

# Run with coverage
python -m pytest tests/ --cov=./
```

---

## 🛠️ How to Add New Tools

### **Step 1: Create Tool Module**

File: `tools/my_new_tool.py`

```python
"""
New Tool Description

Usage:
    from tools.my_new_tool import my_function
    result = my_function(parameter)
"""

def my_function(parameter: str) -> dict:
    """
    What your function does
    
    Args:
        parameter: Description
        
    Returns:
        dict: {status, message, data}
    """
    try:
        # Your implementation
        result = do_something(parameter)
        
        return {
            "status": "success",
            "message": f"✅ Operation completed",
            "data": result
        }
    except Exception as e:
        return {
            "status": "error",
            "message": f"❌ {str(e)}",
            "data": None
        }
```

### **Step 2: Register Intent in Router**

File: `router/intent_router.py`

```python
# Add to VALID_INTENTS
VALID_INTENTS = {
    ...
    "my_tool": {
        "description": "What my tool does",
        "actions": ["action1", "action2"],
        "requires_admin": False
    }
}

# Add handler function
def _handle_my_tool(action: str, params: dict) -> dict:
    """Handler for my_tool intent"""
    if action == "action1":
        from tools.my_new_tool import my_function
        return my_function(params.get("parameter"))
    ...
```

### **Step 3: Add to Main Routing**

```python
elif intent == "my_tool":
    result = self._handle_my_tool(action, parameters)
```

### **Step 4: Test**

```python
# Test in CLI
python main.py
🤖 JARVIS> Use my new tool with action1

# Test programmatically
from router.intent_router import IntentRouter
router = IntentRouter()
result = router.route_intent({
    "intent": "my_tool",
    "action": "action1",
    "parameters": {"parameter": "test"}
})
```

### **Step 5: Document**

Add to WORKFLOW_GUIDE.md in appropriate layer section with:
- Function signature
- Parameters
- Output format
- Example usage

---

## 👁️ Vision Analysis Deep Dive

### **Vision Pipeline**

```
1. Capture Screen
   ↓
   Uses: mss (multiplatform screenshot)
   Output: PNG file
   
2. Convert to Base64
   ↓
   Uses: PIL ImageGrab → base64
   Output: String for API transmission
   
3. Send to Gemini
   ↓
   Uses: Google Gemini 2.5 Flash API
   Input: Base64 image + structured prompt
   
4. Parse Response
   ↓
   Uses: JSON extraction
   Output: Structured analysis
   
5. Save Results
   ↓
   Files: vision_output/{analysis,screenshots}/
   
6. Generate Recommendations
   ↓
   Output: Solutions + auto-fix options
```

### **Analysis Capabilities**

- **UI Detection:** Identifies buttons, text fields, menus
- **Text Extraction:** OCR via Gemini
- **Problem Detection:** Flickering, color issues, errors
- **Accessibility Check:** Validates UI accessibility
- **Screen Reader Output:** Generates descriptions

### **Prompts Used**

```python
VISION_SYSTEM_PROMPT = """
You are a Windows screen analysis expert.
Analyze the provided screenshot and:
1. Describe what's on screen
2. Identify all UI elements (buttons, text fields, etc.)
3. Detect any errors or issues
4. Suggest troubleshooting steps
"""
```

---

## 📊 Intent Statistics

**Total Intents Supported:** 14+

**Most Common Operations:**
1. Diagnostics (CPU, RAM, Disk) - 35% of commands
2. File management (organize, dedupe) - 25% of commands
3. System control (sleep, lock) - 15% of commands
4. Installation (apps, wallpapers) - 15% of commands
5. Personalization/Vision - 10% of commands

---

## 🎯 Summary

**JARVIS Workflow in One Sentence:**
> User input → MasterAgent extracts intent → TaskPlanner decomposes into tasks → ExecutionEngine runs tasks → IntentRouter dispatches to tools → Results stored in memory → Output displayed

**Key Files:**
- `main.py` / `run_jarvis.py` - Entry points
- `agent/master_agent.py` - Brain
- `agent/task_planner.py` - Planning
- `agent/execution_engine.py` - Execution
- `router/intent_router.py` - Routing
- `tools/` - Actual operations
- `file_manager/` - File operations
- `installer/` - App installation
- `troubleshooter/` - Auto-fixing
- `personalisation/` - System customization

**Three Input Modes:**
- CLI text input
- GUI chat interface
- Voice commands (Whisper API)

**Memory Persistence:**
- Execution history
- Context tracking
- User preferences
- Task plans

**Safety Layers:**
- Input validation
- Intent validation
- Safety checks
- Admin elevation
- Confirmation gates
- Audit trails

**Supported Operations:** 40+ distinct operations across 14+ intents

**Performance:** Most operations complete in <30 seconds

---

## 📝 Notes

- All operations are logged for debugging
- User confirmations required for dangerous operations
- Dependencies prevent task execution if parents fail
- LLM provides intelligent reasoning at two points: goal extraction & task planning  
- Vision analysis for troubleshooting via screenshots
- Support for 14+ intents covering most OS automation needs
- Modular architecture allows easy addition of new tools
- Full execution history enables learning and analytics
- Multi-layered error handling with recovery mechanisms

---

## 🔍 Documentation Index

| Section | Purpose | Page** |
|---------|---------|--------|
| Architecture Diagram | Visual overview | Top |
| Input Entry Points | Ways to interact | Input Section |
| Complete Workflow | Step-by-step process | Steps 1-8 |
| Core Components | System architecture | Components Section |
| Concrete Examples | Real-world usage | 4 Examples |
| Tools & Layers | All 12 tool layers | Layers 1-12 |
| File Structure | Directory reference | File Structure |
| Intent Routing | All 14+ intents | Intent Table |
| How to Use | Getting started | Usage Section |
| Safety Features | Security & validation | Safety Section |
| Performance Metrics | Speed benchmarks | Performance Section |
| Testing | Validation approaches | Testing Section |
| Add New Tools | Extension guide | Extension Section |

---

**Last Updated:** March 28, 2026  
**Version:** 3.0 (Complete)  
**Status:** ✅ COMPREHENSIVE - All details captured  
**Author:** JARVIS Development Team  
**Scope:** 100% of codebase referenced with detailed implementations
