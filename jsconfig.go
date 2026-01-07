package main

import (
	"fmt"
	"os"
	"regexp"
)

// LoadConfigFromJS parses your existing JavaScript config files directly
// Works with files in the format:
//
//	module.exports = function () {
//	    return {
//	        server: 'value' OR server: { key: 'value', ... },
//	        questServer: 'value',
//	        ...
//	    }
//	}
func LoadConfigFromJS(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	config := &Config{}

	// DEBUG: Print first 500 chars of config file
	preview := content
	if len(preview) > 500 {
		preview = preview[:500]
	}
	fmt.Printf("\n[DEBUG] Config file content preview:\n%s\n[END DEBUG]\n\n", preview)

	// Try to extract server - could be a string or an object
	serverStringRe := regexp.MustCompile(`server:\s*['"]([^'"]+)['"]`)
	if matches := serverStringRe.FindStringSubmatch(content); len(matches) > 1 {
		fmt.Printf("[DEBUG] Found server as string: %s\n", matches[1])
		config.Server = matches[1]
	} else {
		// Extract server as an object by finding individual fields
		// Look for server: { ... } block and extract key-value pairs
		serverMap := make(map[string]interface{})
		
		// Known server object fields
		serverFields := []string{"quest", "agents", "bo", "tpv", "vault", "front"}
		
		for _, field := range serverFields {
			// Match patterns like: quest: 'value' or quest: "value"
			re := regexp.MustCompile(field + `:\s*['"]([^'"]+)['"]`)
			if matches := re.FindStringSubmatch(content); len(matches) > 1 {
				fmt.Printf("[DEBUG] Found %s: %s\n", field, matches[1])
				serverMap[field] = matches[1]
			} else {
				fmt.Printf("[DEBUG] NOT FOUND: %s\n", field)
			}
		}
		
		if len(serverMap) > 0 {
			config.Server = serverMap
		} else {
			fmt.Printf("[DEBUG] serverMap is empty!\n")
		}
	}

	// Simple extractors for top-level string values
	simpleExtractors := []struct {
		pattern string
		target  *string
	}{
		{`questServer:\s*['"]([^'"]+)['"]`, &config.QuestServer},
		{`questFront:\s*['"]([^'"]+)['"]`, &config.QuestFront},
		{`walkmeUrl:\s*['"]([^'"]+)['"]`, &config.WalkmeUrl},
	}

	for _, ext := range simpleExtractors {
		re := regexp.MustCompile(ext.pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			*ext.target = matches[1]
		}
	}

	// Extract Firebase config
	firebaseExtractors := []struct {
		pattern string
		target  *string
	}{
		{`apiKey:\s*['"]([^'"]+)['"]`, &config.Firebase.ApiKey},
		{`authDomain:\s*['"]([^'"]+)['"]`, &config.Firebase.AuthDomain},
		{`databaseURL:\s*['"]([^'"]+)['"]`, &config.Firebase.DatabaseURL},
		{`storageBucket:\s*['"]([^'"]+)['"]`, &config.Firebase.StorageBucket},
		{`messaginSenderId:\s*['"]([^'"]+)['"]`, &config.Firebase.MessagingSenderId}, // Note: typo in original
		{`messagingSenderId:\s*['"]([^'"]+)['"]`, &config.Firebase.MessagingSenderId},
	}

	for _, ext := range firebaseExtractors {
		re := regexp.MustCompile(ext.pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			*ext.target = matches[1]
		}
	}

	// Extract Google config
	googleExtractors := []struct {
		pattern string
		target  *string
	}{
		{`mapsKey:\s*['"]([^'"]+)['"]`, &config.Google.MapsKey},
		{`analytics:\s*['"]([^'"]+)['"]`, &config.Google.Analytics},
		{`recaptcha:\s*['"]([^'"]+)['"]`, &config.Google.Recaptcha},
	}

	for _, ext := range googleExtractors {
		re := regexp.MustCompile(ext.pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			*ext.target = matches[1]
		}
	}

	return config, nil
}

