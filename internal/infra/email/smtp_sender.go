package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailSenderConfig struct {
	fromEmail string
	password  string
	smtpHost  string
	smtpPort  int
}

func NewMailSenderConfigByKey(fromEmailKey, passwordKey, smtpHostKey, smtpPortKey string) (*MailSenderConfig, error) {
	for _, envKeys := range []string{fromEmailKey, passwordKey, smtpHostKey, smtpPortKey} {
		if value, ok := os.LookupEnv(envKeys); value == "" || !ok {
			return nil, fmt.Errorf("environment variable %s is not set", envKeys)
		}
	}

	smtpPort, err := strconv.Atoi(os.Getenv(smtpPortKey))
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s in .env.email-test file to integer: %w", smtpPortKey, err)
	}

	return &MailSenderConfig{
		fromEmail: os.Getenv(fromEmailKey),
		password:  os.Getenv(passwordKey),
		smtpHost:  os.Getenv(smtpHostKey),
		smtpPort:  smtpPort,
	}, nil
}

type MailSender struct {
	dialer               *gomail.Dialer
	config               *MailSenderConfig
	verificationTemplate *template.Template
}

func NewMailSender(config *MailSenderConfig) (*MailSender, error) {
	dialer := gomail.NewDialer(config.smtpHost, config.smtpPort, config.fromEmail, config.password)

	verificationTemplate, err := template.ParseFiles("template/verification.html")
	if err != nil {
		return nil, fmt.Errorf("could not parse verification template: %w", err)
	}

	return &MailSender{
		dialer:               dialer,
		config:               config,
		verificationTemplate: verificationTemplate,
	}, nil
}

type VerificationTemplateData struct {
	Username         string
	VerificationLink string
}

func (gs *MailSender) SendVerificationEmail(toEmail string, username string, verificationToken string) error {
	baseURL := os.Getenv("APP_BASE_URL")
	verificationLink := fmt.Sprintf("%s/verify?token=%s", baseURL, verificationToken)

	data := VerificationTemplateData{
		Username:         username,
		VerificationLink: verificationLink,
	}

	var body bytes.Buffer
	err := gs.verificationTemplate.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", gs.config.fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Welcome to Wordrop! Please Verify Your Email")
	m.SetBody("text/html", body.String())

	if err := gs.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("✅ Verification email sent successfully to", toEmail)
	return nil
}
