package commands

import (
	"fmt"
	"os"
	"strings"

	"cliscore/internal/client"
	"cliscore/internal/config"
)

type SetupCommand struct{}

func (c *SetupCommand) Name() string {
	return "setup"
}

func (c *SetupCommand) Description() string {
	return "Setup API key and endpoint"
}

func (c *SetupCommand) Execute(args []string) error {
	var baseURL, apiKey, resultsDir, spinnerStyle string
	var saveResults bool
	
	// Load existing config if it exists
	cfg := config.Load()
	
	fmt.Println("🔧 Setting up CliScore configuration...")
	
	// Prompt for base URL with current value as default
	fmt.Printf("Enter API base URL (default: %s): ", cfg.BaseURL)
	fmt.Scanln(&baseURL)
	if baseURL == "" {
		baseURL = cfg.BaseURL
	}

	// Prompt for API key with current value as default (but don't show it)
	if cfg.APIKey != "" {
		fmt.Print("Enter API key (press Enter to keep current): ")
		fmt.Scanln(&apiKey)
		if apiKey == "" {
			apiKey = cfg.APIKey
		}
	} else {
		fmt.Print("Enter API key: ")
		fmt.Scanln(&apiKey)
	}

	if apiKey == "" {
		fmt.Println("API key is required")
		os.Exit(1)
	}

	// Validate API key with the backend
	fmt.Print("Validating API key...")
	apiClient := client.New(cfg)
	err := apiClient.ValidateAPIKey(apiKey)
	if err != nil {
		fmt.Printf("\n❌ API key validation failed: %v\n", err)
		fmt.Println("Please check your API key and try again.")
		os.Exit(1)
	}
	fmt.Println(" ✅ Valid")

	fmt.Print("Save results to files? (y/N): ")
	var saveResponse string
	fmt.Scanln(&saveResponse)
	if saveResponse == "" {
		saveResults = cfg.SaveResults
	} else {
		saveResults = strings.ToLower(saveResponse) == "y" || strings.ToLower(saveResponse) == "yes"
	}

	if saveResults {
		fmt.Printf("Enter results directory (default: %s): ", cfg.ResultsDir)
		fmt.Scanln(&resultsDir)
		if resultsDir == "" {
			resultsDir = cfg.ResultsDir
		}
	}

	// Prompt for spinner style
	fmt.Println("\nAvailable spinner styles:")
	fmt.Println("1. default (braille dots)")
	fmt.Println("2. dots (heavy dots)")
	fmt.Println("3. arrows (rotating arrows)")
	fmt.Println("4. bounce (bouncing dots)")
	fmt.Println("5. simple (progressive dots)")
	fmt.Println("6. none (no spinner)")
	fmt.Printf("Select spinner style (1-6, default: %s): ", cfg.SpinnerStyle)
	
	var styleResponse string
	fmt.Scanln(&styleResponse)
	if styleResponse == "" {
		spinnerStyle = cfg.SpinnerStyle
	} else {
		switch styleResponse {
		case "1":
			spinnerStyle = "default"
		case "2":
			spinnerStyle = "dots"
		case "3":
			spinnerStyle = "arrows"
		case "4":
			spinnerStyle = "bounce"
		case "5":
			spinnerStyle = "simple"
		case "6":
			spinnerStyle = "none"
		default:
			fmt.Printf("Invalid choice, using default: %s\n", cfg.SpinnerStyle)
			spinnerStyle = cfg.SpinnerStyle
		}
	}

	if err := config.SaveFullWithSpinner(baseURL, apiKey, resultsDir, saveResults, spinnerStyle); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration saved successfully!")
	fmt.Printf("📍 Base URL: %s\n", baseURL)
	fmt.Println("🔑 API key: ********")
	fmt.Printf("💾 Save results: %v\n", saveResults)
	if saveResults {
		fmt.Printf("📁 Results directory: %s\n", resultsDir)
	}
	fmt.Printf("🎨 Spinner style: %s\n", spinnerStyle)

	return nil
}