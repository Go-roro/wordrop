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

const (
	fromEmailEnv = "EMAIL_SENDER_ADDRESS"
	passwordEnv  = "EMAIL_SENDER_PASSWORD"
	smtpHostEnv  = "SMTP_HOST"
	smtpPortEnv  = "SMTP_PORT"
)

type GmailSenderConfig struct {
	fromEmail string
	password  string
	smtpHost  string
	smtpPort  int
}

func NewMailSenderConfig() (*GmailSenderConfig, error) {
	return NewMailSenderConfigByKey(fromEmailEnv, passwordEnv, smtpHostEnv, smtpPortEnv)
}

func NewMailSenderConfigByKey(fromEmailKey, passwordKey, smtpHostKey, smtpPortKey string) (*GmailSenderConfig, error) {
	for _, envKeys := range []string{fromEmailKey, passwordKey, smtpHostKey, smtpPortKey} {
		if value, ok := os.LookupEnv(envKeys); value == "" || !ok {
			return nil, fmt.Errorf("environment variable %s is not set", envKeys)
		}
	}

	smtpPort, err := strconv.Atoi(os.Getenv(smtpPortKey))
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s in .env.email-test file to integer: %w", smtpPortKey, err)
	}

	return &GmailSenderConfig{
		fromEmail: os.Getenv(fromEmailKey),
		password:  os.Getenv(passwordKey),
		smtpHost:  os.Getenv(smtpHostKey),
		smtpPort:  smtpPort,
	}, nil
}

type GmailSender struct {
	dialer               *gomail.Dialer
	config               *GmailSenderConfig
	verificationTemplate *template.Template
}

func NewMailSender(config *GmailSenderConfig) (*GmailSender, error) {
	dialer := gomail.NewDialer(config.smtpHost, config.smtpPort, config.fromEmail, config.password)

	verificationTemplate, err := template.ParseFiles("template/verification.html")
	if err != nil {
		return nil, fmt.Errorf("could not parse verification template: %w", err)
	}

	return &GmailSender{
		dialer:               dialer,
		config:               config,
		verificationTemplate: verificationTemplate,
	}, nil
}

type VerificationTemplateData struct {
	Username         string
	VerificationLink string
}

func (gs *GmailSender) SendVerificationEmail(toEmail string, username string, verificationToken string) error {
	baseURL := os.Getenv("APP_BASE_URL")
	verificationLink := fmt.Sprintf("%s/subscriptions/verify?token=%s", baseURL, verificationToken)

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
	m.SetHeader("Subject", "Wordrop - 이메일 주소를 인증해주세요")
	m.SetBody("text/html", body.String())

	if err := gs.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("✅ Verification email sent successfully to", toEmail)
	return nil
}
