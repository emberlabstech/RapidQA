package main

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send email using sendgrid.
// The configuration comes from the config file section [sendgrid.[name]].
func sendEmail(cfg *Config, content string) error {

	apiKey := cfg.String("sendgrid.api_key")

	if apiKey == "" {
		envName := cfg.String("sendgrid.api_key_env")
		apiKey = os.Getenv(envName)

		if apiKey == "" {
			return fmt.Errorf("SendGrid API key is empty: set %s or sendgrid.api_key", envName)
		}
	}

	from := mail.NewEmail(cfg.String("sendgrid.from_name"), cfg.String("sendgrid.from_email"))
	to := mail.NewEmail(cfg.String("sendgrid.to_name"), cfg.String("sendgrid.to_email"))
	subject := cfg.String("sendgrid.subject")

	htmlContent := content
	//htmlContent := "<pre style=\"font-family:monospace;white-space:pre-wrap;\">" + html.EscapeString(content) + "</pre>"

	message := mail.NewSingleEmail(from, subject, to, content, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("send email through SendGrid: %w", err)
	}

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {

		return fmt.Errorf("SendGrid returned HTTP %d: %s", response.StatusCode, response.Body)
	}

	println("Sent email to ", to.Address)

	return nil
}
