package commands

import (
	"fmt"
	"os"

	"cliscore/internal/config"
)

type ConfigCommand struct{}

func (c *ConfigCommand) Name() string {
	return "config"
}

func (c *ConfigCommand) Description() string {
	return "Show current configuration"
}

func (c *ConfigCommand) Execute(args []string) error {
	cfg := config.Load()
	
	fmt.Println("Current configuration:")
	fmt.Printf("Base URL: %s\n", cfg.BaseURL)
	
	if cfg.APIKey != "" {
		fmt.Println("API key: ********")
	} else {
		fmt.Println("API key: Not set")
	}

	fmt.Printf("Save results: %v\n", cfg.SaveResults)
	if cfg.SaveResults {
		fmt.Printf("Results directory: %s\n", cfg.ResultsDir)
	}
	
	fmt.Printf("Spinner style: %s\n", cfg.SpinnerStyle)

	if config.ConfigFileExists() {
		homeDir, _ := os.UserHomeDir()
		configPath := fmt.Sprintf("%s/.keyscore-cli/config.json", homeDir)
		fmt.Printf("Config file: %s\n", configPath)
	} else {
		fmt.Println("Config file: Not found")
	}

	return nil
}