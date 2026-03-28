# import os
# import json
# import time
# import speech_recognition as sr
# from gtts import gTTS
# import pygame
# from mss import mss
# from groq import Groq

# # Import your local vision engine
# from graph_engine import GraphExtractor

# # --- 1. Initialization ---
# print("Initializing Spectra Hybrid Cloud System...")

# # Initialize Pygame for playing audio files safely
# pygame.mixer.init()

# # Initialize STT
# recognizer = sr.Recognizer()
# microphone = sr.Microphone()

# # Initialize Groq Client
# # Ensure you generate an API key from console.groq.com and paste it below
# GROQ_API_KEY = "" 
# groq_client = Groq(api_key=GROQ_API_KEY)

# # Initialize Local Vision Engine
# extractor = GraphExtractor()

# # Global memory for the Agentic conversation
# chat_history = []

# # --- 2. Bilingual Audio & Cues ---

# def speak(text, lang='hi'):
#     """Generates speech using Google TTS and plays it."""
#     print(f"Spectra: {text}")
#     try:
#         # Generate audio file
#         tts = gTTS(text=text, lang=lang)
#         filename = f"temp_speech_{int(time.time())}.mp3"
#         tts.save(filename)
        
#         # Play audio file
#         pygame.mixer.music.load(filename)
#         pygame.mixer.music.play()
        
#         # Wait for audio to finish playing
#         while pygame.mixer.music.get_busy():
#             pygame.time.Clock().tick(10)
            
#         # Cleanup
#         pygame.mixer.music.unload()
#         os.remove(filename)
#     except Exception as e:
#         print(f"TTS Error: {e}")

# def play_cue(type="start"):
#     """Quick audio cues to let the user know the state."""
#     if type == "start":
#         print("🔔 [BEEP: Listening/Processing]")
#     elif type == "error":
#         print("❌ [BEEP: Error]")

# # --- 3. Agentic LLM Logic ---

# def ask_spectra_agent(user_prompt, graph_json=None):
#     """Sends the prompt and context to the Groq API (Llama 3)."""
#     global chat_history
    
#     if graph_json:
#         # Initialize a new context window for a new graph
#         chat_history = [{
#             "role": "system", 
#             "content": (
#                 "You are Spectra, an AI assistant for visually impaired users. "
#                 "You are friendly, concise, and conversational. "
#                 "You can understand and speak both Hindi and English. Match the language the user speaks. "
#                 f"Here is the raw data for the current graph on the user's screen: {graph_json}"
#             )
#         }]
    
#     chat_history.append({"role": "user", "content": user_prompt})
    
#     try:
#         completion = groq_client.chat.completions.create(
#             model="llama-3.3-70b-versatile", 
#             messages=chat_history,
#             temperature=0.5,
#             max_tokens=256,
#         )
        
#         reply = completion.choices[0].message.content
#         chat_history.append({"role": "assistant", "content": reply})
#         return reply
        
#     except Exception as e:
#         print(f"Groq API Error: {e}")
#         return "Network connection issue. Main online nahi hoon."

# # --- 4. Main Conversational Loop ---

# def listen_for_command():
#     with microphone as source:
#         print("\nListening...")
#         recognizer.adjust_for_ambient_noise(source, duration=0.5)
#         try:
#             # phrase_time_limit prevents it from hanging if background noise is continuous
#             audio = recognizer.listen(source, timeout=5, phrase_time_limit=8)
            
#             # Use 'hi-IN' which natively catches both Hindi and Indian English pronunciation perfectly
#             text = recognizer.recognize_google(audio, language="hi-IN").lower() 
#             print(f"User heard: {text}")
#             return text
#         except (sr.WaitTimeoutError, sr.UnknownValueError):
#             return ""
#         except Exception as e:
#             print(f"Microphone Error: {e}")
#             return ""

# def start_system():
#     play_cue("start")
#     speak("Spectra online. Mujhe 'read screen' ya 'screen dekho' bolein.", lang='hi')

