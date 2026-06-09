# Rentflow Backend

This is the backend for Rentflow, a rent management automation tool.

## Features

- Automated monthly invoice generation
- Automated rent reminders via WhatsApp
- Payment gateway integration with Razorpay
- REST API for managing invoices and tenants

## Getting Started

### Prerequisites

- Go 1.18 or higher
- A Google Cloud project with Firestore enabled
- A WhatsApp Business Account with a configured phone number
- A Razorpay account

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/your-username/rentflow-backend.git
    ```

2.  Navigate to the project directory:

    ```bash
    cd rentflow-backend
    ```

3.  Create a `.env` file from the example:

    ```bash
    cp .env.example .env
    ```

4.  Update the `.env` file with your credentials.

5.  Install the dependencies:

    ```bash
    go mod tidy
    ```

### Running the Application

-   **API Server**:

    ```bash
    go run cmd/api/main.go
    ```

-   **Scheduler**:

    ```bash
    go run cmd/scheduler/main.go
    ```

## API Endpoints

-   `GET /health`: Health check
-   `POST /manual/generate-invoices`: Manually trigger invoice generation
-   `POST /manual/send-reminders`: Manually trigger reminder sending
-   `GET /invoices`: Get all invoices
-   `GET /invoices/:id`: Get a single invoice
-   `POST /webhooks/razorpay`: Razorpay webhook
