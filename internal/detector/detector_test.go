package detector

import (
	"testing"
)

func TestDetector_UUID(t *testing.T) {
	d := New()
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"ce1869f2-b922-456b-882c-58aa4ad5f266", true},
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"not-a-uuid", false},
		{"ce1869f2-b922-456b-882c", false}, // too short
		{"ce1869f2-b922-456b-882c-58aa4ad5f266-extra", false}, // too long
		{"gggggggg-gggg-gggg-gggg-gggggggggggg", false}, // invalid hex
	}
	
	for _, test := range tests {
		result := d.isUUID(test.input)
		if result != test.expected {
			t.Errorf("isUUID(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestDetector_Email(t *testing.T) {
	d := New()
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"admin@example.com", true},
		{"user.name@domain.org", true},
		{"test+tag@gmail.com", true},
		{"@example.com", false},
		{"user@", false},
		{"user@.com", false},
		{"user.com@", false},
		{"not-an-email", false},
	}
	
	for _, test := range tests {
		result := d.isEmail(test.input)
		if result != test.expected {
			t.Errorf("isEmail(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestDetector_URL(t *testing.T) {
	d := New()
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://example.com", true},
		{"http://example.org", true},
		{"ftp://files.domain.net", true},
		{"ws://websocket.example.io", true},
		{"custom-protocol://server.com", true},
		{"example.com", false},
		{"sub.domain.net", false},
		{"not-a-url", false},
		{"just-text", false},
		{"https://", false},  // Missing domain
		{"://example.com", false}, // Missing protocol
	}
	
	for _, test := range tests {
		result := d.isURL(test.input)
		if result != test.expected {
			t.Errorf("isURL(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}


func TestDetector_Domain(t *testing.T) {
	d := New()
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"example.com", true},
		{"sub.domain.org", true},
		{"test.net", true},
		{"admin@example.com", false}, // Email
		{"https://example.com", false}, // URL
		{"not-a-domain", false},
		{".com", false},
		{"domain.", false},
	}
	
	for _, test := range tests {
		result := d.isDomain(test.input)
		if result != test.expected {
			t.Errorf("isDomain(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestDetector_DetectTypes(t *testing.T) {
	d := New()
	
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"ce1869f2-b922-456b-882c-58aa4ad5f266"}, []string{"uuid"}},
		{[]string{"admin@example.com"}, []string{"email"}},
		{[]string{"https://example.com"}, []string{"url"}},
		{[]string{"example.com"}, []string{"domain"}},
		{[]string{"ftp://files.domain.net"}, []string{"url"}},
		{[]string{"admin@example.com", "test.domain.com"}, []string{"email", "domain"}},
		{[]string{"unknown"}, []string{}},
	}
	
	for _, test := range tests {
		result := d.DetectTypes(test.input)
		if !stringSlicesEqual(result, test.expected) {
			t.Errorf("DetectTypes(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestDetector_DetectTypes_Helper(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"ce1869f2-b922-456b-882c-58aa4ad5f266"}, []string{"uuid"}},
		{[]string{"admin@example.com"}, []string{"email"}},
		{[]string{"unknown"}, []string{}},
	}
	
	for _, test := range tests {
		result := DetectTypes(test.input)
		if !stringSlicesEqual(result, test.expected) {
			t.Errorf("DetectTypes(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	
	aMap := make(map[string]bool)
	bMap := make(map[string]bool)
	
	for _, s := range a {
		aMap[s] = true
	}
	for _, s := range b {
		bMap[s] = true
	}
	
	if len(aMap) != len(bMap) {
		return false
	}
	
	for k := range aMap {
		if !bMap[k] {
			return false
		}
	}
	
	return true
}