package rabbitmq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/mail"
	"time"

	"github.com/streadway/amqp"
	"github.com/tapiaw38/auth-api/internal/models"
)

// RabbitMQConnection is the connection to RabbitMQ
type RabbitMQConnection struct {
	Conn *amqp.Connection
}

// RabbitMQConfig is the RabbitMQ configuration
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// NewRabbitMQConfig creates a new RabbitMQ configuration
func NewRabbitMQConfig(conf *RabbitMQConfig) *RabbitMQConfig {
	return &RabbitMQConfig{
		Host:     conf.Host,
		Port:     conf.Port,
		User:     conf.User,
		Password: conf.Password,
	}
}

// Connection creates a new RabbitMQ connection
func (c *RabbitMQConfig) Connection() *RabbitMQConnection {
	conn, err := amqp.Dial("amqp://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}

	return &RabbitMQConnection{
		Conn: conn,
	}
}

// Channel creates a new RabbitMQ channel
func (c *RabbitMQConnection) Channel() (*amqp.Channel, error) {
	return c.Conn.Channel()
}

// Close closes the RabbitMQ connection
func (c *RabbitMQConnection) Close() {
	c.Conn.Close()
}

func (c *RabbitMQConnection) PublishJSONMessage(queueName string, message interface{}) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		})
	if err != nil {
		return err
	}

	return nil
}

// PublishEmailMessage publishes an email verification message to RabbitMQ
func (c *RabbitMQConnection) PublishEmailMessage(to string, from string, subject string, tempateName string, variables map[string]string) error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	fromEmail := mail.Address{
		Name:    "Mi Tour",
		Address: from,
	}
	toEmail := mail.Address{
		Name:    "",
		Address: to,
	}

	headers := map[string]string{
		"From":         fromEmail.String(),
		"To":           toEmail.String(),
		"Subject":      subject,
		"Content-Type": "text/html; charset=UTF-8",
	}

	message := bytes.Buffer{}
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	tmpl, err := template.ParseFiles("templates/" + tempateName + ".html")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&message, variables); err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"email_queue", // queue name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	msg := models.EmailMessage{
		To:        to,
		From:      from,
		Subject:   subject,
		Headers:   headers,
		Body:      message.String(),
		Variables: variables,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMsg,
		})
	if err != nil {
		return err
	}

	return nil
}

// ConsumeEmailMessage consumes an email verification message from RabbitMQ
func (c *RabbitMQConnection) ConsumeEmailMessage(sendEmail func(message, toEmail, fromEmail string) error) error {
	for {
		ch, err := c.Conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %s", err)
			time.Sleep(time.Second)
			continue
		}

		q, err := ch.QueueDeclare(
			"email_queue", // queue name
			false,         // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			log.Printf("Failed to declare a queue: %s", err)
			time.Sleep(time.Second)
			continue
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Printf("Failed to consume messages: %s", err)
			time.Sleep(time.Second)
			continue
		}

		for msg := range msgs {
			var email models.EmailMessage
			err := json.Unmarshal(msg.Body, &email)
			if err != nil {
				log.Printf("failed to unmarshal message: %s", err)
				continue
			}

			to := email.To
			from := email.From
			body := email.Body

			if from != "" || to != "" {
				if err := sendEmail(to, from, body); err != nil {
					log.Printf("Failed to send email: %s", err)
				}
				// Acknowledge the message to remove it from the queue
				msg.Ack(false)
			}
		}

		ch.Close()
	}
}
