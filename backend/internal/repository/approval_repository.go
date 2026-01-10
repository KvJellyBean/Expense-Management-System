package repository

import (
	"context"
	"database/sql"
	"errors"
	"expense-management-system/internal/domain"
)

type approvalRepository struct {
	db *sql.DB
}

func NewApprovalRepository(db *sql.DB) domain.ApprovalRepository {
	return &approvalRepository{db: db}
}

func (r *approvalRepository) Create(ctx context.Context, approval *domain.Approval) error {
	query := `
		INSERT INTO approvals (expense_id, approver_id, status, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		approval.ExpenseID,
		approval.ApproverID,
		approval.Status,
		approval.Notes,
	).Scan(&approval.ID, &approval.CreatedAt)

	return err
}

func (r *approvalRepository) GetByExpenseID(ctx context.Context, expenseID int) (*domain.Approval, error) {
	query := `
		SELECT id, expense_id, approver_id, status, notes, created_at
		FROM approvals
		WHERE expense_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	approval := &domain.Approval{}
	err := r.db.QueryRowContext(ctx, query, expenseID).Scan(
		&approval.ID,
		&approval.ExpenseID,
		&approval.ApproverID,
		&approval.Status,
		&approval.Notes,
		&approval.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("approval not found")
	}

	return approval, err
}