#     in_conversation = False

#     while True:
#         command = listen_for_command()
#         if not command:
#             continue

#         # 1. Trigger the Vision Pipeline
#         if any(trigger in command for trigger in ["read screen", "screen dekho", "graph dekho", "analyze", "स्क्रीन देखो", "ग्राफ देखो", "रीड स्क्रीन"]):
#             play_cue("start")
#             speak("Screen check kar rahi hoon. Ek second dijiye.", lang='hi')
            
#             # Capture the screen locally
#             with mss() as sct:
#                 filename = sct.shot(mon=1, output='temp_screen.png')
            
#             # Run Local Vision Engine 
#             print("Extracting Graph Data Locally...")
#             raw_json = extractor.extract(filename)
            
#             # Send to Groq for summary 
#             print("Consulting Agentic Brain...")
#             prompt = "Summarize the key points of this graph briefly. If I spoke in Hindi, answer in Hindi."
#             explanation = ask_spectra_agent(prompt, graph_json=raw_json)
            
#             # Speak the result
#             speak(explanation, lang='hi') 
#             speak("Kya aapko is graph ke baare mein kuch aur poochhna hai?", lang='hi')
            
#             in_conversation = True
#             if os.path.exists(filename):
#                 os.remove(filename)

#         # 2. Agentic Q&A Loop
#         elif in_conversation:
#             if any(stop_word in command for stop_word in ["stop", "cancel", "chup", "band karo", "thank you", "चुप", "बंद करो", "रुक जाओ"]):
#                 speak("Theek hai. Main background mein hoon.", lang='hi')
#                 in_conversation = False
#                 chat_history = [] # Wipe memory
#             else:
#                 play_cue("start")
#                 answer = ask_spectra_agent(command)
#                 speak(answer, lang='hi')

# if __name__ == "__main__":
#     start_system()






















# import os
# import json
# import time
# import speech_recognition as sr
# from gtts import gTTS
# import pygame
# from mss import mss
# from groq import Groq
# import tkinter as tk
# from tkinter import scrolledtext, messagebox
# import threading
# import os
# import sys

# # Import your Hugging Face Vision Engine
# from graph_engine import GraphExtractor

# # --- 1. Initialization ---
# print("Initializing Spectra Hybrid Cloud System...")

# # Initialize Pygame for playing audio files safely
# pygame.mixer.init()

# # Initialize STT
# recognizer = sr.Recognizer()
# microphone = sr.Microphone()

# # Initialize Groq Client
# # INSERT YOUR GROQ API KEY HERE
# GROQ_API_KEY = "" 
# try:
#     groq_client = Groq(api_key=GROQ_API_KEY)
# except Exception as e:
#     print(f"Error initializing Groq: {e}. Please check your API key.")

# # Initialize Vision Engine (This will download weights from Hugging Face on the first run)
# extractor = GraphExtractor()

# # Global memory for the Agentic conversation
# chat_history = []

# # --- 2. Bilingual Audio & Cues ---

# def speak(text, lang='hi'):
#     """Generates speech using Google TTS and plays it."""
#     print(f"\nSpectra: {text}")
#     try:
#         # Generate audio file
#         tts = gTTS(text=text, lang=lang)
#         filename = f"temp_speech_{int(time.time())}.mp3"
#         tts.save(filename)
        
#         # Play audio file
#         pygame.mixer.music.load(filename)
#         pygame.mixer.music.play()
        
#         # Wait for audio to finish playing
#         while pygame.mixer.music.get_busy():
#             pygame.time.Clock().tick(10)
            
#         # Cleanup
#         pygame.mixer.music.unload()
#         os.remove(filename)
#     except Exception as e:
#         print(f"TTS Error: {e}")

# def play_cue(type="start"):
#     """Quick audio cues to let the user know the state."""
#     if type == "start":
#         print("🔔 [BEEP: Listening/Processing]")
#     elif type == "error":
#         print("❌ [BEEP: Error]")

