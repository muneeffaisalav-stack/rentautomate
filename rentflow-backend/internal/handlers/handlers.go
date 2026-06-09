package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"rentflow-backend/internal/config"
	"rentflow-backend/internal/firestore"
	"rentflow-backend/internal/models"
	"rentflow-backend/internal/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HealthCheck is a simple health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GenerateInvoices manually triggers the invoice generation process
func GenerateInvoices(c *gin.Context) {
	logFile, err := os.OpenFile("invoice_generation.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open log file"})
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load config: %v", err)})
		return
	}

	fsClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		log.Printf("Failed to create firestore client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create firestore client: %v", err)})
		return
	}
	defer fsClient.Close()

	invoiceService := services.NewInvoiceService(fsClient)

	if err := invoiceService.GenerateMonthlyInvoices(ctx); err != nil {
		log.Printf("Failed to generate invoices: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate invoices: %v", err)})
		return
	}

	log.Println("Invoice generation process completed.")
	c.JSON(http.StatusOK, gin.H{"message": "invoice generation started"})
}

// SendReminders manually triggers the reminder sending process
func SendReminders(c *gin.Context) {
	// TODO: Implement manual reminder sending
	c.JSON(http.StatusOK, gin.H{"message": "reminder sending started"})
}

// GetInvoices returns a list of invoices
func GetInvoices(c *gin.Context) {
	// TODO: Implement fetching invoices
	c.JSON(http.StatusOK, gin.H{"invoices": []string{}})
}

// GetInvoiceByID returns a single invoice by ID
func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load config: %v", err)})
		return
	}

	fsClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create firestore client: %v", err)})
		return
	}
	defer fsClient.Close()

	dsnap, err := fsClient.Collection("invoices").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get invoice: %v", err)})
		return
	}

	var invoice models.Invoice
	if err := dsnap.DataTo(&invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode invoice: %v", err)})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// GetTenants returns a list of tenants
func GetTenants(c *gin.Context) {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load config: %v", err)})
		return
	}

	fsClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create firestore client: %v", err)})
		return
	}
	defer fsClient.Close()

	var tenants []models.Tenant
	iter := fsClient.Collection("tenants").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get tenants: %v", err)})
			return
		}

		var tenant models.Tenant
		if err := doc.DataTo(&tenant); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode tenant: %v", err)})
			return
		}
		tenants = append(tenants, tenant)
	}

	c.JSON(http.StatusOK, tenants)
}

// ProcessRazorpayWebhook processes the Razorpay webhook
func ProcessRazorpayWebhook(c *gin.Context) {
	// TODO: Implement webhook processing
	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}

// TestFirestore is a handler to test the firestore connection
func TestFirestore(c *gin.Context) {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load config: %v", err)})
		return
	}

	fsClient, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create firestore client: %v", err)})
		return
	}
	defer fsClient.Close()

	// Perform a simple read operation to verify the connection
	// We'll try to get a document that likely doesn't exist.
	// If we get a "NotFound" error, it means the connection and authentication were successful.
	_, err = fsClient.Collection("test").Doc("test").Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusOK, gin.H{"message": "Firestore connection successful! (Verified by a test read operation)"})
			log.Println("Firestore connection successful!")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to perform test read on firestore: %v", err)})
		log.Printf("failed to perform test read on firestore: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Firestore connection successful! (Verified by a test read operation - and the test document exists)"})
	log.Println("Firestore connection successful!")
}
