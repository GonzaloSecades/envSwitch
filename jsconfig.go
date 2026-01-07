package main

import (
	"encoding/json"
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

	// Try to extract server as a string first
	serverStringRe := regexp.MustCompile(`server:\s*['"]([^'"]+)['"]`)
	if matches := serverStringRe.FindStringSubmatch(content); len(matches) > 1 {
		config.Server = matches[1]
	} else {
		// Try to extract server as an object (for envJs format)
		serverObjRe := regexp.MustCompile(`server:\s*(\{[^}]+\})`)
		if matches := serverObjRe.FindStringSubmatch(content); len(matches) > 1 {
			// Parse the JS object - convert single quotes to double quotes for JSON
			jsObj := matches[1]
			jsonObj := regexp.MustCompile(`'`).ReplaceAllString(jsObj, `"`)
			var serverMap map[string]interface{}
			if err := json.Unmarshal([]byte(jsonObj), &serverMap); err == nil {
				config.Server = serverMap
			}
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

