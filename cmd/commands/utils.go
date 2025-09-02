package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"cliscore/internal/models"
)

// PrettyPrint prints data with special handling for machine info file trees
func PrettyPrint(data interface{}) {
	// Check if this is machine info with fileTree
	if info, ok := data.(*models.NormalizedMachineInfo); ok && len(info.FileTree) > 0 {
		// Create a copy without the fileTree for JSON printing
		infoCopy := *info
		fileTree := infoCopy.FileTree
		infoCopy.FileTree = nil
		
		// Print the machine info as JSON
		pretty, err := json.MarshalIndent(infoCopy, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", infoCopy)
			return
		}
		fmt.Println(string(pretty))
		
		// Print the formatted file tree
		fmt.Println("\n📁 File Structure:")
		fmt.Println(FormatFileTree(fileTree))
		return
	}
	
	// Regular JSON pretty print for other data
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(pretty))
}

// FormatFileTree formats a flat file list into a tree structure
func FormatFileTree(fileTree []string) string {
	if len(fileTree) == 0 {
		return "No files found"
	}

	// Build a tree structure
	tree := make(map[string]interface{})
	
	for _, filePath := range fileTree {
		parts := strings.Split(filePath, "/")
		current := tree
		
		for i, part := range parts {
			if i == len(parts)-1 {
				// This is a file
				current[part] = nil // nil indicates a file
			} else {
				// This is a directory
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
	}
	
	return formatTreeLevel(tree, "")
}

func formatTreeLevel(level map[string]interface{}, prefix string) string {
	var result strings.Builder
	keys := make([]string, 0, len(level))
	
	for k := range level {
		keys = append(keys, k)
	}
	
	// Sort keys for consistent output
	sort.Strings(keys)
	
	for i, key := range keys {
		isLast := i == len(keys)-1
		connector := "└── "
		if !isLast {
			connector = "├── "
		}
		
		result.WriteString(prefix + connector + key)
		
		if level[key] == nil {
			// It's a file
			result.WriteString("\n")
		} else {
			// It's a directory
			result.WriteString("/\n")
			newPrefix := prefix
			if !isLast {
				newPrefix += "│   "
			} else {
				newPrefix += "    "
			}
			result.WriteString(formatTreeLevel(level[key].(map[string]interface{}), newPrefix))
		}
	}
	
	return result.String()
}

// DetectOrPromptTypes detects data types from terms or prompts user for selection
func DetectOrPromptTypes(terms []string, detector Detector) []string {
	// Try to detect types based on the input
	detectedTypes := detector.DetectTypes(terms)
	
	if len(detectedTypes) > 0 {
		fmt.Printf("Detected types: %v\n", detectedTypes)
		fmt.Print("Use detected types? (Y/n): ")
		
		var response string
		fmt.Scanln(&response)
		
		if response == "" || strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			return detectedTypes
		}
	}
	
	// Interactive type selection
	return PromptForTypes()
}

// PromptForTypes prompts user to select data types interactively
func PromptForTypes() []string {
	fmt.Println("\nAvailable types:")
	fmt.Println("1. login")
	fmt.Println("2. password")
	fmt.Println("3. url")
	fmt.Println("4. domain")
	fmt.Println("5. username")
	fmt.Println("6. ip")
	fmt.Println("7. hash")
	fmt.Println("8. phone")
	fmt.Println("9. uuid")
	fmt.Println()
	fmt.Print("Select types (comma-separated numbers, e.g., '1,2,3' or 'all'): ")
	
	var input string
	fmt.Scanln(&input)
	
	input = strings.ToLower(strings.TrimSpace(input))
	
	if input == "all" {
		return []string{"login", "password", "url", "domain", "username", "ip", "hash", "phone", "uuid"}
	}
	
	typeMap := map[string]string{
		"1": "login",
		"2": "password",
		"3": "url",
		"4": "domain",
		"5": "username",
		"6": "ip",
		"7": "hash",
		"8": "phone",
		"9": "uuid",
	}
	
	var selectedTypes []string
	parts := strings.Split(input, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if typeName, exists := typeMap[part]; exists {
			selectedTypes = append(selectedTypes, typeName)
		}
	}
	
	if len(selectedTypes) == 0 {
		fmt.Println("No valid types selected, using defaults: login, password, url")
		return []string{"login", "password", "url"}
	}
	
	return selectedTypes
}

// Detector interface for type detection
type Detector interface {
	DetectTypes(terms []string) []string
}