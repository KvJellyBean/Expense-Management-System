package repository

import (
	"context"
	"database/sql"
	"errors"
	"expense-management-system/internal/domain"
	"fmt"
	"strings"
)

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) domain.ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (user_id, amount_idr, description, receipt_url, status, auto_approved, payment_external_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, submitted_at, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		expense.UserID,
		expense.AmountIDR,
		expense.Description,
		expense.ReceiptURL,
		expense.Status,
		expense.AutoApproved,
		expense.PaymentExternalID,
	).Scan(&expense.ID, &expense.SubmittedAt, &expense.CreatedAt, &expense.UpdatedAt)

	return err
}

func (r *expenseRepository) GetByID(ctx context.Context, id int) (*domain.Expense, error) {
	query := `
		SELECT id, user_id, amount_idr, description, receipt_url, status, auto_approved,
		       submitted_at, processed_at, payment_id, payment_external_id, created_at, updated_at
		FROM expenses
		WHERE id = $1`

	expense := &domain.Expense{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.AmountIDR,
		&expense.Description,
		&expense.ReceiptURL,
		&expense.Status,
		&expense.AutoApproved,
		&expense.SubmittedAt,
		&expense.ProcessedAt,
		&expense.PaymentID,
		&expense.PaymentExternalID,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("expense not found")
	}

	return expense, err
}

func (r *expenseRepository) GetByUserID(ctx context.Context, userID int, status string, limit, offset int) ([]*domain.Expense, int, error) {
	var expenses []*domain.Expense
	var total int

	whereClause := "WHERE user_id = $1"
	args := []interface{}{userID}
	argCount := 1

	if status != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses %s", whereClause)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, amount_idr, description, receipt_url, status, auto_approved,
		       submitted_at, processed_at, payment_id, payment_external_id, created_at, updated_at
		FROM expenses
		%s
		ORDER BY submitted_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argCount+1, argCount+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.AmountIDR,
			&expense.Description,
			&expense.ReceiptURL,
			&expense.Status,
			&expense.AutoApproved,
			&expense.SubmittedAt,
			&expense.ProcessedAt,
			&expense.PaymentID,
			&expense.PaymentExternalID,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, total, nil
}

func (r *expenseRepository) GetAll(ctx context.Context, status string, limit, offset int) ([]*domain.Expense, int, error) {
	var expenses []*domain.Expense
	var total int

	whereClause := ""
	args := []interface{}{}
	argCount := 0

	if status != "" {
		argCount++
		whereClause = "WHERE status = $1"
		args = append(args, status)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses %s", whereClause)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, amount_idr, description, receipt_url, status, auto_approved,
		       submitted_at, processed_at, payment_id, payment_external_id, created_at, updated_at
		FROM expenses
		%s
		ORDER BY submitted_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argCount+1, argCount+2)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.AmountIDR,
			&expense.Description,
			&expense.ReceiptURL,
			&expense.Status,
			&expense.AutoApproved,
			&expense.SubmittedAt,
			&expense.ProcessedAt,
			&expense.PaymentID,
			&expense.PaymentExternalID,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, total, nil
}

func (r *expenseRepository) GetPendingApprovals(ctx context.Context, limit, offset int) ([]*domain.Expense, int, error) {
	var expenses []*domain.Expense
	var total int

	countQuery := "SELECT COUNT(*) FROM expenses WHERE status = $1 AND amount_idr >= $2"
	err := r.db.QueryRowContext(ctx, countQuery, domain.StatusAwaitingApproval, domain.ApprovalThreshold).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, user_id, amount_idr, description, receipt_url, status, auto_approved,
		       submitted_at, processed_at, payment_id, payment_external_id, created_at, updated_at
		FROM expenses
		WHERE status = $1 AND amount_idr >= $2
		ORDER BY submitted_at ASC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, query, domain.StatusAwaitingApproval, domain.ApprovalThreshold, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.AmountIDR,
			&expense.Description,
			&expense.ReceiptURL,
			&expense.Status,
			&expense.AutoApproved,
			&expense.SubmittedAt,
			&expense.ProcessedAt,
			&expense.PaymentID,
			&expense.PaymentExternalID,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, total, nil
}

func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	query := `
		UPDATE expenses
		SET status = $1, processed_at = $2, payment_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query,
		expense.Status,
		expense.ProcessedAt,
		expense.PaymentID,
		expense.ID,
	)

	return err
}

func (r *expenseRepository) UpdateStatus(ctx context.Context, id int, status string, processedAt *string) error {
	query := `
		UPDATE expenses
		SET status = $1, processed_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, status, processedAt, id)
	return err
}

func (r *expenseRepository) UpdatePaymentInfo(ctx context.Context, id int, paymentID, externalID string) error {
	setClauses := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{}
	argCount := 0

	if paymentID != "" {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("payment_id = $%d", argCount))
		args = append(args, paymentID)
	}

	if externalID != "" {
		argCount++
		setClauses = append(setClauses, fmt.Sprintf("payment_external_id = $%d", argCount))
		args = append(args, externalID)
	}

	argCount++
	args = append(args, id)

	query := fmt.Sprintf("UPDATE expenses SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argCount)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
