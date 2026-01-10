package usecase

import (
	"context"
	"errors"
	"expense-management-system/internal/domain"
	"expense-management-system/pkg/logger"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type expenseUsecase struct {
	expenseRepo  domain.ExpenseRepository
	approvalRepo domain.ApprovalRepository
	auditRepo    domain.AuditLogRepository
	userRepo     domain.UserRepository
	paymentChan  chan PaymentJob
}

type PaymentJob struct {
	ExpenseID  int
	Amount     int
	ExternalID string
}

func NewExpenseUsecase(
	expenseRepo domain.ExpenseRepository,
	approvalRepo domain.ApprovalRepository,
	auditRepo domain.AuditLogRepository,
	userRepo domain.UserRepository,
	paymentChan chan PaymentJob,
) domain.ExpenseUsecase {
	return &expenseUsecase{
		expenseRepo:  expenseRepo,
		approvalRepo: approvalRepo,
		auditRepo:    auditRepo,
		userRepo:     userRepo,
		paymentChan:  paymentChan,
	}
}

func (u *expenseUsecase) Submit(ctx context.Context, userID int, amountIDR int, description string, receiptURL *string) (*domain.Expense, error) {
	if amountIDR < domain.MinExpenseAmount || amountIDR > domain.MaxExpenseAmount {
		return nil, fmt.Errorf("amount must be between IDR %d and IDR %d", domain.MinExpenseAmount, domain.MaxExpenseAmount)
	}

	if description == "" {
		return nil, errors.New("description is required")
	}

	externalID := uuid.New().String()
	autoApproved := amountIDR < domain.ApprovalThreshold
	status := domain.StatusAwaitingApproval
	if autoApproved {
		status = domain.StatusApproved
	}

	expense := &domain.Expense{
		UserID:            userID,
		AmountIDR:         amountIDR,
		Description:       description,
		ReceiptURL:        receiptURL,
		Status:            status,
		AutoApproved:      autoApproved,
		PaymentExternalID: &externalID,
	}

	if err := u.expenseRepo.Create(ctx, expense); err != nil {
		return nil, err
	}

	auditLog := &domain.AuditLog{
		ExpenseID: expense.ID,
		UserID:    &userID,
		Action:    domain.ActionSubmit,
		NewStatus: &status,
		Metadata: map[string]interface{}{
			"amount_idr":    amountIDR,
			"auto_approved": autoApproved,
		},
	}
	u.auditRepo.Create(ctx, auditLog)

	if autoApproved {
		logger.InfoLogger.Printf("Auto-approved expense %d, sending to payment queue", expense.ID)
		u.sendToPaymentQueue(expense.ID, amountIDR, externalID)

		user, _ := u.userRepo.GetByID(ctx, userID)
		if user != nil {
			logger.InfoLogger.Printf("[EMAIL] Auto-approval notification sent to %s for expense %d (IDR %d)", user.Email, expense.ID, amountIDR)
		}
	} else {
		logger.InfoLogger.Printf("Expense %d requires manager approval (amount: IDR %d >= threshold)", expense.ID, amountIDR)

		user, _ := u.userRepo.GetByID(ctx, userID)
		if user != nil {
			logger.InfoLogger.Printf("[EMAIL] Approval request notification sent to managers for expense %d by %s", expense.ID, user.Email)
		}
	}

	return expense, nil
}

func (u *expenseUsecase) GetByID(ctx context.Context, userID int, expenseID int, isManager bool) (*domain.Expense, error) {
	expense, err := u.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return nil, err
	}

	if !isManager && expense.UserID != userID {
		return nil, errors.New("unauthorized access to expense")
	}

	// Populate approval data if expense has been approved/rejected/completed
	// Completed expenses were previously approved and then payment processed
	if expense.Status == domain.StatusApproved || expense.Status == domain.StatusRejected || expense.Status == domain.StatusCompleted {
		approval, err := u.approvalRepo.GetByExpenseID(ctx, expenseID)
		if err == nil {
			expense.Approval = approval
		}
	}

	return expense, nil
}

