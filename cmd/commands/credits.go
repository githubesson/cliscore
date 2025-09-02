package commands

import (
	"flag"
	"fmt"
	"os"

	"cliscore/internal/client"
	"cliscore/internal/config"
)

type CreditsCommand struct{}

func (c *CreditsCommand) Name() string {
	return "credits"
}

func (c *CreditsCommand) Description() string {
	return "Check your remaining credits"
}

func (c *CreditsCommand) Execute(args []string) error {
	var (
		apiKey string
		quiet  bool
	)

	flagSet := flag.NewFlagSet("credits", flag.ExitOnError)
	flagSet.StringVar(&apiKey, "api-key", "", "API key for authentication (overrides env var)")
	flagSet.BoolVar(&quiet, "quiet", false, "Quiet mode (minimal output)")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	cfg := config.Load()
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	if cfg.APIKey == "" {
		fmt.Println("Error: API key is required. Set CLISCORE_API_KEY environment variable or use --api-key flag")
		os.Exit(1)
	}

	apiClient := client.New(cfg)

	response, err := apiClient.GetCredits(cfg.APIKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if !quiet {
		if response.Message != "" {
			fmt.Printf("%s\n", response.Message)
		}
		fmt.Printf("Credits remaining: %d\n", response.Credits)
	} else {
		fmt.Printf("%d\n", response.Credits)
	}

	return nil
}