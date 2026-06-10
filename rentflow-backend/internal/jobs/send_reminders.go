package jobs

import (
	"context"
	"log"
	"rentflow-backend/internal/config"
	"rentflow-backend/internal/services"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// SendRemindersJob sends rent reminders to tenants
type SendRemindersJob struct {
	Config *config.Config
}

func (j *SendRemindersJob) Run() error {
	log.Println("Starting SendRemindersJob...")

	ctx := context.Background()
	config, err := config.LoadConfig(".")
	if err != nil {
		return err
	}

	log.Println("Config loaded successfully for SendRemindersJob.")

	// Initialize Firestore client
	client, err := firestore.NewClient(ctx, config.FirestoreProjectID, option.WithCredentialsFile(config.GoogleApplicationCredentials))
	if err != nil {
		log.Printf("ERROR: failed to create firestore client in SendRemindersJob: %v", err)
		return err
	}
	defer client.Close()

	log.Println("Firestore client created successfully for SendRemindersJob.")

	whatsappService := services.NewWhatsAppService(config)
	reminderService := services.NewReminderService(client, whatsappService)

	log.Println("Sending rent reminders...")
	if err := reminderService.SendRentReminders(ctx); err != nil {
		log.Printf("failed to send rent reminders: %v", err)
		return err
	}

	log.Println("Successfully sent rent reminders.")

	log.Println("SendRemindersJob finished.")
	return nil
}
