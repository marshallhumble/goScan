package ReadFunctions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"goScan/utilityFunctions"
)

type FileAttributes struct {
	// Basic file metadata
	FilePath     string    `json:"file_path"`
	FileType     string    `json:"file_type"`
	FileSize     int64     `json:"file_size"`
	CreatedDate  time.Time `json:"created_date"`
	ModifiedDate time.Time `json:"modified_date"`

	// Processing metadata
	ProcessedAt    time.Time `json:"processed_at"`
	ProcessingTime int64     `json:"processing_time_ms"`
	ProcessorUsed  string    `json:"processor_used"` // "go-regex", "python-ml", "ocr"

	// Detection results
	PIIDetections []PIIDetection `json:"pii_detections"`
	PHIDetections []PIIDetection `json:"phi_detections"`

	// Summary statistics
	TotalPIICount   int     `json:"total_pii_count"`
	TotalPHICount   int     `json:"total_phi_count"`
	RiskScore       float64 `json:"risk_score"` // 0.0-1.0
	ConfidenceScore float64 `json:"confidence_score"`

	// Processing status
	Status   string   `json:"status"` // "success", "error", "partial"
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`

	// Content analysis
	DocumentType   string `json:"document_type,omitempty"`   // "medical", "financial", "legal"
	ContentPreview string `json:"content_preview,omitempty"` // First 200 chars for context
}

type PIIDetection struct {
	Type            string  `json:"type"`           // "ssn", "email", "phone", "name", "address"
	Value           string  `json:"value"`          // Redacted or full value based on config
	RedactedValue   string  `json:"redacted_value"` // "XXX-XX-1234"
	StartOffset     int     `json:"start_offset"`
	EndOffset       int     `json:"end_offset"`
	LineNumber      int     `json:"line_number,omitempty"`
	Confidence      float64 `json:"confidence"`       // 0.0-1.0
	Context         string  `json:"context"`          // Surrounding text for validation
	DetectionMethod string  `json:"detection_method"` // "regex", "ml", "manual"
}

func DetectFileType(filePath string) (FileAttributes, error) {

	fileAttr := FileAttributes{}

	file, err := os.Open(filePath)
	if err != nil {
		return fileAttr, err
	}
	defer utilityFunctions.SafeClose(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return fileAttr, err
	}

	fileSize := fileInfo.Size()
	if fileSize <= 0 {
		return fileAttr, fmt.Errorf("file is empty: %s", filePath)
	}

	fileAttr.FilePath = filePath
	fileAttr.FileSize = fileInfo.Size()
	fileAttr.CreatedDate = fileInfo.ModTime()
	fileAttr.ModifiedDate = fileInfo.ModTime()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fileAttr, err
	}

	if bytes.HasPrefix(buffer, []byte("%PDF")) {
		fileAttr.FileType = "pdf"
		return fileAttr, nil
	}

	//Matches a ZIP file signature, could be office documents
	//like docx, xlsx, pptx or just a ZIP file
	if bytes.HasPrefix(buffer, []byte{0x50, 0x4B, 0x03, 0x04}) {
		fileAttr.FileType = analyzeZipContent(buffer)
		return fileAttr, nil
	}

	if bytes.HasPrefix(buffer, []byte{0xFF, 0xD8, 0xFF}) {
		fileAttr.FileType = "jpeg"
		return fileAttr, nil
	}

	return fileAttr, fmt.Errorf("unsupported file type for file: %s", filePath)
}

// analyzeZipContent checks the content of a ZIP file determine if it is an office document or a generic ZIP file.
func analyzeZipContent(buffer []byte) string {
	switch {
	case bytes.Contains(buffer, []byte("word/document.xml")):
		return "docx"
	case bytes.Contains(buffer, []byte("xl/workbook.xml")):
		return "xlsx"
	case bytes.Contains(buffer, []byte("ppt/presentation.xml")):
		return "pptx"
	default:
		return "zip"
	}
}

func ReadFile(fileAttr FileAttributes) (FileAttributes, error) {
	file, err := os.Open(fileAttr.FilePath)
	if err != nil {
		return fileAttr, err
	}
	defer utilityFunctions.SafeClose(file)

	fileAttr.ProcessedAt = time.Now()

	switch fileAttr.FileType {
	case "pdf":
		// Read PDF content (placeholder)
		content, err := ReadPDFFile(file)
		if err != nil {
			return fileAttr, err
		}
		fileAttr.ContentPreview = content[:200] // First 200 chars for preview

	case "docx", "xlsx", "pptx":
		// Read Office document content (placeholder)
		content, err := ReadOfficeFile(file)
		if err != nil {
			return fileAttr, err
		}
		fileAttr.ContentPreview = content[:200] // First 200 chars for preview

	case "txt":
		if fileAttr.FileSize > 10*1024*1024 { // 10 MB threshold for large files
			lines, err := readLargeFile(fileAttr.FilePath)
			if err != nil {
				return fileAttr, err
			} else {
				lines, err := readInMemory(fileAttr.FilePath)
			}
		}

		return fileAttr, nil

	case "json":
		content, err := ReadJSONFile(file)
		if err != nil {
			return fileAttr, err
		}
		fileAttr.ContentPreview = fmt.Sprintf("%v", content)[:200] // First 200 chars for preview

	case "csv":
		content, err := ReadCSVFile(file)
		if err != nil {
			return fileAttr, err
		}
		fileAttr.ContentPreview = fmt.Sprintf("%v", content)[:200] // First 200 chars for preview

	case "sql":
		content, err := ReadSQLFile(file)
		if err != nil {
			return fileAttr, err
		}
		fileAttr.ContentPreview = fmt.Sprintf("%v", content)[:200] // First 200 chars for preview

	default:
		return fileAttr, fmt.Errorf("unsupported file type: %s", fileAttr.FileType)
	}

	return fileAttr, nil
}

func ReadPDFFile(OpenFile io.Reader) (string, error) {
	// Implementation for reading a PDF file
	return "", nil // Placeholder return
}

func ReadOfficeFile(OpenFile io.Reader) (string, error) {
	// Implementation for reading an Office file (docx, xlsx, pptx)
	return "", nil // Placeholder return
}

// readLargeFile reads a large file line by line and returns its content as a slice of strings in a buffer
func readLargeFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer utilityFunctions.SafeClose(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// readInMemory reads a small file and returns its content as an array in memory.
func readInMemory(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, os.ErrInvalid // Return an error if the content is empty
	}

	lines := strings.Split(string(data), "n")

	return lines, nil
}
