package utilityFunctions

import (
	"io"
	"log"
	"os"
	"strings"
)

// SafeCloseFile  safe close of file writes
func SafeCloseFile(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("warning: failed to close: %v", err)
	}
}

// ReadFileEnvs pull the sensitive data details from the .ENV file that we are using for Docker init
func ReadFileEnvs(fileName string) (projectId string, err error) {

	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	text := string(data)

	projectId = getVariable(text, "PROJECT_ID")

	return projectId, nil
}

// getVariable get the variables from the ENV file, right now we are assuming they look like this:
func getVariable(text, key string) string {

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if strings.Contains(line, key) {
			// Split the line into key-value pairs
			parts := strings.Split(line, "=")

			// Get the value of the variable
			return parts[1]
		}

	}
	return ""
}

func SafeClose(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("Error closing resource: %v", err)
		}
	}
}
