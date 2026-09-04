package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

//    _                               _
//   /_\  _ __ _ __    _ __  __ _ _ _| |_
//  / _ \| '_ \ '_ \  | '_ \/ _` | '_|  _|
// /_/ \_\ .__/ .__/  | .__/\__,_|_|  \__|
//       |_|  |_|     |_|

// run - Process the files and generate a review report to be printed or sent by email.
func run(files []string) error {
	instructions := []byte(prompt)
	if len(files) == 0 {
		return fmt.Errorf(
			"no files specified\nusage: commit-review <file> [file ...]",
		)
	}

	appRoot, err := applicationRoot()
	if err != nil {
		return err
	}

	// Get the default instructions, and override later if
	instructionsPath := filepath.Join(appRoot, Cfg.String(`instructionsPath`))

	if prompt == "" {
		r, err := os.ReadFile(instructionsPath)
		if err != nil {
			return fmt.Errorf("read instructions %s: %w", instructionsPath, err)
		}
		instructions = r
	}

	extensions := parseExtensions(
		Cfg.String("files.extensions"),
	)

	client := &http.Client{
		Timeout: time.Duration(Cfg.Int("ollama.timeout_seconds")) * time.Second,
	}

	results := make([]ReviewResult, 0, len(files))

	for _, filename := range files {
		result, processed := reviewFile(client, Cfg, string(instructions), extensions, filename)

		if processed {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		fmt.Println("No matching text files to process")

		return nil
	}

	report := buildReport(results)

	if *noMail {
		fmt.Print(report)
		return nil
	}

	if err := sendEmail(Cfg, report); err != nil {
		return err
	}

	fmt.Printf("Review sent for %d file(s)\n", len(results))

	return nil
}

// reviewFile - Process a single file and generate a review report.
func reviewFile(client *http.Client, cfg *Config, instructions string, extensions map[string]struct{}, filename string) (ReviewResult, bool) {

	displayName := filename

	fullPath, err := resolveInputFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", filename, err)

		return ReviewResult{}, false
	}

	if !allowedExtension(fullPath, extensions) {
		fmt.Fprintf(os.Stderr, "Skipping %s: extension not configured\n", filename)

		return ReviewResult{}, false
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", filename, err)

		return ReviewResult{}, false
	}

	if isBinary(data) {
		fmt.Fprintf(os.Stderr, "Skipping %s: binary or non-UTF-8 content\n", filename)

		return ReviewResult{}, false
	}

	fmt.Printf("Processing %s\n", fullPath)

	output, err := processWithOllama(client, cfg, instructions, displayName, string(data))

	return ReviewResult{Filename: displayName, Output: output, Err: err}, true
}

// resolveInputFile - Resolve the input file path and check if it is a regular file.
func resolveInputFile(filename string) (string, error) {

	filename = strings.TrimSpace(filename)

	if filename == "" {
		return "", fmt.Errorf("empty filename")
	}

	absolute, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}

	info, err := os.Stat(absolute)
	if err != nil {
		return "", err
	}

	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular file")
	}

	return absolute, nil
}

// parseExtensions - Parse the comma-separated list of extensions and return a map of extensions.
func parseExtensions(value string) map[string]struct{} {
	result := make(map[string]struct{})

	for _, extension := range strings.Split(value, ",") {

		extension = strings.TrimSpace(extension)
		extension = strings.TrimPrefix(extension, ".")
		extension = strings.ToLower(extension)
		if extension == "" {
			continue
		}

		result[extension] = struct{}{}
	}

	return result
}

// allowedExtension - Check if the file extension is allowed.
func allowedExtension(filename string, extensions map[string]struct{}) bool {

	extension := filepath.Ext(filename)

	if extension == "" {
		return false
	}

	extension = strings.TrimPrefix(extension, ".")
	extension = strings.ToLower(extension)
	_, found := extensions[extension]

	return found
}

// isBinary - Check if the file content is binary.
func isBinary(data []byte) bool {

	if len(data) == 0 {
		return false
	}

	sample := data

	if len(sample) > 8192 {
		sample = sample[:8192]
	}

	if bytes.IndexByte(sample, 0) >= 0 {
		return true
	}

	if !utf8.Valid(sample) {
		return true
	}

	controlCharacters := 0

	for _, b := range sample {
		if b < 0x20 &&
			b != '\n' &&
			b != '\r' &&
			b != '\t' &&
			b != '\f' {

			controlCharacters++
		}
	}

	return controlCharacters*100/len(sample) > 2
}

// buildReport - Generates a summary report from a list of review results and returns it as a formatted string.
func buildReport(results []ReviewResult) string {
	var report strings.Builder

	report.WriteString("Source review\n\n")

	for _, result := range results {
		report.WriteString(delimiter)
		report.WriteString("\nFILE: ")
		report.WriteString(result.Filename)
		report.WriteString("\n")
		report.WriteString(delimiter)
		report.WriteString("\n\n")

		if result.Err != nil {
			report.WriteString("ERROR: ")
			report.WriteString(result.Err.Error())
		} else {
			report.WriteString(result.Output)
		}

		report.WriteString("\n\n")
	}

	return report.String()
}
