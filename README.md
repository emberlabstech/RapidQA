# RapidQA
Author: 2026 - EmberLabs / Chris Sprucefield.

Licensed for public use under MIT License. 

## Dependencies for building

* Go 1.20 or later
* Ollama (local instance)

#### For email sending
* Sendgrid API key

#### For credentials used in the config.json
* 1Password account for secrets
* 1Password API key.
* 1PassMapper https://github.com/emberlabstech/1PassMapper

## Arguments

`commit-review <flags> <file1> <file2> <file3> ...`

Standard wildcards are permissible. 

## GIT Post-Commit hook. 

The sample post commit shellscript can be found in post-commit.sh

It will analyze the files committed by the Git commit. 


## Flags

    -n          Print the review to console instead of sending email
    -p  string  "CLI AI Prompt in quotes"
    -i  string  Instruction file to use for the review (full path) if not default. 
                                                              

## Configuration

```json
{
  "ollama": {
    "base_url": "http://127.0.0.1:11434",
    "generate_path": "/api/generate",
    "model": "qwen3-coder:30b",
    "timeout_seconds": 600,
    "keep_alive": "15m",
    "temperature": 0.2,
    "num_ctx": 262144
  },
  "files": {
    "defaultInstruction": "config/instructions.txt",
    "extensions": "go,txt,css,js,ts,html,htm,json,xml,yaml,yml,sql,md,sh"
  },
  "sendgrid": {
    "api_key": "<sendgrid_api_key>",
    "from_email": "<sender_email>",
    "from_name": "<sender_name>",
    "to_email": "<default_recipient_email>",
    "to_name": "<default_recipient_name>",
    "subject": "<email_subject>"
  }
}
```

## Files

```plain text
.
├── commit-review                   Compiled binary
├── config
│   ├── config.json                 Configuration file generated from the template.  
│   ├── config-template.json        Configuration file template. 
│   └── instructions.txt            Default instruction file for the review. 
│    
├── email.go                        The sendmail component for SendGrid. 
├── go.mod
├── go.sum
├── main.go                         Main application. 
├── Makefile                        Build file. 
├── ollama.go                       Ollama Processor. 
├── post-commit.sh                  Shell script to support git commits as a post commit event. 
├── process.go                      Process the files and compile the report. 
├── README.md                       Readme.
└── support.go                      Support functions for the other parts of the application. 
```

## Makefile

Modify the initial parameters of the Makefile to use your defaults.

```plain text
onePassVault 	?= CICD
onePassItem 	?= Labs-general
env				?= dev
```
This build makes use of the 1PassMapper from EmberLabs.
https://github.com/emberlabstech/1passmapper

### make
Builds the application and produces the output Binary.

### make mapconfig
Builds the configuration file using one pass mapper, 
using the defaults from the initial settings in the make file. 


### make init
Initialize the module and name the package as commit-review. 

### make setup
Installs the required libraries for it to compile. SendGrid and gjson, sjson. 
