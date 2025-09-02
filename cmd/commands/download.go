package commands

import (
	"flag"
	"fmt"
	"os"

	"cliscore/internal/client"
	"cliscore/internal/config"
	"cliscore/internal/spinner"
)

type DownloadCommand struct{}

func (c *DownloadCommand) Name() string {
	return "download"
}

func (c *DownloadCommand) Description() string {
	return "Download log files"
}

func (c *DownloadCommand) Execute(args []string) error {
	var (
		uuid       string
		filePath   string
		outputPath string
		apiKey     string
		quiet      bool
	)

	flag.StringVar(&uuid, "uuid", "", "UUID of the log file")
	flag.StringVar(&filePath, "file", "", "Specific file to extract from the archive")
	flag.StringVar(&outputPath, "output", "", "Output file path")
	flag.StringVar(&apiKey, "api-key", "", "API key for authentication (overrides env var)")
	flag.BoolVar(&quiet, "quiet", false, "Quiet mode (minimal output)")

	flag.Parse()

	if len(args) < 1 {
		fmt.Println("Usage: cliscore download [options] <uuid>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get UUID from flag or argument
	if uuid == "" {
		uuid = args[0]
	}

	if uuid == "" {
		fmt.Println("UUID is required")
		os.Exit(1)
	}

	cfg := config.Load()
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	apiClient := client.New(cfg)

	// Build description for spinner
	description := fmt.Sprintf("Downloading file for UUID: %s", uuid)
	if filePath != "" {
		description = fmt.Sprintf("Downloading %s from UUID: %s", filePath, uuid)
	}

	// Start spinner if enabled
	var spin *spinner.Spinner
	if !quiet {
		spin = config.CreateSpinner(description)
		if spin != nil {
			spin.Start()
		}
	}

	err := apiClient.DownloadFile(uuid, filePath, cfg.APIKey, outputPath)
	
	// Stop spinner
	if spin != nil {
		spin.Stop()
	}
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return nil
}