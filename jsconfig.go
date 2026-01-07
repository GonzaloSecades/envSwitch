package main

import (
	"os"
	"regexp"
)

// LoadConfigFromJS parses your existing JavaScript config files directly
// Works with files in the format:
//
//	module.exports = function () {
//	    return {
//	        server: 'value',
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

	// Simple extractors for top-level string values
	simpleExtractors := []struct {
		pattern string
		target  *string
	}{
		{`server:\s*['"]([^'"]+)['"]`, &config.Server},
		{`questServer:\s*['"]([^'"]+)['"]`, &config.QuestServer},
		{`questFront:\s*['"]([^'"]+)['"]`, &config.QuestFront},
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

