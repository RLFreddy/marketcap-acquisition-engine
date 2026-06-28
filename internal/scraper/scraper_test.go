package scraper

import "testing"

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{`"quoted"`, "quoted"},
		{"&amp;", "&"},
		{"&lt;tag&gt;", "<tag>"},
		{"  &quot;spaced&quot;  ", "spaced"},
	}
	for _, tt := range tests {
		result := sanitize(tt.input)
		if result != tt.expected {
			t.Errorf("sanitize(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
