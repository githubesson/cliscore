package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cliscore/internal/spinner"
)

type Config struct {
	BaseURL      string `json:"baseURL"`
	APIKey       string `json:"apiKey"`
	ResultsDir   string `json:"resultsDir"`
	SaveResults  bool   `json:"saveResults"`
	SpinnerStyle string `json:"spinnerStyle"`
}

func Load() *Config {
	cfg := &Config{
		BaseURL:      "https://api.keysco.re",
		ResultsDir:   getDefaultResultsDir(),
		SaveResults:  false,
		SpinnerStyle: "default",
	}

	if config := loadFromFile(); config != nil {
		cfg.BaseURL = config.BaseURL
		cfg.APIKey = config.APIKey
		cfg.ResultsDir = config.ResultsDir
		cfg.SaveResults = config.SaveResults
		cfg.SpinnerStyle = config.SpinnerStyle
	}

	if url := os.Getenv("CLISCORE_BASE_URL"); url != "" {
		cfg.BaseURL = url
	}

	if key := os.Getenv("CLISCORE_API_KEY"); key != "" {
		cfg.APIKey = key
	}

	if resultsDir := os.Getenv("CLISCORE_RESULTS_DIR"); resultsDir != "" {
		cfg.ResultsDir = resultsDir
	}

	if saveResults := os.Getenv("CLISCORE_SAVE_RESULTS"); saveResults != "" {
		cfg.SaveResults = saveResults == "true" || saveResults == "1"
	}

	if spinnerStyle := os.Getenv("CLISCORE_SPINNER_STYLE"); spinnerStyle != "" {
		cfg.SpinnerStyle = spinnerStyle
	}

	return cfg
}

func getDefaultResultsDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./results"
	}
	return filepath.Join(homeDir, ".keyscore-cli", "results")
}

func loadFromFile() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	configPath := filepath.Join(homeDir, ".keyscore-cli", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}

	return &cfg
}

func Save(baseURL, apiKey string) error {
	return SaveFull(baseURL, apiKey, getDefaultResultsDir(), false)
}

func SaveFull(baseURL, apiKey, resultsDir string, saveResults bool) error {
	return SaveFullWithSpinner(baseURL, apiKey, resultsDir, saveResults, "default")
}

func SaveFullWithSpinner(baseURL, apiKey, resultsDir string, saveResults bool, spinnerStyle string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".keyscore-cli")
	if err = os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	cfg := Config{
		BaseURL:      baseURL,
		APIKey:       apiKey,
		ResultsDir:   resultsDir,
		SaveResults:  saveResults,
		SpinnerStyle: spinnerStyle,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func ConfigFileExists() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configPath := filepath.Join(homeDir, ".keyscore-cli", "config.json")
	_, err = os.Stat(configPath)
	return err == nil
}

func SaveResults(data interface{}, command string, terms []string, types []string) error {
	cfg := Load()

	if !cfg.SaveResults {
		return nil
	}

	if err := os.MkdirAll(cfg.ResultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %v", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	safeTerms := makeSafeFilename(strings.Join(terms, "_"))
	safeTypes := makeSafeFilename(strings.Join(types, "_"))

	filename := fmt.Sprintf("%s_%s_%s_%s.json", command, safeTerms, safeTypes, timestamp)
	filePath := filepath.Join(cfg.ResultsDir, filename)

	resultData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"command":   command,
		"terms":     terms,
		"types":     types,
		"results":   data,
	}

	jsonData, err := json.MarshalIndent(resultData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %v", err)
	}

	return nil
}

func GetResultsFilePath(command string, terms []string, types []string) string {
	cfg := Load()

	timestamp := time.Now().Format("20060102-150405")
	safeTerms := makeSafeFilename(strings.Join(terms, "_"))
	safeTypes := makeSafeFilename(strings.Join(types, "_"))

	filename := fmt.Sprintf("%s_%s_%s_%s.json", command, safeTerms, safeTypes, timestamp)
	return filepath.Join(cfg.ResultsDir, filename)
}

func makeSafeFilename(s string) string {
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " ", "@", "#", "$", "%", "^", "&", "*", "(", ")", "+", "=", "[", "]", "{", "}", "|", ";", ":", "'", ",", ".", "<", ">", "/", "?"}
	result := s
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func CreateSpinner(message string) *spinner.Spinner {
	cfg := Load()

	switch cfg.SpinnerStyle {
	case "dots":
		return spinner.WithDots(message)
	case "arrows":
		return spinner.WithArrows(message)
	case "bounce":
		return spinner.WithBounce(message)
	case "simple":
		return spinner.WithSimple(message)
	case "none":
		return nil
	default:
		return spinner.New(message)
	}
}
