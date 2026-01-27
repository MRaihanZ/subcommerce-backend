package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/MRaihanZ/subcommerce-backend/internal/entity"
	"github.com/joho/godotenv"
)

func init() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func RenderSubscriptionEmail(data entity.SubscriptionEmailData, htmlPath string) (string, error) {
	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderSubscriptionCancelationByUserEmail(data entity.CancelationSubscriptionByUserEmailData, htmlPath string) (string, error) {
	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderPayoutEmail(data entity.PayoutEmail, htmlPath string) (string, error) {
	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type BrevoSender struct {
	apiKey string
}

func NewBrevoSender() *BrevoSender {
	return &BrevoSender{
		apiKey: os.Getenv("BREVO_API_KEY"),
	}
}

func (b *BrevoSender) SendMail(to, subject, html string) error {
	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Subcommerce",
			"email": "no-reply@subcommerce.mraihanz.my.id",
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject":     subject,
		"htmlContent": html,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", b.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("brevo error: %s", resp.Status)
		return fmt.Errorf("brevo error: %s", resp.Status)
	}

	return nil
}
