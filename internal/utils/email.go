package utils

import (
	"crypto/tls"
	"net/smtp"
)

// EmailSMTPConfig is the configuration for the SMTP server
type EmailSMTPConfig struct {
	Host         string
	Port         string
	HostUser     string
	HostPassword string
}

// NewEmailSMTPConfig creates a new EmailSMTPConfig
func NewEmailSMTPConfig(config *EmailSMTPConfig) *EmailSMTPConfig {
	return &EmailSMTPConfig{
		Host:         config.Host,
		Port:         config.Port,
		HostUser:     config.HostUser,
		HostPassword: config.HostPassword,
	}
}

// SendEmail sends an email
func (conf *EmailSMTPConfig) SendEmail(toEmail, fromEmail, message string) error {

	host := conf.Host
	servername := conf.Host + ":" + conf.Port

	auth := smtp.PlainAuth("", conf.HostUser, conf.HostPassword, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(fromEmail); err != nil {
		return err
	}

	if err = client.Rcpt(toEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}
