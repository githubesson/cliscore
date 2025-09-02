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
	return "Search for terms across different data types (with pagination support)"
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
		page         int
		pages        string
		pageSize     int
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
	flagSet.IntVar(&page, "page", 0, "Specific page number to retrieve (1-10)")
	flagSet.StringVar(&pages, "pages", "", "Pages to retrieve (e.g., '1,2,3' or '1-5')")
	flagSet.IntVar(&pageSize, "page-size", 0, "Number of results per page (max: 10000)")

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

	// Parse pagination parameters
	var pagination *models.SearchPaginationParams
	if page > 0 || pages != "" || pageSize > 0 {
		pagination = &models.SearchPaginationParams{}
		
		if page > 0 {
			if page > 10 {
				page = 10
			}
			pagination.Page = &page
		}
		
		if pages != "" {
			pageList := parsePages(pages)
			if len(pageList) > 0 {
				pagination.Pages = pageList
			}
		}
		
		if pageSize > 0 {
			if pageSize > 10000 {
				pageSize = 10000
			}
			pagination.PageSize = &pageSize
		}
	}

	var spin *spinner.Spinner
	if showSpinner && !quiet {
		searchMsg := fmt.Sprintf("Searching for %s in %s...", strings.Join(terms, ", "), strings.Join(types, ", "))
		spin = config.CreateSpinner(searchMsg)
		if spin != nil {
			spin.Start()
		}
	}

	var response *models.SearchResponse
	var err error
	
	if pagination != nil {
		response, err = apiClient.SearchWithPagination(req, pagination, cfg.APIKey)
	} else {
		response, err = apiClient.Search(req, cfg.APIKey)
	}
	
	if spin != nil {
		spin.Stop()
	}
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Handle different response formats for paginated vs regular searches
	var resultCount int
	var resultsToSave interface{}

	if len(response.Pages) > 0 {
		// Paginated response
		resultCount = int(response.Size)
		
		if !quiet {
			fmt.Printf("🔍 Search Results (Paginated)\n")
			fmt.Printf("%s", formatPaginationInfo(response))
			
			// Display each page's results
			for pageNum, pageResults := range response.Pages {
				fmt.Printf("\n=== Page %d ===\n", pageNum)
				PrettyPrint(pageResults)
			}
		} else {
			// Quiet mode - just show total count
			fmt.Printf("%d\n", response.Size)
		}
		
		// For saving, include the full paginated response structure
		resultsToSave = response
	} else {
		// Regular response
		resultCount = len(response.Results)
		resultsToSave = response.Results
		
		if !quiet {
			fmt.Printf("Found %d results\n", resultCount)
			fmt.Printf("Search Results:\n")
			PrettyPrint(response.Results)
		} else {
			// Quiet mode - just show count
			fmt.Printf("%d\n", resultCount)
		}
	}

	if cfg.SaveResults {
		if err := config.SaveResults(resultsToSave, "search", terms, types); err != nil {
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