import os
import json
from fastapi import FastAPI, File, UploadFile
from fastapi.responses import PlainTextResponse
from fastapi.middleware.cors import CORSMiddleware
from groq import Groq

# Import your newly updated PyPI library!
from graphvision import GraphExtractor

app = FastAPI(title="STEM Sight Backend")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"], # Allows any browser extension to connect
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize the Groq Client (Looks for the GROQ_API_KEY environment variable)
groq_client = Groq()

# Initialize your custom PyPI library
print("Initializing STEM Sight Vision Engine...")
vision_engine = GraphExtractor()

@app.get("/")
async def root():
    return {"message": "STEM Sight API is online and ready."}

# @app.post("/analyze-graph", response_class=PlainTextResponse)
# async def analyze_graph(file: UploadFile = File(...)):
#     try:
#         # 1. Save the uploaded image temporarily
#         temp_image_path = f"temp_{file.filename}"
#         with open(temp_image_path, "wb") as buffer:
#             buffer.write(await file.read())

#         # 2. Extract structured data using your library
#         print(f"Extracting data from {file.filename}...")
        
#         # 🚨 UPDATED: Call the new extract method
#         extraction_json_string = vision_engine.extract(temp_image_path)
        
#         # Clean up the temporary file immediately
#         if os.path.exists(temp_image_path):
#             os.remove(temp_image_path)

#         # 🚨 UPDATED: Parse the JSON string back into a Python dictionary
#         extraction_result = json.loads(extraction_json_string)

#         # 🚨 UPDATED: Check for the new error format from your library
#         if "error" in extraction_result:
#             return f"I'm sorry, I couldn't clearly identify the data in this graph. Reason: {extraction_result['error']}"

#         # 3. Format the JSON data into a prompt
#         graph_type = extraction_result.get("chart_type", "unknown")
#         graph_data = extraction_result.get("data", [])
        
#         # Grab optional labels/titles if they exist (good for context!)
#         x_label = extraction_result.get("x_axis_label", "Unknown X-Axis")
#         y_label = extraction_result.get("y_axis_label", "Unknown Y-Axis")
#         title = extraction_result.get("title", "Untitled Graph")
        
#         prompt = f"""
#         You are an accessibility assistant for visually impaired students. 
#         I am giving you extracted data from a {graph_type} chart. 
#         Title: {title}
#         X-Axis Label: {x_label}
#         Y-Axis Label: {y_label}
        
#         Please summarize this data in one short, conversational, and easy-to-understand paragraph.
#         Point out the largest and smallest values if relevant.
#         Do not use markdown, bold text, or asterisks. Write it exactly as it should be spoken out loud by a text-to-speech engine.
        
#         Data:
#         {graph_data}
#         """

#         # 4. Send to Groq for lightning-fast inference
#         print("Generating audio script with Groq Llama 3...")
#         chat_completion = groq_client.chat.completions.create(
#             messages=[
#                 {
#                     "role": "user",
#                     "content": prompt,
#                 }
#             ],
#             model="llama-3.1-8b-instant",
#             temperature=0.4, # Lowered slightly for more factual summaries
#         )

#         # 5. Return strictly the text response for the Chrome extension to speak
#         return chat_completion.choices[0].message.content.strip()

#     except Exception as e:
#         return f"An error occurred while analyzing the graph: {str(e)}"

# if __name__ == "__main__":
#     import uvicorn
#     # Runs the server on port 8000
#     uvicorn.run(app, host="0.0.0.0", port=8000)


import time # Add this to your imports

@app.post("/analyze-graph", response_class=PlainTextResponse)
async def analyze_graph(file: UploadFile = File(...)):
    try:
        start_time = time.time()
        
        # 1. Save the uploaded image temporarily
        temp_image_path = f"temp_{file.filename}"
        with open(temp_image_path, "wb") as buffer:
            buffer.write(await file.read())
            
        print(f"⏱️ Image received and saved in: {time.time() - start_time:.2f} seconds")

        # 2. Extract structured data
        extract_start = time.time()
        print(f"Extracting data from {file.filename}...")
        extraction_json_string = vision_engine.extract(temp_image_path)
        print(f"⏱️ AI Extraction finished in: {time.time() - extract_start:.2f} seconds")
        
        if os.path.exists(temp_image_path):
            os.remove(temp_image_path)

        extraction_result = json.loads(extraction_json_string)
        print(f"Extracted data: {extraction_result}")

        if "error" in extraction_result:
            return f"I'm sorry, I couldn't clearly identify the data in this graph. Reason: {extraction_result['error']}"

        graph_type = extraction_result.get("chart_type", "unknown")
        graph_data = extraction_result.get("data", [])
        x_label = extraction_result.get("x_axis_label", "Unknown X-Axis")
        y_label = extraction_result.get("y_axis_label", "Unknown Y-Axis")
        title = extraction_result.get("title", "Untitled Graph")
        
        prompt = f"""
        You are an accessibility assistant for visually impaired students. 
        I am giving you extracted data from a {graph_type} chart. 
        Title: {title}
        X-Axis Label: {x_label}
        Y-Axis Label: {y_label}
        
        Please summarize this data in one short, conversational, and easy-to-understand paragraph.
        Point out the largest and smallest values if relevant.
        Do not use markdown, bold text, or asterisks. Write it exactly as it should be spoken out loud by a text-to-speech engine.
        
        Data:
        {graph_data}
        """

        # 3. Send to Groq
        groq_start = time.time()
        print("Generating audio script with Groq Llama 3...")
        chat_completion = groq_client.chat.completions.create(
            messages=[{"role": "user", "content": prompt}],
            model="llama-3.1-8b-instant",
            temperature=0.4,
        )
        print(f"⏱️ Groq Llama 3 finished in: {time.time() - groq_start:.2f} seconds")
        print(f"✅ TOTAL TIME: {time.time() - start_time:.2f} seconds")

        return chat_completion.choices[0].message.content.strip()

    except Exception as e:
        return f"An error occurred while analyzing the graph: {str(e)}"