(function() {
    // 1. THE TOGGLE SWITCH: If clicked while already running, stop it entirely!
    if (window.STEM_SIGHT_ACTIVE) {
        console.log("🛑 Extension clicked again. Stopping STEM Sight.");
        window.STEM_SIGHT_ACTIVE = false;
        window.speechSynthesis.cancel();
        return; 
    }

    // 2. Turn it on
    window.STEM_SIGHT_ACTIVE = true;
    console.log("🧠 STEM Sight content script has booted up!");

    // 3. The Panic Button: Press 'Escape' to stop the AI completely
    document.addEventListener('keydown', function escapeListener(event) {
        if (event.key === 'Escape') {
            console.log("🛑 Escape key pressed. Stopping STEM Sight.");
            window.STEM_SIGHT_ACTIVE = false;
            window.speechSynthesis.cancel();
            document.removeEventListener('keydown', escapeListener); // Clean up
        }
    });

    // 4. Ultra-Safe Speech Engine
    function speakAndWait(text) {
        if (!window.STEM_SIGHT_ACTIVE) return Promise.resolve(); 

        console.log("🔊 Attempting to speak:", text); 
        
        return new Promise((resolve) => {
            window.speechSynthesis.cancel(); 

            setTimeout(() => {
                if (!window.STEM_SIGHT_ACTIVE) return resolve(); 

                // Attach to 'window' so Chrome's memory manager NEVER deletes the voice
                window.stemSightUtterance = new SpeechSynthesisUtterance(text);
                window.stemSightUtterance.rate = 1.0; 
                
                // Safety timer just in case Chrome freezes again
                const maxReadTime = (text.length * 100) + 3000;
                let safetyTimer = setTimeout(() => {
                    if (window.STEM_SIGHT_ACTIVE) {
                        console.warn("⚠️ Chrome voice got stuck! Forcing skip to next element.");
                    }
                    resolve(); 
                }, maxReadTime);
                
                window.stemSightUtterance.onend = () => {
                    clearTimeout(safetyTimer); 
                    console.log("✅ Finished speaking.");
                    resolve(); 
                };
                
                window.stemSightUtterance.onerror = (e) => {
                    clearTimeout(safetyTimer);
                    if (e.error !== 'canceled' && e.error !== 'interrupted') {
                        console.error("❌ Speech error:", e.error);
                    }
                    resolve(); 
                };
                
                window.speechSynthesis.speak(window.stemSightUtterance);
            }, 100); 
        });
    }

    // 5. Secure Backend Request
    async function analyzeImage(imgSrc) {
        return new Promise((resolve) => {
            chrome.runtime.sendMessage({ action: "analyzeImage", src: imgSrc }, (response) => {
                if (chrome.runtime.lastError) {
                    resolve("Sorry, the extension disconnected from the server.");
                } else {
                    resolve(response.text);
                }
            });
        });
    }

    // 6. Main Reading Loop
    async function startReading() {
        await speakAndWait("STEM sight activated. Reading page. Press escape to stop.");

        const elements = document.querySelectorAll('h1, h2, h3, h4, p, img');
        console.log(`📄 Found ${elements.length} elements to read.`);

        for (let i = 0; i < elements.length; i++) {
            if (!window.STEM_SIGHT_ACTIVE) break; // Emergency Stop Check

            const el = elements[i];

            if (el.tagName === 'P' || el.tagName.startsWith('H')) {
                const text = el.innerText.trim();
                if (text.length > 10) { 
                    await speakAndWait(text);
                }
            } 
            else if (el.tagName === 'IMG') {
                if (el.width < 50 || el.height < 50) continue;
                
                console.log("🔎 Found an image, sending to backend:", el.src);
                await speakAndWait("Image encountered. Analyzing.");
                
                if (!window.STEM_SIGHT_ACTIVE) break;
                
                const explanation = await analyzeImage(el.src);
                
                if (!window.STEM_SIGHT_ACTIVE) break;
                
                console.log("🧠 Backend response:", explanation);
                await speakAndWait(explanation);
            }
        }
        
        if (window.STEM_SIGHT_ACTIVE) {
            await speakAndWait("End of page reached.");
            window.STEM_SIGHT_ACTIVE = false;
        }
    }

    // Start!
    startReading();
})();