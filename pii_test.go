package main

import (
	"testing"
)

func TestDOB(t *testing.T) {
	testCases := []struct {
		text     string
		expected bool
	}{

		{"01/02/1960", true},
		{"02/10/2000", true},
		{"11-03-90", true},
		{"123ABC", false},
		{"some text", false},

		{"12/31/99", true},    // 2-digit year with /
		{"03-15-2023", true},  // 4-digit year with -
		{"13/01/2000", false}, // Invalid month
		{"01/32/2000", false}, // Invalid day
		{"1/2/2000", false},   // Single digit month/day (not matched by this regex)
		{"01.02.2000", false}, // Different separator

	}
	for _, tt := range testCases {
		result := dobRegex.MatchString(tt.text)
		if result != tt.expected {
			t.Errorf("dobRegex.MatchString(%q) = %v; want %v", tt.text, result, tt.expected)
		}
	}
}

func TestEmail(t *testing.T) {
	testCases := []struct {
		text     string
		expected bool
		comment  string
	}{
		{"user@example.com", true, "Basic email"},
		{"john.doe@company.org", true, "Dot in local part"},
		{"jane+newsletter@site.co.uk", true, "Plus sign and subdomain"},
		{"test_user@domain-name.info", true, "Underscore and hyphen"},
		{"123@numbers.net", true, "Numeric local part"},
		{"user123@test123.com", true, "Mixed alphanumeric"},
		{"first.last+tag@example-domain.travel", true, "Complex valid email"},
		{"a@b.co", true, "Minimal valid email"},
		{"very.long.email.address@very-long-domain-name.museum", true, "Long email"},

		// Realistic invalid cases - should NOT match
		{"plaintext", false, "No @ symbol"},
		{"@domain.com", false, "Missing local part"},
		{"user@", false, "Missing domain"},
		{"user@domain", false, "Missing TLD"},
		{"user@.com", false, "Missing domain name"},
		{"user@domain.c", false, "TLD too short"},
		{"user@@domain.com", false, "Double @"},
		{"user@domain@com", false, "Multiple @"},
		{"", false, "Empty string"},
		{"123ABC", false, "Random text"},

		// Realistic edge cases that SHOULD be caught by word boundaries
		{"Contact me at john@example.com today", true, "Email in sentence"},
		{"Email: admin@company.org", true, "Email with label"},
		{"(support@help.net)", true, "Email in parentheses"},
		{"Visit http://user@domain.com/path", true, "Email in URL context"},
	}
	for _, tt := range testCases {
		result := emailRegex.MatchString(tt.text)
		if result != tt.expected {
			t.Errorf("emailRegex.MatchString(%q) = %v; want %v", tt.text, result, tt.expected)
		}
	}
}

