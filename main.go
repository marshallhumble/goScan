package main

import (
	"flag"
	"fmt"

	"goScan/ReadFunctions"
)

func main() {

	// Define command-line flags
	fsFile := flag.String("file", "", "Scan a single file for sensitive data, use -file <filename>")
	fsScan := flag.Bool("scan", false, "Enable scanning on the file system, requires -path")
	fsPath := flag.String("path", "", "Path to scan for files")
	help := flag.Bool("help", false, "Show help")
	writeJSON := flag.Bool("writeJSON", false, "Write results to a JSON file (not implemented yet)")
	flag.Parse()

	if *fsFile == "" && !*fsScan && *fsPath == "" || *help {
		flag.Usage()
		return
	}

	if *fsFile != "" {
		var fileAttr = ReadFunctions.FileAttributes{}
		fileAttr, err := ReadFunctions.DetectFileType(*fsFile)
		if err != nil {
			return
		}

		fileAttr, err = ReadFunctions.ReadFile(fileAttr)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", *fsFile, err)
			return
		}

		if fileAttr.TotalPHICount > 0 || fileAttr.TotalPIICount > 0 {
			showDetections(fileAttr)
		}

	}

	if *fsScan && *fsPath != "" {

	} else if *fsScan && *fsPath != "" {
		fmt.Println("You must specify a path to scan for files")
		flag.Usage()
		return
	}

	if *writeJSON {
		fmt.Println("Writing results to JSON is not implemented yet")
		// Implement JSON writing logic here
	}

}

func showDetections(fileAttr ReadFunctions.FileAttributes) {
	if fileAttr.TotalPIICount > 0 {
		fmt.Printf("Total PII Count: %d\n", fileAttr.TotalPIICount)
		for _, pii := range fileAttr.PIIDetections {
			fmt.Printf("PII Detected: Type: %s, Value: %s, Redacted: %s, Confidence: %.2f\n",
				pii.Type, pii.Value, pii.RedactedValue, pii.Confidence)
		}
	} else {
		fmt.Println("No PII detected in the file.")
	}

	if fileAttr.TotalPHICount > 0 {
		fmt.Printf("Total PHI Count: %d\n", fileAttr.TotalPHICount)
		for _, phi := range fileAttr.PHIDetections {
			fmt.Printf("PHI Detected: Type: %s, Value: %s, Redacted: %s, Confidence: %.2f\n",
				phi.Type, phi.Value, phi.RedactedValue, phi.Confidence)
		}
	} else {
		fmt.Println("No PHI detected in the file.")
	}
}
