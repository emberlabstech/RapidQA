onePassVault 	?= CICD
onePassItem 	?= Labs-general
env				?= dev
target			?= commit-review
 
all:
	mkdir -p bin
	go build -o $(target) .
	make mapConfig

	
mapConfig: ## Update the config, using 1Password to fill in the tags from the template config.
	rm -f config/config.json
	1PassMapper -vv -prefix $(env) -vault $(onePassVault) -item $(onePassItem) -in config/config-template.json -out config/config.json

init:
	go mod init $(target)

setup:
	go get -u "github.com/sendgrid/sendgrid-go"
	go get -u "github.com/sendgrid/sendgrid-go/helpers/mail"
	go get -u "github.com/tidwall/gjson"
	go get -u "github.com/tidwall/sjson"

clean: 
	rm config/config.json
	rm $(target) 
	
