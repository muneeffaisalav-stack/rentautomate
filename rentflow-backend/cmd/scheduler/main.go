package main

import (
	"log"
	"os"
	"rentflow-backend/internal/scheduler"
)

func main() {
	// Log to standard output, which is ideal for containerized environments
	log.SetOutput(os.Stdout)

	log.Println("Starting scheduler...")

	log.Println("Running job: generate-invoices")
	if err := scheduler.GenerateInvoicesJob(); err != nil {
		log.Printf("Job 'generate-invoices' failed: %v", err)
	} else {
		log.Println("Job 'generate-invoices' finished successfully.")
	}

	log.Println("Running job: send-reminders")
	if err := scheduler.SendRemindersJob(); err != nil {
		log.Printf("Job 'send-reminders' failed: %v", err)
	} else {
		log.Println("Job 'send-reminders' finished successfully.")
	}

	log.Println("Scheduler finished.")
}
