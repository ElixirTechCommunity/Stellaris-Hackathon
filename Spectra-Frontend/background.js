const API_URL = "https://shadowgard3n-spectra-backend.hf.space/analyze-graph";

// 1. Listen for the icon click
chrome.action.onClicked.addListener((tab) => {
  console.log("🔘 Extension icon clicked on URL:", tab.url);
  
  chrome.scripting.executeScript({
    target: { tabId: tab.id, allFrames: false },
    files: ["content.js"]
  })
  .then(() => console.log("✅ content.js injected successfully!"))
  .catch((err) => console.error("❌ Failed to inject content.js:", err));
});

// 2. Listen for image URLs from the webpage
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "analyzeImage") {
    console.log("🖼️ Received image to analyze:", request.src);
    
    processImage(request.src)
      .then(explanation => sendResponse({ text: explanation }))
      .catch(error => {
        console.error("❌ Background Fetch Error:", error);
        sendResponse({ text: "Sorry, I could not process this image." });
      });
      
    return true; // Keeps the message channel open while waiting for FastAPI
  }
});

// 3. Send to Hugging Face
async function processImage(imgSrc) {
  const imageResponse = await fetch(imgSrc);
  const imageBlob = await imageResponse.blob();
  const formData = new FormData();
  formData.append("file", imageBlob, "image.png");

  const apiResponse = await fetch(API_URL, {
    method: "POST",
    body: formData
  });

  if (!apiResponse.ok) throw new Error(`API Error: ${apiResponse.status}`);
  return await apiResponse.text();
}