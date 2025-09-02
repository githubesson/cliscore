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

type SearchCommand struct{}

func (c *SearchCommand) Name() string {
	return "search"
}

func (c *SearchCommand) Description() string {
	return "Search for terms across different data types"
}

func (c *SearchCommand) Execute(args []string) error {
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

	flagSet := flag.NewFlagSet("search", flag.ExitOnError)
	flagSet.StringVar(&source, "source", "xkeyscore", "Source to search from")
	flagSet.BoolVar(&wildcard, "wildcard", false, "Enable wildcard search")
	flagSet.StringVar(&apiKey, "api-key", "", "API key for authentication (overrides env var)")
	flagSet.BoolVar(&saveResults, "save", false, "Save results to file")
	flagSet.BoolVar(&noSave, "no-save", false, "Don't save results to file")
	flagSet.StringVar(&resultsDir, "results-dir", "", "Results directory (overrides config)")
	flagSet.BoolVar(&showSpinner, "spinner", true, "Show loading spinner")
	flagSet.BoolVar(&quiet, "quiet", false, "Quiet mode (no spinner)")
	flagSet.StringVar(&operator, "operator", "", "Search operator (AND, LOGS)")

	if err := flagSet.Parse(args); err != nil {
		return err
	}
	
	terms = flagSet.Args()

	if len(terms) < 1 {
		fmt.Println("Usage: cliscore search [options] <terms...>")
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

	mappedTypes := make([]string, len(types))
	for i, t := range types {
		if t == "login" {
			mappedTypes[i] = "email"
		} else {
			mappedTypes[i] = t
		}
	}

	req := &models.SearchRequest{
		Terms:    terms,
		Types:    mappedTypes,
		Wildcard: wildcard,
		Source:   source,
	}

	if operator != "" {
		req.Operator = &operator
	}

	var spin *spinner.Spinner
	if showSpinner && !quiet {
		searchMsg := fmt.Sprintf("Searching for %s in %s...", strings.Join(terms, ", "), strings.Join(types, ", "))
		spin = config.CreateSpinner(searchMsg)
		if spin != nil {
			spin.Start()
		}
	}

	response, err := apiClient.Search(req, cfg.APIKey)
	
	if spin != nil {
		spin.Stop()
	}
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	resultCount := len(response.Results)
	fmt.Printf("Found %d results\n", resultCount)

	if !quiet && !cfg.SaveResults {
		fmt.Printf("Search Results:\n")
		PrettyPrint(response.Results)
	}

	if cfg.SaveResults {
		if err := config.SaveResults(response.Results, "search", terms, types); err != nil {
			if !quiet {
				fmt.Printf("Warning: Failed to save results: %v\n", err)
			}
		} else {
			if !quiet {
				fmt.Printf("Full response saved to: %s\n", config.GetResultsFilePath("search", terms, types))
			}
		}
	}

	return nil
}