# # --- 3. Agentic LLM Logic ---

# def ask_spectra_agent(user_prompt, graph_json=None):
#     """Sends the prompt and context to the Groq API (Llama 3.3)."""
#     global chat_history
    
#     if graph_json:
#         # Initialize a new context window for a new graph
#         chat_history = [{
#             "role": "system", 
#             "content": (
#                 "You are Spectra, an AI assistant for visually impaired users. "
#                 "You are friendly, concise, and conversational. "
#                 "You can understand and speak both Hindi and English. Match the language the user speaks. "
#                 "Do not read the raw JSON to the user. Explain the insights naturally. "
#                 f"Here is the data for the current graph on the user's screen: {graph_json}"
#             )
#         }]
    
#     chat_history.append({"role": "user", "content": user_prompt})
    
#     try:
#         # Using the latest supported Groq model
#         completion = groq_client.chat.completions.create(
#             model="llama-3.3-70b-versatile", 
#             messages=chat_history,
#             temperature=0.5,
#             max_tokens=300,
#         )
        
#         reply = completion.choices[0].message.content
#         chat_history.append({"role": "assistant", "content": reply})
#         return reply
        
#     except Exception as e:
#         print(f"Groq API Error: {e}")
#         return "Mujhe network se connect karne mein problem ho rahi hai. Kripya apna internet check karein."

# # --- 4. Main Conversational Loop ---

# def listen_for_command():
#     with microphone as source:
#         print("\nListening...")
#         # Quick calibration for background noise
#         recognizer.adjust_for_ambient_noise(source, duration=0.5)
#         try:
#             audio = recognizer.listen(source, timeout=5, phrase_time_limit=8)
#             # hi-IN captures both spoken Hindi, English, and Hinglish
#             text = recognizer.recognize_google(audio, language="hi-IN").lower() 
#             print(f"User heard: {text}")
#             return text
#         except (sr.WaitTimeoutError, sr.UnknownValueError):
#             return ""
#         except Exception as e:
#             print(f"Microphone Error: {e}")
#             return ""

# def start_system():
#     play_cue("start")
#     speak("Spectra online. Mujhe 'read screen' ya 'screen dekho' bolein.", lang='hi')

#     in_conversation = False

#     while True:
#         command = listen_for_command()
#         if not command:
#             continue

#         # 1. Trigger the Vision Pipeline
#         # Now natively handles both English and actual Devanagari text
#         triggers = ["read screen", "screen dekho", "graph dekho", "analyze", "स्क्रीन देखो", "ग्राफ देखो", "रीड स्क्रीन"]
        
#         if any(trigger in command for trigger in triggers):
#             play_cue("start")
#             speak("Screen check kar rahi hoon. Ek second dijiye.", lang='hi')
            
#             # Capture the screen locally
#             with mss() as sct:
#                 filename = sct.shot(mon=1, output='temp_screen.png')

#             # # TEMPORARY TEST OVERRIDE:
#             # filename = "sample1.png"
            
#             # Run Vision Engine (Downloads from HF if first time, otherwise uses cache)
#             print("\nExtracting Graph Data...")
#             raw_json = extractor.extract(filename)
            
#             # Send to Groq for summary 
#             print("Consulting Agentic Brain...")
#             prompt = "Please summarize the key takeaways of this graph. If I spoke in Hindi, answer in Hindi."
#             explanation = ask_spectra_agent(prompt, graph_json=raw_json)
            
#             # Speak the result
#             speak(explanation, lang='hi') 
#             speak("Kya aapko is graph ke baare mein kuch aur poochhna hai?", lang='hi')
            
#             in_conversation = True
#             if os.path.exists(filename):
#                 os.remove(filename)

#         # 2. Agentic Q&A Loop
#         elif in_conversation:
#             stop_words = ["stop", "cancel", "chup", "band karo", "thank you", "चुप", "बंद करो", "रुक जाओ"]
#             if any(stop_word in command for stop_word in stop_words):
#                 speak("Theek hai. Main background mein dhyan rakh rahi hoon.", lang='hi')
#                 in_conversation = False
#                 chat_history = [] # Wipe memory for privacy
#             else:
#                 play_cue("start")
#                 answer = ask_spectra_agent(command)
#                 speak(answer, lang='hi')

