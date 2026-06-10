package scheduler

import (
	"context"
	"log"
	"rentflow-backend/internal/config"
	"rentflow-backend/internal/firestore"
	"rentflow-backend/internal/services"
)

// GenerateInvoicesJob is the cron job for generating monthly invoices
func GenerateInvoicesJob() error {
	log.Println("Starting GenerateInvoicesJob...")
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	firestoreClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		return err
	}

	invoiceService := services.NewInvoiceService(firestoreClient)
	if err := invoiceService.GenerateMonthlyInvoices(ctx); err != nil {
		return err
	}

	log.Println("GenerateInvoicesJob finished.")
	return nil
}

// SendRemindersJob is the cron job for sending rent reminders
func SendRemindersJob() error {
	log.Println("Starting SendRemindersJob...")
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	firestoreClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		return err
	}

	whatsappService := services.NewWhatsAppService(cfg)
	reminderService := services.NewReminderService(firestoreClient, whatsappService, cfg)

	if err := reminderService.SendRentReminders(ctx); err != nil {
		return err
	}

	log.Println("SendRemindersJob finished.")
	return nil
}
