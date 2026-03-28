## VITALS: Agentic AI Healthcare
**Team UNFAZED**

---

### 👥 Team Members

| Name | Role | GitHub |
| :--- | :--- | :--- |
| **Aryan Saini** | Team Leader / Full Stack | @username |
| **Aryan Gusain** | Backend / AI Logic | @username |
| **Ansh Thakur** | Frontend / UI Design | @username |
| **Anshika Garg** | Documentation / Testing | @username |

---

### 💡 Problem Statement
Current healthcare is reactive and manual, leading to three primary systemic failures:

* **Ignoring Symptoms:** Chronic patients often overlook subtle symptom links (e.g., a "metallic taste" in Type 2 Diabetes indicating a shift to Type 3), which leads to preventable emergencies.
* **Staff Overload:** Medical professionals are too burdened to provide continuous manual monitoring for every chronic patient at home.
* **Efficiency Loss:** Manual history-taking consumes the majority of a patient's visit, leaving minimal time for actual treatment and consultation.

---

### 🛠️ Tech Stack

* **Core Intelligence:** Gemini API (Agentic Reasoning & Pattern Recognition)
* **Orchestration:** n8n (Workflow Automation) & LangChain
* **Voice Pipeline:** Browser Web Speech API (Mocking Twilio/Vapi/ElevenLabs for web-based calling)
* **Database:** Supabase (Patient History & Vector DB)
* **Transcription:** Deepgram (Speech-to-Text)
* **Frontend:** React/Vite with Tailwind CSS

---

### 🔗 Links
* **Live Demo:** [Insert Link Here]
* **Video Demo:** [Insert Link Here]
* **Presentation:** [Link to Team Unfazed Stellaris.pdf]

---

### 📸 Screenshots
*(Placeholder for UI screenshots: Highlighting the **Call Interface**, **Patient History Sidebar**, and the **Doctor Approval Dashboard**)*

---

### 🚀 How to Run Locally

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/username/vitals-agentic-ai.git
    ```

2.  **Install Dependencies:**
    ```bash
    npm install
    ```

3.  **Environment Setup:**
    Create a `.env` file and add your **Gemini API Key** and **Supabase URL/Key**.

4.  **Database Sync:**
    Connect your existing patient database to the workflow as defined in the technical approach.

6.  **Launch Frontend:**
    ```bash
    npm run dev
    ```

7.  **Initiate Mock Call:**
    Use the web interface to start a voice session. The system will autonomously retrieve patient history and begin context-aware questioning.

---