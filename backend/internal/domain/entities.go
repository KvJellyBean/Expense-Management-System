package domain

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Expense struct {
	ID                int        `json:"id"`
	UserID            int        `json:"user_id"`
	AmountIDR         int        `json:"amount_idr"`
	Description       string     `json:"description"`
	ReceiptURL        *string    `json:"receipt_url,omitempty"`
	Status            string     `json:"status"`
	AutoApproved      bool       `json:"auto_approved"`
	SubmittedAt       time.Time  `json:"submitted_at"`
	ProcessedAt       *time.Time `json:"processed_at,omitempty"`
	PaymentID         *string    `json:"payment_id,omitempty"`
	PaymentExternalID *string    `json:"payment_external_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Approval          *Approval  `json:"approval,omitempty"`
}

type Approval struct {
	ID         int       `json:"id"`
	ExpenseID  int       `json:"expense_id"`
	ApproverID int       `json:"approver_id"`
	Status     string    `json:"status"`
	Notes      *string   `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuditLog struct {
	ID        int                    `json:"id"`
	ExpenseID int                    `json:"expense_id"`
	UserID    *int                   `json:"user_id,omitempty"`
	Action    string                 `json:"action"`
	OldStatus *string                `json:"old_status,omitempty"`
	NewStatus *string                `json:"new_status,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}
