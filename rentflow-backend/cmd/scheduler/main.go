package main

import (
	"log"
	"os"
	"rentflow-backend/internal/scheduler"
)

func main() {
	// Log to standard output, which is ideal for GitHub Actions
	log.SetOutput(os.Stdout)

	log.Println("Starting scheduler jobs...")

	log.Println("Running job: generate-invoices")
	scheduler.GenerateInvoicesJob()
	log.Println("Job 'generate-invoices' finished.")

	log.Println("Running job: send-reminders")
	scheduler.SendRemindersJob()
	log.Println("Job 'send-reminders' finished.")

	log.Println("All scheduler jobs finished successfully.")
}
