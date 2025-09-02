package commands

import (
	"flag"
	"fmt"
	"os"

	"cliscore/internal/client"
	"cliscore/internal/config"
	"cliscore/internal/spinner"
)

type MachineInfoCommand struct{}

func (c *MachineInfoCommand) Name() string {
	return "machineinfo"
}

func (c *MachineInfoCommand) Description() string {
	return "Get machine information from log files"
}

func (c *MachineInfoCommand) Execute(args []string) error {
	var (
		uuid      string
		apiKey    string
		quiet     bool
	)

	flag.StringVar(&uuid, "uuid", "", "UUID of the log file")
	flag.StringVar(&apiKey, "api-key", "", "API key for authentication (overrides env var)")
	flag.BoolVar(&quiet, "quiet", false, "Quiet mode (minimal output)")

	flag.Parse()

	if len(args) < 1 {
		fmt.Println("Usage: cliscore machineinfo [options] <uuid>")
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

	// Start spinner if enabled
	var spin *spinner.Spinner
	if !quiet {
		spin = config.CreateSpinner(fmt.Sprintf("Retrieving machine info for UUID: %s", uuid))
		if spin != nil {
			spin.Start()
		}
	}

	response, err := apiClient.GetMachineInfo(uuid, cfg.APIKey)
	
	// Stop spinner
	if spin != nil {
		spin.Stop()
	}
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if response.Error != "" {
		fmt.Printf("Error: %s\n", response.Error)
		os.Exit(1)
	}

	if !quiet {
		fmt.Println("Machine Information:")
		PrettyPrint(response.Data)
	}

	// Save results if enabled
	if err := config.SaveResults(response.Data, "machineinfo", []string{uuid}, []string{"log"}); err != nil {
		if !quiet {
			fmt.Printf("Warning: Failed to save results: %v\n", err)
		}
	}

	return nil
}