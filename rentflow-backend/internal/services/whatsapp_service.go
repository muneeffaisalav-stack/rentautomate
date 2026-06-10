package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rentflow-backend/internal/config"
)

// IWhatsAppSender defines the interface for sending WhatsApp messages.
type IWhatsAppSender interface {
	SendRentReminder(tenantName, phone, amount, dueDate, paymentLink, templateName string) error
}

// WhatsAppService handles sending WhatsApp messages
type WhatsAppService struct {
	config *config.Config
	client *http.Client
}

// ensure WhatsAppService implements IWhatsAppSender
var _ IWhatsAppSender = (*WhatsAppService)(nil)

// NewWhatsAppService creates a new WhatsAppService
func NewWhatsAppService(config *config.Config) *WhatsAppService {
	return &WhatsAppService{
		config: config,
		client: &http.Client{},
	}
}

// SendRentReminder sends a rent reminder to a tenant
func (s *WhatsAppService) SendRentReminder(tenantName, phone, amount, dueDate, paymentLink, templateName string) error {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", s.config.WhatsAppAPIVersion, s.config.WhatsAppPhoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                phone,
		"type":              "template",
		"template": map[string]interface{}{
			"name": templateName,
			"language": map[string]string{
				"code": s.config.WhatsAppLanguageCode,
			},
			"components": []map[string]interface{}{
				{
					"type": "body",
					"parameters": []map[string]interface{}{
						{
							"type": "text",
							"text": tenantName,
						},
						{
							"type": "text",
							"text": amount,
						},
						{
							"type": "text",
							"text": dueDate,
						},
						{
							"type": "text",
							"text": paymentLink,
						},
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.WhatsAppAPIToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to send whatsapp message: %s", resp.Status)
		}
		return fmt.Errorf("failed to send whatsapp message: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
