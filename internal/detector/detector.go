package detector

import (
	"strings"
)

type Detector struct{}

func New() *Detector {
	return &Detector{}
}

func (d *Detector) DetectTypes(terms []string) []string {
	types := make(map[string]bool)
	
	for _, term := range terms {
		// Check UUID first (highest priority)
		if d.isUUID(term) {
			types["uuid"] = true
			continue // Skip other checks for UUIDs
		}
		
		// Check email
		if d.isEmail(term) {
			types["email"] = true
		}
		
		// Check URL
		if d.isURL(term) {
			types["url"] = true
		}
		
		// Check domain
		if d.isDomain(term) {
			types["domain"] = true
		}
	}
	
	// Convert map keys to slice
	result := make([]string, 0, len(types))
	for t := range types {
		result = append(result, t)
	}
	
	return result
}

func (d *Detector) isEmail(term string) bool {
	return strings.Contains(term, "@") && strings.Contains(term, ".") && 
	       !strings.HasPrefix(term, "@") && !strings.HasSuffix(term, "@") &&
	       !strings.HasPrefix(term, ".") && !strings.HasSuffix(term, ".")
}

func (d *Detector) isURL(term string) bool {
	// Check for protocol pattern: (string)://(string)
	return strings.Contains(term, "://") && len(strings.Split(term, "://")) == 2
}


func (d *Detector) isDomain(term string) bool {
	// Must contain a dot but not be an email or URL
	return strings.Contains(term, ".") && !d.isEmail(term) && !d.isURL(term) &&
	       !strings.HasPrefix(term, ".") && !strings.HasSuffix(term, ".") &&
	       !strings.Contains(term, "://")
}

func (d *Detector) isUUID(term string) bool {
	if len(term) != 36 {
		return false
	}
	
	parts := strings.Split(term, "-")
	if len(parts) != 5 {
		return false
	}
	
	// Check the pattern: 8-4-4-4-12
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		return false
	}
	
	// Check that all parts are hexadecimal
	for _, part := range parts {
		if !d.isHex(part) {
			return false
		}
	}
	
	return true
}

func (d *Detector) isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// Helper function for direct use without detector instance
func DetectTypes(terms []string) []string {
	d := New()
	return d.DetectTypes(terms)
}