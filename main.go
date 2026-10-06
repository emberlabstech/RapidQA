package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const delimiter = "================================================================================"

type Config struct {
	raw []byte
}

type ReviewResult struct {
	Filename string
	Output   string
	Err      error
}

var instructionFile string // Override of instruction file source.
var noMail *bool
var asHtml *bool // Generate output as HTML.
var prompt = ""  // Manual CLI prompt.

func main() {
	// Flags
	noMail = flag.Bool("n", false, "print the review to console instead of sending email")
	asHtml = flag.Bool("h", false, "Generate the output as HTML instead of plaintext.")
	fPrompt := flag.String("p", "", "CLI AI Prompt in quotes")
	fInstruction := flag.String("i", "", "Instruction file to use for the review")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-nomail] <file> [file ...]\n", filepath.Base(os.Args[0]))

		flag.PrintDefaults()
	}

	flag.Parse()

	// ----------------------------------------------------------------------------------------
	// App root and configuration loading
	appRoot, err := applicationRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve application root: %v\n", err)
		os.Exit(1)
	}
	// Load the default config.
	configPath := filepath.Join(appRoot, "config", "config.json")
	Cfg, _ := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}
	if err := Cfg.applyDefaults(); err != nil {
		println(`Unable to apply defaults.`)
		os.Exit(1)
	}
	if err := Cfg.validate(); err != nil {
		println(`Bad validation of configuration file.`)
		os.Exit(1)
	}

	// ----------------------------------------------------------------------------------------
	// Set up the default instruction file path, and potentially override later.
	instructionFile = filepath.Join(appRoot, Cfg.String(`files.defaultInstruction`))
	println(`Instructions file:`, instructionFile)
	// ----------------------------------------------------------------------------------------
	// Overrides
	// Instructions -in flag, and this would be an absolute path, so override.
	if *fInstruction != "" {
		instructionFile = *fInstruction
	}
	// Prompt -p flag
	if *fPrompt != "" {
		prompt = *fPrompt
	}

	// ----------------------------------------------------------------------------------------
	// Get the list of remaining files listed on the command line behind the original flags.
	files := flag.Args()

	if err := run(Cfg, files); err != nil {
		fmt.Fprintf(os.Stderr, "commit-review: %v\n", err)

		os.Exit(1)
	}
}
