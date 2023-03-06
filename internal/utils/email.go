package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"net/mail"
	"net/smtp"
)

type EmailSMTPConfig struct {
	Host         string
	Port         string
	HostUser     string
	HostPassword string
}

func NewEmailSMTPConfig(config *EmailSMTPConfig) *EmailSMTPConfig {
	return &EmailSMTPConfig{
		Host:         config.Host,
		Port:         config.Port,
		HostUser:     config.HostUser,
		HostPassword: config.HostPassword,
	}
}

func (conf *EmailSMTPConfig) SendEmail(to, subject, templateName string, variables map[string]interface{}, c chan error) {

	fromEmail := mail.Address{
		Name:    "Mi Tour",
		Address: conf.HostUser,
	}
	toEmail := mail.Address{
		Name:    "",
		Address: to,
	}

	subjectEmail := subject
	bodyEmail := variables

	headers := make(map[string]string)
	headers["From"] = fromEmail.String()
	headers["To"] = toEmail.String()
	headers["Subject"] = subjectEmail
	headers["content-type"] = "text/html; charset=UTF-8"

	message := ""

	for k, v := range headers {
		message += k + ": " + v + "\r\n"
	}

	t, err := template.ParseFiles("templates/" + templateName + ".html")
	if err != nil {
		c <- err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, bodyEmail)

	if err != nil {
		c <- err
	}

	message += "\r\n" + buf.String()

	host := conf.Host
	servername := conf.Host + ":" + conf.Port

	auth := smtp.PlainAuth("", conf.HostUser, conf.HostPassword, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		c <- err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		c <- err
	}

	if err = client.Auth(auth); err != nil {
		c <- err
	}

	if err = client.Mail(fromEmail.Address); err != nil {
		c <- err
	}

	if err = client.Rcpt(toEmail.Address); err != nil {
		c <- err
	}

	w, err := client.Data()
	if err != nil {
		c <- err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		c <- err
	}

	err = w.Close()
	if err != nil {
		c <- err
	}

	err = client.Quit()
	if err != nil {
		c <- err
	}

	c <- nil
}
