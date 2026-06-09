package main

import (
	"log"
	"os"
	"rentflow-backend/internal/scheduler"

	"github.com/robfig/cron/v3"
)

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("scheduler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	
	c := cron.New()

	// Schedule cron jobs
	// Every month on the 1st day
	_, err = c.AddFunc("* * * * *", scheduler.GenerateInvoicesJob)
	if err != nil {
		log.Fatalf("failed to add generate invoices job: %v", err)
	}

	// Every day at 9am
	_, err = c.AddFunc("0 9 * * *", scheduler.SendRemindersJob)
	if err != nil {
		log.Fatalf("failed to add send reminders job: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	log.Println("scheduler started")

	// Keep the scheduler running
	select {}
}
