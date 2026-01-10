package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense *Expense) error
	GetByID(ctx context.Context, id int) (*Expense, error)
	GetByUserID(ctx context.Context, userID int, status string, limit, offset int) ([]*Expense, int, error)
	GetAll(ctx context.Context, status string, limit, offset int) ([]*Expense, int, error)
	GetPendingApprovals(ctx context.Context, limit, offset int) ([]*Expense, int, error)
	Update(ctx context.Context, expense *Expense) error
	UpdateStatus(ctx context.Context, id int, status string, processedAt *string) error
	UpdatePaymentInfo(ctx context.Context, id int, paymentID, externalID string) error
}

type ApprovalRepository interface {
	Create(ctx context.Context, approval *Approval) error
	GetByExpenseID(ctx context.Context, expenseID int) (*Approval, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	GetByExpenseID(ctx context.Context, expenseID int) ([]*AuditLog, error)
}