# #[Section 3: The New Chat UI Class]
# class SpectraChatUI:

#     def update_chat(self, sender, message):
#         self.chat_area.config(state=tk.NORMAL)
#         color = "#55ff55" if sender == "You" else "#55ccff"
#         self.chat_area.insert(tk.END, f"{sender}: ", tag:=sender)
#         self.chat_area.tag_config(tag, foreground=color, font=("Arial", 10, "bold"))
#         self.chat_area.insert(tk.END, f"{message}\n\n")
#         self.chat_area.config(state=tk.DISABLED)
#         self.chat_area.yview(tk.END)

#     def start_thread(self):
#         if not self.is_running:
#             self.is_running = True
#             self.update_chat("System", "Spectra Online. Say 'Screen Dekho'.")
#             # This launches the voice loop without freezing the window
#             threading.Thread(target=self.run_spectra, daemon=True).start()

#     def run_spectra(self):
#         # This links the UI to your voice functions
#         global speak
#         original_speak = speak
#         def ui_speak(text, lang='hi'):
#             self.root.after(0, self.update_chat, "Spectra", text)
#             original_speak(text, lang)
#         speak = ui_speak
        
#         # Start your main conversational loop
#         start_system()

# # [Section 4: The Main Execution]
# if __name__ == "__main__":
#     root = tk.Tk()
#     gui = SpectraChatUI(root)
#     root.mainloop()


import os
import time
import speech_recognition as sr
from gtts import gTTS
import pygame
from mss import mss
from groq import Groq
import tkinter as tk
from tkinter import scrolledtext
import threading
import sys

# Import your Vision Engine
from graph_engine import GraphExtractor

# --- GLOBAL SETUP ---
pygame.mixer.init()
recognizer = sr.Recognizer()
microphone = sr.Microphone()
GROQ_API_KEY = "" 
groq_client = Groq(api_key=GROQ_API_KEY)
extractor = GraphExtractor()
chat_history = []

def original_speak(text, lang='hi'):
    try:
        tts = gTTS(text=text, lang=lang)
        filename = f"temp_{int(time.time())}.mp3"
        tts.save(filename)
        pygame.mixer.music.load(filename)
        pygame.mixer.music.play()
        while pygame.mixer.music.get_busy(): pygame.time.Clock().tick(10)
        pygame.mixer.music.unload()
        os.remove(filename)
    except Exception as e: print(f"Audio Error: {e}")

speak = original_speak

def ask_spectra_agent(user_prompt, graph_json=None):
    global chat_history
    if graph_json:
        chat_history = [{"role": "system", "content": f"You are Spectra. Data: {graph_json}"}]
    chat_history.append({"role": "user", "content": user_prompt})
    completion = groq_client.chat.completions.create(
        model="llama-3.3-70b-versatile", messages=chat_history, temperature=0.5
    )
    reply = completion.choices[0].message.content
    chat_history.append({"role": "assistant", "content": reply})
    return reply

