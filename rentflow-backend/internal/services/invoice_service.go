package services

import (
	"context"
	"fmt"
	"log"
	"rentflow-backend/internal/models"
	"time"

	"cloud.google.com/go/firestore"
)

// InvoiceService handles invoice generation
type InvoiceService struct {
	firestore *firestore.Client
}

// NewInvoiceService creates a new InvoiceService
func NewInvoiceService(firestore *firestore.Client) *InvoiceService {
	return &InvoiceService{firestore: firestore}
}

// GenerateMonthlyInvoices generates monthly invoices for all active tenants
func (s *InvoiceService) GenerateMonthlyInvoices(ctx context.Context) error {
	tenants, err := s.firestore.Collection("tenants").Where("status", "==", "active").Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to get tenants: %v", err)
	}

	today := time.Now().Day()

	for _, tenantSnap := range tenants {
		var tenant models.Tenant
		if err := tenantSnap.DataTo(&tenant); err != nil {
			log.Printf("failed to decode tenant data: %v", err)
			continue
		}
		tenant.ID = tenantSnap.Ref.ID

		if tenant.DueDate == int64(today) {
			if err := s.generateInvoiceForTenant(ctx, &tenant); err != nil {
				log.Printf("failed to generate invoice for tenant %s: %v", tenant.ID, err)
			}
		}
	}

	return nil
}

func (s *InvoiceService) generateInvoiceForTenant(ctx context.Context, tenant *models.Tenant) error {
	now := time.Now()
	month := now.Format("2006-01")
	invoiceID := fmt.Sprintf("%s-%s", tenant.ID, month)

	// Check if invoice already exists
	_, err := s.firestore.Collection("invoices").Doc(invoiceID).Get(ctx)
	if err == nil {
		log.Printf("invoice already exists for tenant %s for month %s", tenant.ID, month)
		return nil // Not an error, just means it's already generated.
	}

	// Create invoice
	invoice := &models.Invoice{
		ID:         invoiceID,
		TenantID:   tenant.ID,
		LandlordID: tenant.LandlordID,
		PropertyID: tenant.PropertyID,
		Month:      month,
		Amount:     tenant.RentAmount,
		Status:     "pending",
		CreatedAt:  now.Format(time.RFC3339Nano),
	}

	_, err = s.firestore.Collection("invoices").Doc(invoiceID).Set(ctx, invoice)
	if err != nil {
		return fmt.Errorf("failed to create invoice: %v", err)
	}

	log.Printf("generated invoice for tenant %s for month %s", tenant.ID, month)

	return nil
}
