package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config holds the application configuration

type Config struct {
	Port                   string
	FirestoreProjectID     string
	WhatsAppAPIToken       string
	WhatsAppAPIVersion     string
	WhatsAppPhoneNumberID string
	RazorpayWebhookSecret string
	GeminiAPIKey           string
}

// Load loads the configuration from environment variables

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Port:                   getEnv("PORT", "8080"),
		FirestoreProjectID:     getEnv("FIRESTORE_PROJECT_ID", ""),
		WhatsAppAPIToken:       getEnv("WHATSAPP_API_TOKEN", ""),
		WhatsAppAPIVersion:     getEnv("WHATSAPP_API_VERSION", "v18.0"),
		WhatsAppPhoneNumberID: getEnv("WHATSAPP_PHONE_NUMBER_ID", ""),
		RazorpayWebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),
		GeminiAPIKey:           getEnv("GEMINI_API_KEY", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
