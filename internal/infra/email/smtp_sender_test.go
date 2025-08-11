package email

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/Go-roro/wordrop/internal/infra/testhelper"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

const testEnvPath = "../testhelper/.env.email-test"
const (
	fromEmailEnv = "EMAIL_SENDER_ADDRESS"
	passwordEnv  = "EMAIL_SENDER_PASSWORD"
	smtpHostEnv  = "SMTP_HOST"
	smtpPortEnv  = "SMTP_PORT"
)

type EmailSenderTestSuite struct {
	suite.Suite
	sender     *GmailSender
	mailServer *testhelper.TestMailServer
}

func (suite *EmailSenderTestSuite) SetupSuite() {
	log.Println("Setting up EmailServiceTestSuite...")
	err := godotenv.Load(testEnvPath)
	if err != nil {
		log.Fatal("Error loading .env.email-test file")
	}

	suite.mailServer = testhelper.SetupTestMailServer()

	config, err := NewMailSenderConfigByKey(fromEmailEnv, passwordEnv, smtpHostEnv, smtpPortEnv)
	if err != nil {
		log.Fatalf("Failed to create GmailSenderConfig: %v", err)
	}

	suite.sender, err = NewMailSender(config)
	if err != nil {
		log.Fatalf("Failed to create GmailSender: %v", err)
	}
}

func (suite *EmailSenderTestSuite) TearDownSuite() {
	log.Println("Tearing down EmailServiceTestSuite...")
	suite.mailServer.TearDown()
}

func (suite *EmailSenderTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("Before test: %s - %s\n", suiteName, testName)
	suite.mailServer.CleanUp()
}

func TestEmailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderTestSuite))
}

func (suite *EmailSenderTestSuite) TestSendVerificationEmail() {
	suite.Run("TestEmailSender_SendVerificationEmail", func() {
		toEmail := "new-user@example.com"
		username := "test-user"
		token := "test-verification-token"
		err := suite.sender.SendVerificationEmail(toEmail, username, token)
		suite.Require().NoError(err, "Expected no error when sending verification email")

		//time.Sleep(200 * time.Millisecond) // Give MailHog a moment to process

		apiUrl := fmt.Sprintf("%s/api/v2/messages", suite.mailServer.ApiUrl)
		resp, err := http.Get(apiUrl)
		suite.NoError(err, "Failed to connect to MailHog API")
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var mailHogResp MailHogResponse

		if err := json.Unmarshal(body, &mailHogResp); err != nil {
			return
		}
		suite.Equal(1, mailHogResp.Total, "Expected 1 email to be caught by MailHog")

		// Verify email content
		latestEmail := mailHogResp.Items[0]
		suite.Contains(latestEmail.Content.Body, username, "Email body should contain the correct username")
		suite.Contains(latestEmail.Content.Body, token, "Email body should contain the correct verification link")
	})
}

// https://github.com/mailhog/MailHog/blob/master/docs/APIv2/swagger-2.0.json
type MailHogResponse struct {
	Total int              `json:"total"`
	Items []MailHogMessage `json:"items"`
}

type MailHogMessage struct {
	ID      string `json:"ID"`
	To      []any  `json:"To"`
	Content struct {
		Body string `json:"Body"`
	} `json:"Content"`
}
