package services

import (
	"context"
	"fmt"
	"log"
	"rentflow-backend/internal/models"
	"time"

	"cloud.google.com/go/firestore"
)

// ReminderService handles sending rent reminders
type ReminderService struct {
	firestore       *firestore.Client
	whatsappService IWhatsAppSender
}

// NewReminderService creates a new ReminderService
func NewReminderService(firestore *firestore.Client, whatsappService IWhatsAppSender) *ReminderService {
	return &ReminderService{
		firestore:       firestore,
		whatsappService: whatsappService,
	}
}

// SendRentReminders sends rent reminders to tenants
func (s *ReminderService) SendRentReminders(ctx context.Context) error {
	invoices, err := s.firestore.Collection("invoices").Where("status", "==", "pending").Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get invoices: %v", err)
	}

	for _, invoiceSnap := range invoices {
		var invoice models.Invoice
		if err := invoiceSnap.DataTo(&invoice); err != nil {
			log.Printf("failed to decode invoice data: %v", err)
			continue
		}

		if err := s.sendReminderForInvoice(ctx, &invoice); err != nil {
			log.Printf("failed to send reminder for invoice %s: %v", invoice.ID, err)
		}
	}

	return nil
}

func (s *ReminderService) sendReminderForInvoice(ctx context.Context, invoice *models.Invoice) error {
	now := time.Now()

	// Check if TenantID is empty
	if invoice.TenantID == "" {
		return fmt.Errorf("invoice %s has an empty TenantID", invoice.ID)
	}

	// Get tenant to find out the due date
	tenantSnap, err := s.firestore.Collection("tenants").Doc(invoice.TenantID).Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %v", err)
	}
	var tenant models.Tenant
	if err := tenantSnap.DataTo(&tenant); err != nil {
		return fmt.Errorf("failed to decode tenant data: %v", err)
	}

	dueDate := tenant.DueDate

	// Send reminder on the due date
	if int64(now.Day()) == dueDate {
		if err := s.sendReminder(ctx, invoice, &tenant, "reminder_due_date"); err != nil {
			return err
		}
	}

	return nil
}

func (s *ReminderService) sendReminder(ctx context.Context, invoice *models.Invoice, tenant *models.Tenant, templateName string) error {

	// Send WhatsApp reminder
	if err := s.whatsappService.SendRentReminder(tenant.Name, tenant.Phone, fmt.Sprintf("%d", invoice.Amount), "", "", templateName);
	err != nil {
		return fmt.Errorf("failed to send whatsapp reminder: %v", err)
	}

	log.Printf("sent %s reminder for invoice %s", templateName, invoice.ID)

	return nil
}