func (u *expenseUsecase) GetUserExpenses(ctx context.Context, userID int, status string, page, limit int, isManager bool) ([]*domain.Expense, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Managers can see all expenses
	if isManager {
		return u.expenseRepo.GetAll(ctx, status, limit, offset)
	}

	// Regular users see only their own expenses
	return u.expenseRepo.GetByUserID(ctx, userID, status, limit, offset)
}

func (u *expenseUsecase) GetPendingApprovals(ctx context.Context, page, limit int) ([]*domain.Expense, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return u.expenseRepo.GetPendingApprovals(ctx, limit, offset)
}

func (u *expenseUsecase) Approve(ctx context.Context, managerID, expenseID int, notes *string) error {
	expense, err := u.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return err
	}

	if expense.Status != domain.StatusAwaitingApproval {
		return errors.New("expense is not awaiting approval")
	}

	approval := &domain.Approval{
		ExpenseID:  expenseID,
		ApproverID: managerID,
		Status:     domain.StatusApproved,
		Notes:      notes,
	}

	if err := u.approvalRepo.Create(ctx, approval); err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	if err := u.expenseRepo.UpdateStatus(ctx, expenseID, domain.StatusApproved, &now); err != nil {
		return err
	}

	oldStatus := expense.Status
	newStatus := domain.StatusApproved
	auditLog := &domain.AuditLog{
		ExpenseID: expenseID,
		UserID:    &managerID,
		Action:    domain.ActionApprove,
		OldStatus: &oldStatus,
		NewStatus: &newStatus,
		Metadata: map[string]interface{}{
			"notes": notes,
		},
	}
	u.auditRepo.Create(ctx, auditLog)

	logger.InfoLogger.Printf("Expense %d approved by manager %d, sending to payment queue", expenseID, managerID)
	u.sendToPaymentQueue(expenseID, expense.AmountIDR, *expense.PaymentExternalID)

	user, _ := u.userRepo.GetByID(ctx, expense.UserID)
	if user != nil {
		logger.InfoLogger.Printf("[EMAIL] Approval notification sent to %s for expense %d", user.Email, expenseID)
	}

	return nil
}

func (u *expenseUsecase) Reject(ctx context.Context, managerID, expenseID int, notes *string) error {
	expense, err := u.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		return err
	}

	if expense.Status != domain.StatusAwaitingApproval {
		return errors.New("expense is not awaiting approval")
	}

	approval := &domain.Approval{
		ExpenseID:  expenseID,
		ApproverID: managerID,
		Status:     domain.StatusRejected,
		Notes:      notes,
	}

	if err := u.approvalRepo.Create(ctx, approval); err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	if err := u.expenseRepo.UpdateStatus(ctx, expenseID, domain.StatusRejected, &now); err != nil {
		return err
	}

	oldStatus := expense.Status
	newStatus := domain.StatusRejected
	auditLog := &domain.AuditLog{
		ExpenseID: expenseID,
		UserID:    &managerID,
		Action:    domain.ActionReject,
		OldStatus: &oldStatus,
		NewStatus: &newStatus,
		Metadata: map[string]interface{}{
			"notes": notes,
		},
	}
	u.auditRepo.Create(ctx, auditLog)

	logger.InfoLogger.Printf("Expense %d rejected by manager %d", expenseID, managerID)

	user, _ := u.userRepo.GetByID(ctx, expense.UserID)
	if user != nil {
		logger.InfoLogger.Printf("[EMAIL] Rejection notification sent to %s for expense %d", user.Email, expenseID)
	}

	return nil
}

func (u *expenseUsecase) sendToPaymentQueue(expenseID, amount int, externalID string) {
	select {
	case u.paymentChan <- PaymentJob{
		ExpenseID:  expenseID,
		Amount:     amount,
		ExternalID: externalID,
	}:
		logger.InfoLogger.Printf("Payment job queued for expense %d", expenseID)
	default:
		logger.ErrorLogger.Printf("Payment queue full, could not queue expense %d", expenseID)
	}
}
