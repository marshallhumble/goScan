package RegexProcessing

import (
	"goScan/ReadFunctions"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`\b[A-Za-z0-9](?:[A-Za-z0-9._%+-]*[A-Za-z0-9])?@[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?\.[A-Za-z]{2,6}\b`)
	dobRegex   = regexp.MustCompile(`\b(0[1-9]|1[0-2])[-/](0[1-9]|[12][0-9]|3[01])[-/](\d{2}|\d{4})\b`)
	ssnRegex   = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b|\b\d{3}\s\d{2}\s\d{4}\b|\b\d{9}\b`)
	phoneRegex = regexp.MustCompile(`\b\(?[0-9]{3}\)?[-.\s]?[0-9]{3}[-.\s]?[0-9]{4}\b|\b[0-9]{10,11}\b`)
	nameRegex  = regexp.MustCompile(`\b[A-Z][a-zA-Z'-]{1,}(?:\s[A-Z][a-zA-Z'-]{1,})*\b`)
)

func checkStrings(line string) (ReadFunctions.PIIDetection, error) {

}
