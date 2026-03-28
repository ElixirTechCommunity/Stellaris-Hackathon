from graph_engine import GraphExtractor

def test_pipeline():
    print("--- Spectra Engine Diagnostics ---")
    
    # 1. Initialize the AI (this will trigger the HF download if it hasn't already)
    extractor = GraphExtractor()
    
    # 2. Define your test image
    # IMPORTANT: Place a picture of a graph in your folder and rename it to 'sample.png'
    # Or change this variable to match your image's actual name!
    image_to_test = "Sample2.png" 
    
    print(f"\nFeeding '{image_to_test}' into the vision engine...")
    
    # 3. Extract the data
    try:
        raw_json = extractor.extract(image_to_test)
        
        print("\n" + "="*40)
        print("🧠 RAW JSON OUTPUT:")
        print("="*40)
        
        # It's already formatted as a nice JSON string from your engine
        print(raw_json) 
        
        print("="*40)
        
    except Exception as e:
        print(f"\n❌ Error during extraction: {e}")
        print("Did you make sure the image file exists in the folder?")

if __name__ == "__main__":
    test_pipeline()