# --- THE UI CLASS ---
class SpectraChatUI:
    def __init__(self, root_window):
        self.root = root_window
        self.root.title("Spectra Agentic AI")
        self.root.geometry("400x650")
        self.root.configure(bg="#1e1e1e")
        self.is_running = False

        # Status Indicator (The "Listening" Dialog)
        self.status_frame = tk.Frame(self.root, bg="#1e1e1e")
        self.status_frame.pack(fill=tk.X, pady=5)
        self.status_light = tk.Label(self.status_frame, text="●", fg="#444444", bg="#1e1e1e", font=("Arial", 14))
        self.status_light.pack(side=tk.LEFT, padx=(10, 0))
        self.status_text = tk.Label(self.status_frame, text="OFFLINE", fg="white", bg="#1e1e1e", font=("Arial", 10))
        self.status_text.pack(side=tk.LEFT, padx=5)

        # Chat Area
        self.chat_area = scrolledtext.ScrolledText(self.root, bg="#2d2d2d", fg="white", font=("Segoe UI", 10))
        self.chat_area.pack(padx=10, pady=10, fill=tk.BOTH, expand=True)
        
        # Buttons
        btn_frame = tk.Frame(self.root, bg="#1e1e1e")
        btn_frame.pack(pady=10)
        tk.Button(btn_frame, text="START ENGINE", command=self.start_thread, bg="#0078d7", fg="white").pack(side=tk.LEFT, padx=5)
        tk.Button(btn_frame, text="EXIT", command=lambda: os._exit(0), bg="#444444", fg="white").pack(side=tk.LEFT, padx=5)

    def set_status(self, state):
        """Updates the visual 'Listening' light."""
        if state == "listening":
            self.status_light.config(fg="#ff5555") # Red for recording
            self.status_text.config(text="LISTENING...")
        elif state == "processing":
            self.status_light.config(fg="#55ff55") # Green for AI thinking
            self.status_text.config(text="PROCESSING...")
        else:
            self.status_light.config(fg="#444444")
            self.status_text.config(text="IDLE")

    def update_chat(self, sender, message):
        self.chat_area.config(state=tk.NORMAL)
        color = "#55ff55" if sender == "You" else "#55ccff"
        self.chat_area.insert(tk.END, f"{sender}: ", sender)
        self.chat_area.tag_config(sender, foreground=color, font=("Arial", 10, "bold"))
        self.chat_area.insert(tk.END, f"{message}\n\n")
        self.chat_area.config(state=tk.DISABLED)
        self.chat_area.yview(tk.END)

    def start_thread(self):
        if not self.is_running:
            self.is_running = True
            threading.Thread(target=self.run_logic, daemon=True).start()

    def run_logic(self):
        global speak
        def ui_speak(text, lang='hi'):
            self.root.after(0, self.set_status, "idle")
            self.root.after(0, self.update_chat, "Spectra", text)
            original_speak(text, lang)
        speak = ui_speak
        
        speak("Spectra online. Screen dekho bolein.", lang='hi')
        in_conversation = False

        while True:
            self.root.after(0, self.set_status, "listening")
            with microphone as source:
                recognizer.adjust_for_ambient_noise(source, duration=0.5)
                try:
                    audio = recognizer.listen(source, timeout=None, phrase_time_limit=5)
                    command = recognizer.recognize_google(audio, language="hi-IN").lower()
                    self.root.after(0, self.update_chat, "You", command)
                except: continue

            # --- STOP COMMANDS ---
            stop_triggers = ["stop", "cancel", "chup", "band karo", "चुप", "बंद करो"]
            if any(s in command for s in stop_triggers):
                speak("Theek hai, main chup ho rahi hoon.", lang='hi')
                in_conversation = False
                continue

            # --- ANALYZE COMMANDS ---
            if any(t in command for t in ["read screen", "screen dekho", "graph dekho", "analyze", "स्क्रीन देखो", "ग्राफ देखो", "रीड स्क्रीन"]):
                self.root.after(0, self.set_status, "processing")
                speak("Scanning...", lang='hi')
                with mss() as sct:
                    fname = sct.shot(mon=1, output='temp_screen.png')
                raw_json = extractor.extract(fname)
                explanation = ask_spectra_agent("Summarize.", graph_json=raw_json)
                speak(explanation, lang='hi')
                in_conversation = True
            
            # --- CONVERSATION ---
            elif in_conversation and len(command) > 2:
                self.root.after(0, self.set_status, "processing")
                answer = ask_spectra_agent(command)
                speak(answer, lang='hi')

# --- ENTRY POINT ---
if __name__ == "__main__":
    main_window = tk.Tk()
    app_instance = SpectraChatUI(main_window)
    main_window.mainloop()