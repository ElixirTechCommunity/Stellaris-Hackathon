# 🤖 Kiro - Intelligent OS Automation Agent

## 👥 Team Name
**Horizon**

## 🧑‍💻 Team Members
| Name | Role | GitHub |
|------|------|--------|
| Saubhgya Sharma | Full Stack | @somcpp |
| Prashant Gupta | Backend | @PrashantGuptacodecraft |
| Saurabh Rajput | Backend | @Dev-saurabhraj |
| Saurya Tripathi | Frontend | @saurya992004 |

## 💡 Problem Statement
Modern users struggle with repetitive, time-consuming desktop tasks. Kiro solves this by providing an intelligent AI-powered OS automation agent that understands natural language commands and executes complex system tasks with safety, reliability, and automatic learning. 

The system intelligently breaks down tasks into executable subtasks, validates commands before execution, manages files efficiently, performs system diagnostics, integrates vision/OCR capabilities, and maintains a persistent memory of user preferences and execution history—all while prioritizing safety through comprehensive action logging and automatic rollback mechanisms.

## 🛠️ Tech Stack
- **Core Framework:** Python 3.x, LangGraph
- **LLMs:** Google Gemini, Groq
- **UI Frameworks:** CustomTkinter, PyQt5, FastAPI
- **Computer Vision:** OpenCV, EasyOCR, Tesseract, PIL
- **System Monitoring:** Psutil
- **Audio Processing:** SoundDevice, Numpy, SciPy
- **Web Stack:** FastAPI, Uvicorn, WebSockets
- **API Communication:** Requests
- **Configuration:** Python-dotenv
- **Graph Visualization:** LangGraph
- **Language Chain:** LangChain-Core

## 🔗 Links
- **Live Demo:** [Not Available]
- **Video Demo:** [https://www.youtube.com/watch?v=DDbIej8rzKk](https://www.youtube.com/watch?v=DDbIej8rzKk)
- **Presentation (PPT/PDF):** [Link]

## 📸 Screenshots
<!-- Screenshots to be added -->

## 🚀 How to Run Locally

### Prerequisites
- Python 3.8 or higher
- pip package manager
- Internet connection for LLM APIs (Gemini, Groq)

### Installation Steps

1. **Clone the Repository**
   ```bash
   cd c:\Users\ASUS\Desktop\Hackathon\Stellaris-Hackathon\submissions\Horizon_KIRO
   ```

2. **Create a Virtual Environment** (Optional but recommended)
   ```bash
   python -m venv venv
   
   # On Windows
   venv\Scripts\activate
   
   # On macOS/Linux
   source venv/bin/activate
   ```

3. **Install Dependencies**
   ```bash
   pip install -r requirements.txt
   ```

4. **Set Up Environment Variables**
   Create a `.env` file in the project root with your API keys:
   ```
   GEMINI_API_KEY=your_google_gemini_api_key
   GROQ_API_KEY=your_groq_api_key
   ```

5. **Run the Application**

   **Option A: CLI Interface**
   ```bash
   python main.py
   ```

   **Option B: GUI Interface**
   ```bash
   python ui/app.py
   ```

   **Option C: Start with JARVIS**
   ```bash
   python run_jarvis.py
   ```

### Key Features to Try

- **File Management:** Organize files, detect duplicates, find large files
- **System Diagnostics:** Monitor CPU, RAM, disk usage and system health
- **Task Planning:** Give complex commands and watch the agent break them down
- **Vision Analysis:** Upload screenshots for AI-powered analysis
- **Natural Language Commands:** Use conversational language to control your system

### Project Structure
- `main.py` - Main entry point
- `agent/` - Core agent components (master agent, task planner, execution engine)
- `file_manager/` - File organization and management tools
- `llm/` - LLM client integrations (Gemini, Groq)
- `tools/` - System, file, diagnostic, and app tools
- `ui/` - User interface components
- `vision/` - Vision and OCR capabilities
- `tests/` - Test suite
- `prompts/` - System and command prompts
- `agent_memory/` - Persistent memory storage

---

**For detailed workflow and architecture information, see [WORKFLOW_GUIDE.md](WORKFLOW_GUIDE.md)**
