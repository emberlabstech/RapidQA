//  ___                         _       __              _   _
// / __|_  _ _ __ _ __  ___ _ _| |_    / _|_  _ _ _  __| |_(_)___ _ _  ___
// \__ \ || | '_ \ '_ \/ _ \ '_|  _|  |  _| || | ' \/ _|  _| / _ \ ' \(_-<
// |___/\_,_| .__/ .__/\___/_|  \__|  |_|  \_,_|_||_\__|\__|_\___/_||_/__/
//          |_|  |_|

// This makes use of the gjson package. https://github.com/tidwall/gjson
// It allows for raw reading of JSON files in a very simple format of a dotted path.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func applicationRoot() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}

	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("resolve executable path: %w", err)
	}

	return filepath.Dir(exe), nil
}

//   ___           __ _
//  / __|___ _ _  / _(_)__ _
// | (__/ _ \ ' \|  _| / _` |
//  \___\___/_||_|_| |_\__, |
//                     |___/

// loadConfig - Load the configuration from a JSON file and apply default values.
func loadConfig(filename string) (*Config, error) {

	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", filename, err)
	}

	if !gjson.ValidBytes(raw) {
		return nil, fmt.Errorf("invalid JSON in %s", filename)
	}

	return &Config{raw: raw}, nil
}

// applyDefaults - Apply default values to the configuration.
func (c *Config) applyDefaults() error {
	defaults := []struct {
		path  string
		value interface{}
	}{
		{
			path:  "ollama.base_url",
			value: "http://127.0.0.1:11434",
		},
		{
			path:  "ollama.generate_path",
			value: "/api/generate",
		},
		{
			path:  "ollama.timeout_seconds",
			value: 300,
		},
		{
			path:  "sendgrid.api_key_env",
			value: "SENDGRID_API_KEY",
		},
	}

	var err error

	for _, item := range defaults {
		if gjson.GetBytes(c.raw, item.path).Exists() {
			continue
		}

		c.raw, err = sjson.SetBytes(c.raw, item.path, item.value)

		if err != nil {
			return fmt.Errorf("set config default %s: %w", item.path, err)
		}
	}

	return nil
}

// validate - Validate the configuration.
func (c *Config) validate() error {
	required := []string{
		"ollama.model",
		"files.extensions",
		"sendgrid.from_email",
		"sendgrid.to_email",
		"sendgrid.subject",
	}

	for _, path := range required {
		if strings.TrimSpace(c.String(path)) == "" {
			return fmt.Errorf("missing required configuration value %q", path)
		}
	}

	return nil
}

// String - Retrieve a string value from the configuration.
func (c *Config) String(path string) string {
	return gjson.GetBytes(c.raw, path).String()
}

// Int - Retrieve an integer value from the configuration.
func (c *Config) Int(path string) int64 {
	return gjson.GetBytes(c.raw, path).Int()
}

// Result - Retrieve a gjson.Result value from the configuration.
func (c *Config) Result(path string) gjson.Result {
	return gjson.GetBytes(c.raw, path)
}
