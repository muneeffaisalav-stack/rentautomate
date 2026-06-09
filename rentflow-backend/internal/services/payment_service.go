package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rentflow-backend/internal/config"

	"cloud.google.com/go/firestore"
)

// PaymentService handles payment gateway webhooks

type PaymentService struct {
	firestore *firestore.Client
	config    *config.Config
}

// NewPaymentService creates a new PaymentService

func NewPaymentService(firestore *firestore.Client, config *config.Config) *PaymentService {
	return &PaymentService{
		firestore: firestore,
		config:    config,
	}
}

// ProcessRazorpayWebhook processes incoming webhooks from Razorpay

func (s *PaymentService) ProcessRazorpayWebhook(ctx context.Context, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}

	if !s.verifyWebhookSignature(r.Header.Get("X-Razorpay-Signature"), string(body)) {
		return fmt.Errorf("invalid webhook signature")
	}

	// Parse the webhook payload
	// ...

	// Find the invoice
	// ...

	// Update the invoice status
	// ...

	log.Printf("processed razorpay webhook")

	return nil
}

func (s *PaymentService) verifyWebhookSignature(signature, payload string) bool {
	mac := hmac.New(sha256.New, []byte(s.config.RazorpayWebhookSecret))
	mac.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
