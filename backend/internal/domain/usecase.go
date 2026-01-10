package domain

import "context"

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, *User, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}

type ExpenseUsecase interface {
	Submit(ctx context.Context, userID int, amountIDR int, description string, receiptURL *string) (*Expense, error)
	GetByID(ctx context.Context, userID int, expenseID int, isManager bool) (*Expense, error)
	GetUserExpenses(ctx context.Context, userID int, status string, page, limit int, isManager bool) ([]*Expense, int, error)
	GetPendingApprovals(ctx context.Context, page, limit int) ([]*Expense, int, error)
	Approve(ctx context.Context, managerID, expenseID int, notes *string) error
	Reject(ctx context.Context, managerID, expenseID int, notes *string) error
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, expenseID int, amount int, externalID string) (string, error)
}
