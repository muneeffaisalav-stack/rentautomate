package scheduler

import (
	"context"
	"log"
	"rentflow-backend/internal/config"
	"rentflow-backend/internal/firestore"
	"rentflow-backend/internal/services"
)

// GenerateInvoicesJob is the cron job for generating monthly invoices
func GenerateInvoicesJob() {
	log.Println("Starting GenerateInvoicesJob...")
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("ERROR: failed to load config in GenerateInvoicesJob: %v", err)
		return
	}
	log.Println("Config loaded successfully for GenerateInvoicesJob.")

	firestoreClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		log.Printf("ERROR: failed to create firestore client in GenerateInvoicesJob: %v", err)
		return
	}
	log.Println("Firestore client created successfully for GenerateInvoicesJob.")

	invoiceService := services.NewInvoiceService(firestoreClient)
	log.Println("Generating monthly invoices...")
	if err := invoiceService.GenerateMonthlyInvoices(ctx); err != nil {
		log.Printf("ERROR: failed to generate monthly invoices: %v", err)
	} else {
		log.Println("Successfully generated monthly invoices.")
	}
	log.Println("GenerateInvoicesJob finished.")
}

// SendRemindersJob is the cron job for sending rent reminders
func SendRemindersJob() {
	log.Println("Starting SendRemindersJob...")
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("ERROR: failed to load config in SendRemindersJob: %v", err)
		return
	}
	log.Println("Config loaded successfully for SendRemindersJob.")

	firestoreClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		log.Printf("ERROR: failed to create firestore client in SendRemindersJob: %v", err)
		return
	}
	log.Println("Firestore client created successfully for SendRemindersJob.")

	whatsappService := services.NewWhatsAppService(cfg)
	reminderService := services.NewReminderService(firestoreClient, whatsappService)

	log.Println("Sending rent reminders...")
	if err := reminderService.SendRentReminders(ctx); err != nil {
		log.Printf("ERROR: failed to send rent reminders: %v", err)
	} else {
		log.Println("Successfully sent rent reminders.")
	}
	log.Println("SendRemindersJob finished.")
}