func TestSSN(t *testing.T) {
	testCases := []struct {
		text     string
		expected bool
		comment  string
	}{
		// Valid SSN formats - should match
		{"123-45-6789", true, "Standard SSN format"},
		{"987-65-4321", true, "Another standard SSN"},
		{"555-12-3456", true, "Valid SSN pattern"},
		{"123 45 6789", true, "Space-separated SSN"},
		{"987 65 4321", true, "Another space-separated"},
		{"123456789", true, "No separators"},
		{"987654321", true, "Nine digits no separators"},

		// Valid in context - should match
		{"My SSN is 123-45-6789", true, "SSN in sentence"},
		{"SSN: 987-65-4321", true, "SSN with label"},
		{"ID 123 45 6789 on file", true, "SSN with spaces in context"},
		{"Reference #123456789", true, "SSN without separators in context"},

		// THESE MATCH THE PATTERN (handle validation separately)
		{"000-00-0000", true, "All zeros - pattern match, validate separately"},
		{"666-12-3456", true, "Starts with 666 - pattern match, validate separately"},
		{"900-12-3456", true, "Starts with 9xx - pattern match, validate separately"},
		{"123-00-4567", true, "Middle group zeros - pattern match, validate separately"},
		{"123-45-0000", true, "Last group zeros - pattern match, validate separately"},

		// Wrong formats - should NOT match
		{"12-345-6789", false, "Wrong grouping 2-3-4"},
		{"1234-5-6789", false, "Wrong grouping 4-1-4"},
		{"123-456-789", false, "Wrong grouping 3-3-3"},
		{"1234-56-789", false, "Wrong grouping 4-2-3"},
		{"12-34-56789", false, "Wrong grouping 2-2-5"},

		// Too short/long - should NOT match
		{"12-34-567", false, "Too short"},
		{"1234-56-7890", false, "Too long"},
		{"12345678", false, "8 digits only"},
		{"1234567890", false, "10 digits"},

		// Non-numeric - should NOT match
		{"123-45-67AB", false, "Letters in SSN"},
		{"ABC-45-6789", false, "Letters at start"},
		{"123-AB-6789", false, "Letters in middle"},

		// Wrong separators - should NOT match
		{"123.45.6789", false, "Dot separators"},
		{"123/45/6789", false, "Slash separators"},
		{"123_45_6789", false, "Underscore separators"},
		{"123:45:6789", false, "Colon separators"},

		// Mixed separators - should NOT match (regex now handles this correctly)
		{"123-45 6789", false, "Mixed dash and space"},
		{"123 45-6789", false, "Mixed space and dash"},

		// Partial matches - should NOT match
		{"123-45-", false, "Incomplete SSN"},
		{"-45-6789", false, "Missing first group"},
		{"123--6789", false, "Missing middle group"},

		// Random numbers - should NOT match
		{"123", false, "Just 3 digits"},
		{"12345", false, "5 digits"},
		{"1234567", false, "7 digits"},
		{"123ABC", false, "Mixed alphanumeric"},
		{"", false, "Empty string"},
		{"some text", false, "Random text"},

		// Phone numbers that might be confused - should NOT match
		{"(123) 456-7890", false, "Phone number format"},
		{"123-456-7890", false, "10-digit phone (3-3-4)"},

		// Edge cases with context
		{"Call 123-456-7890 for info", false, "Phone number in context"},
		{"Born 12-34-5678 (wrong date format)", false, "Date-like pattern"},

		// THIS ONE SHOULD NOW CORRECTLY FAIL:
		{"File #123-45-6789A contains data", false, "SSN with suffix letter"},
	}

	for _, tt := range testCases {
		result := ssnRegex.MatchString(tt.text)
		if result != tt.expected {
			t.Errorf("ssnRegex.MatchString(%q) = %v; want %v (%s)",
				tt.text, result, tt.expected, tt.comment)
		}
	}
}
func TestPhoneNumber(t *testing.T) {
	testCases := []struct {
		text     string
		expected bool
		comment  string
	}{
		// Valid US phone number formats - should match
		{"(555) 123-4567", true, "Standard (XXX) XXX-XXXX format"},
		{"555-123-4567", true, "XXX-XXX-XXXX format"},
		{"555.123.4567", true, "XXX.XXX.XXXX format"},
		{"555 123 4567", true, "XXX XXX XXXX format"},
		{"5551234567", true, "XXXXXXXXXX format (10 digits)"},

		// With country code
		{"1-555-123-4567", true, "1-XXX-XXX-XXXX format"},
		{"1 555 123 4567", true, "1 XXX XXX XXXX format"},
		{"1.555.123.4567", true, "1.XXX.XXX.XXXX format"},
		{"1 (555) 123-4567", true, "1 (XXX) XXX-XXXX format"},
		{"15551234567", true, "1XXXXXXXXXX format (11 digits)"},
		{"+1 555 123 4567", true, "International format with +1"},
		{"+1-555-123-4567", true, "International format +1-XXX-XXX-XXXX"},
		{"+15551234567", true, "International format no spaces"},

		// In context - should match
		{"Call me at (555) 123-4567", true, "Phone in sentence"},
		{"Phone: 555-123-4567", true, "Phone with label"},
		{"Contact 1-800-555-1234 for help", true, "Toll-free number"},
		{"Mobile: +1 (555) 123-4567", true, "International mobile"},

		// Extensions (optional - decide based on your needs)
		{"555-123-4567 ext 123", true, "Phone with extension"},
		{"(555) 123-4567 x456", true, "Phone with x extension"},

		// Invalid formats - should NOT match
		{"555-1234", false, "Too short (7 digits)"},
		{"555-12-34567", false, "Wrong grouping"},
		{"55-123-4567", false, "Wrong area code length"},
		{"555-123-456", false, "Too short last group"},
		{"555-123-45678", false, "Too long last group"},

		// Non-numeric characters
		{"555-ABC-4567", false, "Letters in phone number"},
		{"(555) 12A-4567", false, "Letter in exchange"},
		{"555-123-456X", false, "Letter in last group"},

		// Wrong separators/format
		{"555_123_4567", false, "Underscore separators"},
		{"555/123/4567", false, "Slash separators"},
		{"555:123:4567", false, "Colon separators"},
		{"555-(123)-4567", false, "Parentheses in wrong place"},

		// Too many/few digits
		{"555-123-45678", false, "Too many digits in last group"},
		{"1234567890123", false, "Too many total digits"},
		{"12345", false, "Way too short"},
		{"555123", false, "6 digits only"},

		// Edge cases
		{"", false, "Empty string"},
		{"some text", false, "Random text"},
		{"123ABC", false, "Mixed alphanumeric"},

		// Ambiguous with other numbers (these should NOT match phone patterns)
		{"123-45-6789", false, "SSN format"},
		{"01/02/1960", false, "Date format"},
		{"Account #555-123-4567", true, "Phone number as account (still matches)"},

		// Partial matches - should NOT match
		{"Call (555) 123-", false, "Incomplete phone"},
		{"555-", false, "Just area code"},
		{"(555)", false, "Just area code in parens"},

		// Multiple formats in text
		{"Call (555) 123-4567 or 1-800-555-9999", true, "Multiple phones (matches first)"},
	}

	for _, tt := range testCases {
		result := phoneRegex.MatchString(tt.text)
		if result != tt.expected {
			t.Errorf("phoneRegex.MatchString(%q) = %v; want %v (%s)",
				tt.text, result, tt.expected, tt.comment)
		}
	}
}
func TestNames(t *testing.T) {
	testCases := []struct {
		text     string
		expected bool
		comment  string
	}{
		// Valid names - should match
		{"John", true, "Single first name"},
		{"Jane", true, "Single first name"},
		{"John Smith", true, "First and last name"},
		{"Mary Johnson", true, "First and last name"},
		{"Robert Williams", true, "First and last name"},
		{"Sarah Davis", true, "First and last name"},

		// Names with common variations - should match
		{"John-Paul", true, "Hyphenated first name"},
		{"Mary-Kate", true, "Hyphenated first name"},
		{"O'Connor", true, "Irish name with apostrophe"},
		{"D'Angelo", true, "Italian name with apostrophe"},
		{"McBride", true, "Scottish/Irish surname"},

		// Full names - should match
		{"John Michael Smith", true, "First middle last"},
		{"Mary Elizabeth Johnson", true, "First middle last"},

		// Names in context - should match
		{"Hello John", true, "Name in greeting"},
		{"Contact Sarah Davis", true, "Name in instruction"},
		{"Patient: Mary Johnson", true, "Name with label"},

		// Clear non-names - should NOT match
		{"john", false, "All lowercase (not typical name format)"},
		{"john smith", false, "All lowercase"},
		{"123", false, "Numbers only"},
		{"John123", false, "Name with numbers"},
		{"john_doe", false, "Underscore format"},
		{"", false, "Empty string"},

		// THESE MATCH THE PATTERN (regex can't distinguish context):
		// Regex sees "capitalized word" pattern - it can't know if it's a name or not
		{"JOHN", true, "All caps - matches pattern (validate separately)"},
		{"Jo", true, "Short word - matches pattern (validate separately)"},
		{"The", true, "Article - matches pattern (validate separately)"},
		{"And", true, "Conjunction - matches pattern (validate separately)"},
		{"This", true, "Pronoun - matches pattern (validate separately)"},
		{"From", true, "Preposition - matches pattern (validate separately)"},
		{"With", true, "Preposition - matches pattern (validate separately)"},
		{"API", true, "Acronym - matches pattern (validate separately)"},
		{"HTTP", true, "Acronym - matches pattern (validate separately)"},
		{"JSON", true, "Acronym - matches pattern (validate separately)"},
		{"An", true, "Article - matches pattern (validate separately)"},

		// Ambiguous cases that legitimately could be names
		{"Apple", true, "Could be surname (Johnny Apple)"},
		{"Ford", true, "Could be surname (Harrison Ford)"},
		{"Rose", true, "Could be first name (Rose Wilson)"},
		{"Will", true, "Could be first name (Will Smith)"},
		{"May", true, "Could be first name (May Johnson)"},
	}

	for _, tt := range testCases {
		result := nameRegex.MatchString(tt.text)
		if result != tt.expected {
			t.Errorf("nameRegex.MatchString(%q) = %v; want %v (%s)",
				tt.text, result, tt.expected, tt.comment)
		}
	}
}
