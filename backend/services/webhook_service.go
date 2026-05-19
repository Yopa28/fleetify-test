package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type WebhookPayload struct {
	ReportID uint   `json:"report_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

func SendStatusWebhook(reportID uint, status string) {
	webhookURL := os.Getenv("WEBHOOK_URL")

	if webhookURL == "" {
		log.Println("WEBHOOK_URL is not configured, skipping webhook")
		return
	}

	payload := WebhookPayload{
		ReportID: reportID,
		Status:   status,
		Message:  "Maintenance report status updated",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal webhook payload:", err)
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Post(
		webhookURL,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		log.Println("Failed to send webhook:", err)
		return
	}

	defer response.Body.Close()

	log.Println("Webhook sent with status:", response.Status)
}
