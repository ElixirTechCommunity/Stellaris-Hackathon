package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aurva/compliance-engine/shared/pii"
)

func main() {
	// -------- Step 1: Generate large data --------
	fmt.Println("Generating large dataset...")

	base := "User Aadhaar: 1234 5678 9012 PAN: ABCDE1234F Phone: 9876543210\n"

	var builder strings.Builder

	targetSizeMB := 100 // 🔥 change this (100MB, 500MB, etc.)
	targetBytes := targetSizeMB * 1024 * 1024

	for builder.Len() < targetBytes {
		builder.WriteString(base)
	}

	data := builder.String()

	fmt.Printf("Generated data size: %.2f MB\n", float64(len(data))/(1024*1024))

	// -------- Step 2: Init classifier --------
	classifier := pii.NewClassifier()

	// -------- Step 3: Measure time --------
	fmt.Println("Starting scan...")
	start := time.Now()

	results := classifier.Scan(data)

	elapsed := time.Since(start)

	// -------- Step 4: Output --------
	fmt.Println("Scan complete")
	fmt.Printf("Time taken: %s\n", elapsed)
	fmt.Printf("Detections found: %d\n", len(results))

	// Throughput
	mbProcessed := float64(len(data)) / (1024 * 1024)
	seconds := elapsed.Seconds()
	fmt.Printf("Throughput: %.2f MB/s\n", mbProcessed/seconds)

	// Optional: write to file (simulate real pipeline)
	_ = os.WriteFile("large_input.txt", []byte(data), 0644)
}
