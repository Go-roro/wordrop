package testhelper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestMailServer struct {
	Container testcontainers.Container
	SmtpHost  string
	SmtpPort  string
	ApiUrl    string
}

// SetupTestMailServer sets up a MailHog test container for email testing.
func SetupTestMailServer() *TestMailServer {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Println("Setting up MailHog TestContainer...")
	smtpPort := "1025/tcp"
	apiPort := "8025/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "mailhog/mailhog:v1.0.1",
		ExposedPorts: []string{smtpPort, apiPort}, // SMTP and API ports
		WaitingFor:   wait.ForLog("Creating API v2 with WebPath"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start MailHog container: %s", err)
	}

	// Get the mapped SMTP port (1025)
	smtpMappedPort, err := container.MappedPort(ctx, nat.Port(smtpPort))
	if err != nil {
		log.Fatalf("Failed to get mapped SMTP port: %s", err)
	}

	// Get the mapped API port (8025)
	apiMappedPort, err := container.MappedPort(ctx, nat.Port(apiPort))
	if err != nil {
		log.Fatalf("Failed to get mapped API port: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %s", err)
	}

	apiUrl := fmt.Sprintf("http://%s:%s", host, apiMappedPort.Port()) // localhost:8025
	if err = os.Setenv("SMTP_HOST", host); err != nil {
		log.Fatalf("Failed to set SMTP_HOST environment variable: %s", err)
	}
	if err = os.Setenv("SMTP_PORT", smtpMappedPort.Port()); err != nil {
		log.Fatalf("Failed to set SMTP_PORT environment variable: %s", err)
	}

	log.Printf("✅ MailHog container running. SMTP on %s:%s, API on %s:%s", host, smtpMappedPort.Port(), host, apiMappedPort.Port())
	return &TestMailServer{
		Container: container,
		SmtpHost:  host,
		SmtpPort:  smtpPort,
		ApiUrl:    apiUrl,
	}
}

func (ms *TestMailServer) TearDown() {
	log.Println("Tearing down MailHog container...")
	if err := ms.Container.Terminate(context.Background()); err != nil {
		log.Printf("Failed to terminate MailHog container: %s", err)
	}
}

// CleanUp https://github.com/mailhog/MailHog/blob/master/docs/APIv1.md#delete-apiv1messages
func (ms *TestMailServer) CleanUp() {
	log.Println("Cleaning up mails...")
	apiUrl := fmt.Sprintf("%s/api/v1/messages", ms.ApiUrl)
	req, _ := http.NewRequest(http.MethodDelete, apiUrl, nil)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to clean up mails: %s", err)
	}
}
