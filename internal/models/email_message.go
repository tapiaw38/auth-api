package models

type EmailMessage struct {
	To        string            `json:"to"`
	From      string            `json:"from"`
	Subject   string            `json:"subject"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Variables map[string]string `json:"variables"`
}
