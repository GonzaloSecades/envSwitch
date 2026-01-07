package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Config represents the environment configuration
type Config struct {
	Server      string       `json:"server"`
	QuestServer string       `json:"questServer"`
	QuestFront  string       `json:"questFront"`
	Firebase    FirebaseConf `json:"firebase"`
	Google      GoogleConf   `json:"google"`
}

type FirebaseConf struct {
	ApiKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	DatabaseURL       string `json:"databaseURL"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderId string `json:"messagingSenderId"`
}

type GoogleConf struct {
	MapsKey   string `json:"mapsKey"`
	Analytics string `json:"analytics"`
	Recaptcha string `json:"recaptcha"`
}

// Replacement defines a regex pattern and its replacement value
type Replacement struct {
	Pattern     *regexp.Regexp
	Replacement string
}

func main() {
	// CLI flags
	env := flag.String("env", "", "Environment name (test, stress, cfg, prod, etc.)")
	configDir := flag.String("config-dir", "./configs", "Directory containing config.{env}.json files")
	targetFile := flag.String("target", "./app/shared/services/web/serverConfig.js", "Target file to modify")
	isDist := flag.Bool("dist", false, "Set isDist to true")
	useJS := flag.Bool("js", false, "Use .js config files instead of .json (parses your existing JS configs)")
	dryRun := flag.Bool("dry-run", false, "Show what would be changed without modifying the file")
	interactive := flag.Bool("i", false, "Run in interactive mode with visual CLI")
	flag.Parse()

	// Interactive mode
	if *interactive {
		if err := RunInteractiveCLI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running interactive CLI: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *env == "" {
		fmt.Fprintln(os.Stderr, "Error: --env flag is required (or use -i for interactive mode)")
		fmt.Fprintln(os.Stderr, "Usage: envswitch --env test [--config-dir ./configs] [--target ./path/to/serverConfig.js] [--dist] [--js] [--dry-run]")
		fmt.Fprintln(os.Stderr, "       envswitch -i  (interactive mode)")
		os.Exit(1)
	}

	// Load config file (JSON or JS)
	var configPath string
	var config *Config
	var err error

	if *useJS {
		configPath = filepath.Join(*configDir, fmt.Sprintf("config.%s.js", *env))
		config, err = LoadConfigFromJS(configPath)
	} else {
		configPath = filepath.Join(*configDir, fmt.Sprintf("config.%s.json", *env))
		config, err = loadConfig(configPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config %s: %v\n", configPath, err)
		os.Exit(1)
	}

	// Read target file
	content, err := os.ReadFile(*targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading target file %s: %v\n", *targetFile, err)
		os.Exit(1)
	}

	// Build replacements (preserve trailing commas where they exist)
	result := applyReplacements(string(content), config, *isDist)

	// Dry-run mode: show diff and exit
	if *dryRun {
		fmt.Printf("Dry-run mode - showing changes for environment: %s\n", *env)
		fmt.Printf("Config: %s\n", configPath)
		fmt.Printf("Target: %s\n\n", *targetFile)
		printDiff(string(content), result)
		return
	}

	// Write back to file
	err = os.WriteFile(*targetFile, []byte(result), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing target file %s: %v\n", *targetFile, err)
		os.Exit(1)
	}

	fmt.Printf("✓ Switched to environment: %s\n", *env)
	fmt.Printf("  Config: %s\n", configPath)
	fmt.Printf("  Target: %s\n", *targetFile)
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// applyReplacements applies all environment-specific replacements to content
func applyReplacements(content string, config *Config, isDist bool) string {
	replacements := []Replacement{
		{
			Pattern:     regexp.MustCompile(`baseUrl:\s*['"][^'"]*['"],?`),
			Replacement: fmt.Sprintf(`baseUrl: "%s",`, config.Server),
		},
		{
			Pattern:     regexp.MustCompile(`questUrl:\s*['"][^'"]*['"],?`),
			Replacement: fmt.Sprintf(`questUrl: "%s",`, config.QuestServer),
		},
		{
			Pattern:     regexp.MustCompile(`questFront:\s*['"][^'"]*['"],?`),
			Replacement: fmt.Sprintf(`questFront: "%s",`, config.QuestFront),
		},
		{
			Pattern:     regexp.MustCompile(`isDist:\s*(true|false),?`),
			Replacement: fmt.Sprintf(`isDist: %t,`, isDist),
		},
		{
			Pattern:     regexp.MustCompile(`recaptchaApiKey:\s*['"][^'"]*['"],?`),
			Replacement: fmt.Sprintf(`recaptchaApiKey: "%s",`, config.Google.Recaptcha),
		},
	}

	result := content
	for _, r := range replacements {
		result = r.Pattern.ReplaceAllString(result, r.Replacement)
	}

	return result
}

// Helper to print what will be replaced (dry-run mode)
func printDiff(original, modified string) {
	origLines := strings.Split(original, "\n")
	modLines := strings.Split(modified, "\n")

	for i := 0; i < len(origLines) && i < len(modLines); i++ {
		if origLines[i] != modLines[i] {
			fmt.Printf("- %s\n", origLines[i])
			fmt.Printf("+ %s\n", modLines[i])
		}
	}
}

