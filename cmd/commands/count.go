package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cliscore/internal/client"
	"cliscore/internal/config"
	"cliscore/internal/detector"
	"cliscore/internal/models"
	"cliscore/internal/spinner"
)

type CountCommand struct{}

func (c *CountCommand) Name() string {
	return "count"
}

func (c *CountCommand) Description() string {
	return "Count occurrences of terms"
}

func (c *CountCommand) Execute(args []string) error {
	var (
		terms        []string
		types        []string
		wildcard     bool
		source       string
		apiKey       string
		saveResults  bool
		noSave       bool
		resultsDir   string
		showSpinner  bool
		quiet        bool
		operator     string
	)

	flagSet := flag.NewFlagSet("count", flag.ExitOnError)
	flagSet.StringVar(&source, "source", "xkeyscore", "Source to count from")
	flagSet.BoolVar(&wildcard, "wildcard", false, "Enable wildcard search")
	flagSet.StringVar(&apiKey, "api-key", "", "API key for authentication (overrides env var)")
	flagSet.BoolVar(&saveResults, "save", false, "Save results to file")
	flagSet.BoolVar(&noSave, "no-save", false, "Don't save results to file")
	flagSet.StringVar(&resultsDir, "results-dir", "", "Results directory (overrides config)")
	flagSet.BoolVar(&showSpinner, "spinner", true, "Show loading spinner")
	flagSet.BoolVar(&quiet, "quiet", false, "Quiet mode (no spinner)")
	flagSet.StringVar(&operator, "operator", "", "Search operator (AND, LOGS)")

	flagSet.Parse(args)

	terms = flagSet.Args()

	if len(terms) < 1 {
		fmt.Println("Usage: cliscore count [options] <terms...>")
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	if len(types) == 0 {
		types = DetectOrPromptTypes(terms, detector.New())
	}

	cfg := config.Load()
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	
	// Override save settings with command line flags
	if saveResults {
		cfg.SaveResults = true
	}
	if noSave {
		cfg.SaveResults = false
	}
	if resultsDir != "" {
		cfg.ResultsDir = resultsDir
	}

	apiClient := client.New(cfg)

	// Map types like frontend does (login -> email)
	mappedTypes := make([]string, len(types))
	for i, t := range types {
		if t == "login" {
			mappedTypes[i] = "email"
		} else {
			mappedTypes[i] = t
		}
	}

	req := &models.CountRequest{
		Terms:    terms,
		Types:    mappedTypes,
		Wildcard: wildcard,
		Source:   source,
	}

	// Set operator if provided
	if operator != "" {
		req.Operator = &operator
	}

	// Start spinner if enabled
	var spin *spinner.Spinner
	if showSpinner && !quiet {
		countMsg := fmt.Sprintf("Counting %s in %s...", strings.Join(terms, ", "), strings.Join(types, ", "))
		spin = config.CreateSpinner(countMsg)
		if spin != nil {
			spin.Start()
		}
	}

	response, err := apiClient.Count(req, cfg.APIKey)
	
	// Stop spinner
	if spin != nil {
		spin.Stop()
	}
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if !quiet {
		fmt.Printf("Count Results: %d\n", response.Count)
	} else {
		// In quiet mode, just print the count
		fmt.Printf("%d\n", response.Count)
	}

	// Save results if enabled
	countResult := map[string]interface{}{
		"count": response.Count,
	}
	if err := config.SaveResults(countResult, "count", terms, types); err != nil {
		if !quiet {
			fmt.Printf("Warning: Failed to save results: %v\n", err)
		}
	}

	return nil
}