package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"expense-management-system/internal/domain"
	"expense-management-system/internal/usecase"
	"expense-management-system/pkg/config"
	"expense-management-system/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PaymentService struct {
	client      *http.Client
	cfg         *config.Config
	expenseRepo domain.ExpenseRepository
	auditRepo   domain.AuditLogRepository
}

type PaymentRequest struct {
	Amount     int    `json:"amount"`
	ExternalID string `json:"external_id"`
}

type PaymentResponse struct {
	Data struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

func NewPaymentService(cfg *config.Config, expenseRepo domain.ExpenseRepository, auditRepo domain.AuditLogRepository) *PaymentService {
	return &PaymentService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cfg:         cfg,
		expenseRepo: expenseRepo,
		auditRepo:   auditRepo,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, expenseID int, amount int, externalID string) (string, error) {
	reqBody := PaymentRequest{
		Amount:     amount,
		ExternalID: externalID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/v1/payments", s.cfg.PaymentAPIURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return "", err
	}

	if resp.StatusCode == 400 && paymentResp.Message == "external id already exists" {
		return "", errors.New("idempotency_error: external_id already exists")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("payment failed with status %d", resp.StatusCode)
	}

	if paymentResp.Data.Status != "success" {
		return "", fmt.Errorf("payment status: %s", paymentResp.Data.Status)
	}

	return paymentResp.Data.ID, nil
}

func (s *PaymentService) ProcessPaymentWithRetry(ctx context.Context, job usecase.PaymentJob, maxRetries int) error {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		paymentID, err := s.ProcessPayment(ctx, job.ExpenseID, job.Amount, job.ExternalID)

		if err == nil {
			if err := s.expenseRepo.UpdatePaymentInfo(ctx, job.ExpenseID, paymentID, job.ExternalID); err != nil {
				logger.ErrorLogger.Printf("Failed to update payment info for expense %d: %v", job.ExpenseID, err)
				return err
			}

			now := time.Now().Format(time.RFC3339)
			if err := s.expenseRepo.UpdateStatus(ctx, job.ExpenseID, domain.StatusCompleted, &now); err != nil {
				logger.ErrorLogger.Printf("Failed to update status for expense %d: %v", job.ExpenseID, err)
				return err
			}

			newStatus := domain.StatusCompleted
			auditLog := &domain.AuditLog{
				ExpenseID: job.ExpenseID,
				Action:    domain.ActionComplete,
				NewStatus: &newStatus,
				Metadata: map[string]interface{}{
					"payment_id":  paymentID,
					"external_id": job.ExternalID,
					"amount":      job.Amount,
				},
			}
			s.auditRepo.Create(ctx, auditLog)

			logger.InfoLogger.Printf("Payment successful for expense %d, payment_id: %s", job.ExpenseID, paymentID)
			return nil
		}

		lastErr = err

		if err.Error() == "idempotency_error: external_id already exists" {
			logger.InfoLogger.Printf("Expense %d already processed (idempotency check), marking as completed", job.ExpenseID)

			now := time.Now().Format(time.RFC3339)
			s.expenseRepo.UpdateStatus(ctx, job.ExpenseID, domain.StatusCompleted, &now)
			return nil
		}

		if attempt < maxRetries {
			backoff := time.Duration(attempt*2) * time.Second
			logger.InfoLogger.Printf("Payment attempt %d/%d failed for expense %d, retrying in %v: %v",
				attempt, maxRetries, job.ExpenseID, backoff, err)
			time.Sleep(backoff)
		}
	}

	logger.ErrorLogger.Printf("Payment failed for expense %d after %d attempts: %v", job.ExpenseID, maxRetries, lastErr)
	return lastErr
}
