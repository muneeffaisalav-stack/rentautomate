package main

import (
	"log"
	"os"

	"rentflow-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("rentflow.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Set up Gin router
	router := gin.Default()

	// Register handlers
	router.GET("/health", handlers.HealthCheck)
	router.POST("/api/manual/generate-invoices", handlers.GenerateInvoices)
	router.POST("/api/manual/send-reminders", handlers.SendReminders)
	router.GET("/api/invoices", handlers.GetInvoices)
	router.GET("/api/invoices/:id", handlers.GetInvoiceByID)
	router.POST("/api/webhooks/razorpay", handlers.ProcessRazorpayWebhook)
	router.GET("/test-firestore", handlers.TestFirestore)
	router.GET("/api/tenants", handlers.GetTenants)

	// Start the server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
