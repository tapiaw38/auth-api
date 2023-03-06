package utils

import (
	"context"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunConfig struct {
	Domain        string
	PrivateAPIKey string
}

func NewMailgunConfig(config *MailgunConfig) *MailgunConfig {
	return &MailgunConfig{
		Domain:        config.Domain,
		PrivateAPIKey: config.PrivateAPIKey,
	}
}

func (c *MailgunConfig) SendEmail(subjet string, body string, recipient string, templateName string, variables map[string]interface{}) error {

	sender := "example.com"
	mg := mailgun.NewMailgun(c.Domain, c.PrivateAPIKey)

	message := mg.NewMessage(sender, subjet, body, recipient)
	message.SetTemplate(templateName)

	for key, value := range variables {
		message.AddTemplateVariable(key, value)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ID: %s Resp: %s", id, resp)

	return nil
}
