package services

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"rentflow-backend/internal/config"
	"rentflow-backend/internal/models"
	"time"

	"cloud.google.com/go/firestore"
)

// ReminderService handles sending rent reminders
type ReminderService struct {
	firestore       *firestore.Client
	whatsappService IWhatsAppSender
	config          *config.Config
}

// NewReminderService creates a new ReminderService
func NewReminderService(firestore *firestore.Client, whatsappService IWhatsAppSender, config *config.Config) *ReminderService {
	return &ReminderService{
		firestore:       firestore,
		whatsappService: whatsappService,
		config:          config,
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
			return fmt.Errorf("failed to send reminder for invoice %s: %v", invoice.ID, err)
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
		if err := s.sendReminder(ctx, invoice, &tenant, s.config.WhatsAppDueDateTemplate); err != nil {
			return err
		}
	}

	return nil
}

func (s *ReminderService) sendReminder(ctx context.Context, invoice *models.Invoice, tenant *models.Tenant, templateName string) error {
	if invoice.LandlordID == "" {
		return fmt.Errorf("invoice %s has an empty LandlordID", invoice.ID)
	}

	// Get landlord to find out the UPI ID
	landlordSnap, err := s.firestore.Collection("users").Doc(invoice.LandlordID).Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get landlord: %v", err)
	}
	var landlord models.User
	if err := landlordSnap.DataTo(&landlord); err != nil {
		return fmt.Errorf("failed to decode landlord data: %v", err)
	}

	if landlord.UpiID == "" {
		return fmt.Errorf("landlord %s has an empty UpiID", landlord.ID)
	}

	dueDateStr := fmt.Sprintf("%d", tenant.DueDate)
	amountStr := fmt.Sprintf("%d", invoice.Amount)

	note := fmt.Sprintf("Rent for %s for %s", tenant.Name, invoice.Month)
	paymentLink := fmt.Sprintf("upi://pay?pa=%s&pn=%s&am=%s&cu=INR&tn=%s", landlord.UpiID, landlord.Name, amountStr, url.QueryEscape(note))

	log.Printf("Sending reminder to recipient: %s", tenant.Phone)

	// Send WhatsApp reminder
	if err := s.whatsappService.SendRentReminder(tenant.Name, tenant.Phone, amountStr, dueDateStr, paymentLink, templateName);
	err != nil {
		return fmt.Errorf("failed to send whatsapp reminder: %v", err)
	}

	log.Printf("sent %s reminder for invoice %s", templateName, invoice.ID)

	return nil
}
