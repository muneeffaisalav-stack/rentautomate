package main

import (
	"context"
	"log"
	"os"
	"rentflow-backend/internal/services"
	"testing"

	"cloud.google.com/go/firestore"
)

// MockWhatsAppSender is a mock implementation of the IWhatsAppSender interface
type MockWhatsAppSender struct{}

// SendRentReminder is a mock implementation of the SendRentReminder method
func (s *MockWhatsAppSender) SendRentReminder(to, name, amount, dueDate, paymentLink, templateName string) error {
	log.Printf("Sending rent reminder to %s (%s) with template %s", name, to, templateName)
	return nil
}

func TestScheduler(t *testing.T) {
	ctx := context.Background()

	// Initialize Firestore client
	firestoreClient, err := firestore.NewClient(ctx, "studio-2364224586-e182c")
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// Initialize services
	invoiceService := services.NewInvoiceService(firestoreClient)
	reminderService := services.NewReminderService(firestoreClient, &MockWhatsAppSender{})

	// Set up logger
	log.SetOutput(os.Stdout)

	log.Println("Starting manual scheduler run...")

	log.Println("Running monthly invoice generation...")
	if err := invoiceService.GenerateMonthlyInvoices(ctx); err != nil {
		log.Fatalf("Failed to generate monthly invoices: %v", err)
	}

	log.Println("Running rent reminders...")
	if err := reminderService.SendRentReminders(ctx); err != nil {
		log.Fatalf("Failed to send rent reminders: %v", err)
	}

	log.Println("Manual scheduler run finished.")
}
