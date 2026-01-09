package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"expense-management-system/internal/domain"
)

type auditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) domain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	var metadataJSON []byte
	var err error

	if log.Metadata != nil {
		metadataJSON, err = json.Marshal(log.Metadata)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO audit_logs (expense_id, user_id, action, old_status, new_status, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err = r.db.QueryRowContext(ctx, query,
		log.ExpenseID,
		log.UserID,
		log.Action,
		log.OldStatus,
		log.NewStatus,
		metadataJSON,
	).Scan(&log.ID, &log.CreatedAt)

	return err
}

func (r *auditLogRepository) GetByExpenseID(ctx context.Context, expenseID int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, expense_id, user_id, action, old_status, new_status, metadata, created_at
		FROM audit_logs
		WHERE expense_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog

	for rows.Next() {
		log := &domain.AuditLog{}
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.ExpenseID,
			&log.UserID,
			&log.Action,
			&log.OldStatus,
			&log.NewStatus,
			&metadataJSON,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			err = json.Unmarshal(metadataJSON, &log.Metadata)
			if err != nil {
				return nil, err
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}
