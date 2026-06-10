package models

// User represents a user in the system, who can be a landlord or other roles.

type User struct {
	ID        string `firestore:"id,omitempty"`
	Name      string `firestore:"name,omitempty"`
	Email     string `firestore:"email,omitempty"`
	Role      string `firestore:"role,omitempty"`
	UpiID     string `firestore:"upiId,omitempty"`
	CreatedAt string `firestore:"createdAt,omitempty"`
}

// Property represents a property in the system

type Property struct {
	ID           string `firestore:"id,omitempty"`
	LandlordID   string `firestore:"landlordId,omitempty"`
	PropertyName string `firestore:"propertyName,omitempty"`
	Address      string `firestore:"address,omitempty"`
	CreatedAt    string `firestore:"createdAt,omitempty"`
}

// Tenant represents a tenant in the system

type Tenant struct {
	ID         string `firestore:"id,omitempty"`
	LandlordID string `firestore:"landlordId,omitempty"`
	PropertyID string `firestore:"propertyId,omitempty"`
	Name       string `firestore:"name,omitempty"`
	Phone      string `firestore:"phone,omitempty"`
	RentAmount int64  `firestore:"rentAmount,omitempty"`
	DueDate    int64  `firestore:"dueDate,omitempty"`
	UpiID      string `firestore:"upiId,omitempty"`
	Status     string `firestore:"status,omitempty"`
}

// Invoice represents an invoice in the system

type Invoice struct {
	ID         string `firestore:"id,omitempty"`
	TenantID   string `firestore:"tenantId,omitempty"`
	LandlordID string `firestore:"landlordId,omitempty"`
	PropertyID string `firestore:"propertyId,omitempty"`
	Month      string `firestore:"month,omitempty"`
	Amount     int64  `firestore:"amount,omitempty"`
	Status     string `firestore:"status,omitempty"`
	CreatedAt  string `firestore:"createdAt,omitempty"`
}
