package usecase

import (
	"context"
	"expense-management-system/internal/domain"
	"expense-management-system/pkg/logger"
	"testing"
)

func init() {
	// Initialize logger for tests
	logger.Init()
}

// Mock Repositories
type mockExpenseRepo struct {
	createFunc          func(ctx context.Context, expense *domain.Expense) error
	getByIDFunc         func(ctx context.Context, id int) (*domain.Expense, error)
	updateFunc          func(ctx context.Context, expense *domain.Expense) error
	getByUserIDFunc     func(ctx context.Context, userID int, status string, limit, offset int) ([]*domain.Expense, int, error)
	getAllFunc          func(ctx context.Context, status string, limit, offset int) ([]*domain.Expense, int, error)
	getPendingApprovals func(ctx context.Context, limit, offset int) ([]*domain.Expense, int, error)
}

func (m *mockExpenseRepo) Create(ctx context.Context, expense *domain.Expense) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, expense)
	}
	expense.ID = 1
	return nil
}

func (m *mockExpenseRepo) GetByID(ctx context.Context, id int) (*domain.Expense, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockExpenseRepo) Update(ctx context.Context, expense *domain.Expense) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, expense)
	}
	return nil
}

func (m *mockExpenseRepo) UpdateStatus(ctx context.Context, id int, status string, processedAt *string) error {
	return nil
}

func (m *mockExpenseRepo) UpdatePaymentInfo(ctx context.Context, id int, paymentID, externalID string) error {
	return nil
}

func (m *mockExpenseRepo) GetByUserID(ctx context.Context, userID int, status string, limit, offset int) ([]*domain.Expense, int, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(ctx, userID, status, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockExpenseRepo) GetAll(ctx context.Context, status string, limit, offset int) ([]*domain.Expense, int, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, status, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockExpenseRepo) GetPendingApprovals(ctx context.Context, limit, offset int) ([]*domain.Expense, int, error) {
	if m.getPendingApprovals != nil {
		return m.getPendingApprovals(ctx, limit, offset)
	}
	return nil, 0, nil
}

type mockApprovalRepo struct {
	getByExpenseIDFunc func(ctx context.Context, expenseID int) (*domain.Approval, error)
}

func (m *mockApprovalRepo) Create(ctx context.Context, approval *domain.Approval) error {
	approval.ID = 1
	return nil
}

func (m *mockApprovalRepo) GetByExpenseID(ctx context.Context, expenseID int) (*domain.Approval, error) {
	if m.getByExpenseIDFunc != nil {
		return m.getByExpenseIDFunc(ctx, expenseID)
	}
	return nil, nil
}

type mockAuditRepo struct {
	createFunc func(ctx context.Context, log *domain.AuditLog) error
}

func (m *mockAuditRepo) Create(ctx context.Context, log *domain.AuditLog) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, log)
	}
	log.ID = 1
	return nil
}

func (m *mockAuditRepo) GetByExpenseID(ctx context.Context, expenseID int) ([]*domain.AuditLog, error) {
	return nil, nil
}

type mockUserRepo struct {
	getByIDFunc    func(ctx context.Context, id int) (*domain.User, error)
	getByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return &domain.User{ID: id, Email: "test@example.com"}, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	return nil
}

// Tests
func TestExpenseUsecase_Submit(t *testing.T) {
	tests := []struct {
		name            string
		userID          int
		amountIDR       int
		description     string
		receiptURL      *string
		wantErr         bool
		wantStatus      string
		wantAutoApprove bool
	}{
		{
			name:            "Valid expense below threshold (auto-approved)",
			userID:          1,
			amountIDR:       500000,
			description:     "Office supplies",
			receiptURL:      strPtr("https://example.com/receipt.jpg"),
			wantErr:         false,
			wantStatus:      domain.StatusApproved,
			wantAutoApprove: true,
		},
		{
			name:            "Valid expense at threshold requires approval",
			userID:          1,
			amountIDR:       1000000,
			description:     "Client meeting",
			receiptURL:      nil,
			wantErr:         false,
			wantStatus:      domain.StatusAwaitingApproval,
			wantAutoApprove: false,
		},
		{
			name:            "Valid expense above threshold requires approval",
			userID:          1,
			amountIDR:       2000000,
			description:     "Team building event",
			receiptURL:      strPtr("/mock-receipt.pdf"),
			wantErr:         false,
			wantStatus:      domain.StatusAwaitingApproval,
			wantAutoApprove: false,
		},
		{
			name:        "Amount below minimum",
			userID:      1,
			amountIDR:   5000,
			description: "Small purchase",
			wantErr:     true,
		},
		{
			name:        "Amount above maximum",
			userID:      1,
			amountIDR:   60000000,
			description: "Expensive item",
			wantErr:     true,
		},
		{
			name:        "Empty description",
			userID:      1,
			amountIDR:   100000,
			description: "",
			wantErr:     true,
		},
		{
			name:            "Minimum valid amount",
			userID:          1,
			amountIDR:       10000,
			description:     "Minimum expense",
			wantErr:         false,
			wantStatus:      domain.StatusApproved,
			wantAutoApprove: true,
		},
		{
			name:            "Maximum valid amount",
			userID:          1,
			amountIDR:       50000000,
			description:     "Maximum expense",
			wantErr:         false,
			wantStatus:      domain.StatusAwaitingApproval,
			wantAutoApprove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			paymentChan := make(chan PaymentJob, 10)

			expenseRepo := &mockExpenseRepo{}
			approvalRepo := &mockApprovalRepo{}
			auditRepo := &mockAuditRepo{}
			userRepo := &mockUserRepo{}

			uc := NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

			expense, err := uc.Submit(ctx, tt.userID, tt.amountIDR, tt.description, tt.receiptURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("Submit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if expense == nil {
					t.Error("Expected expense to be created, got nil")
					return
				}

				if expense.Status != tt.wantStatus {
					t.Errorf("Submit() status = %v, want %v", expense.Status, tt.wantStatus)
				}

				if expense.AutoApproved != tt.wantAutoApprove {
					t.Errorf("Submit() autoApproved = %v, want %v", expense.AutoApproved, tt.wantAutoApprove)
				}

				if expense.AmountIDR != tt.amountIDR {
					t.Errorf("Submit() amountIDR = %v, want %v", expense.AmountIDR, tt.amountIDR)
				}

				if expense.Description != tt.description {
					t.Errorf("Submit() description = %v, want %v", expense.Description, tt.description)
				}

				if expense.PaymentExternalID == nil {
					t.Error("Expected PaymentExternalID to be set")
				}

				// Check payment queue for auto-approved expenses
				if tt.wantAutoApprove {
					select {
					case job := <-paymentChan:
						if job.Amount != tt.amountIDR {
							t.Errorf("Payment job amount = %v, want %v", job.Amount, tt.amountIDR)
						}
					default:
						t.Error("Expected payment job to be queued for auto-approved expense")
					}
				}
			}
		})
	}
}

func TestExpenseUsecase_Approve(t *testing.T) {
	tests := []struct {
		name       string
		expenseID  int
		approverID int
		notes      string
		setupMock  func(*mockExpenseRepo, *mockApprovalRepo, *mockAuditRepo)
		wantErr    bool
	}{
		{
			name:       "Successfully approve pending expense",
			expenseID:  1,
			approverID: 3,
			notes:      "Approved for Q1 budget",
			setupMock: func(expRepo *mockExpenseRepo, apprRepo *mockApprovalRepo, auditRepo *mockAuditRepo) {
				expRepo.getByIDFunc = func(ctx context.Context, id int) (*domain.Expense, error) {
					return &domain.Expense{
						ID:                1,
						UserID:            1,
						AmountIDR:         1500000,
						Status:            domain.StatusAwaitingApproval,
						PaymentExternalID: strPtr("test-external-id"),
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name:       "Cannot approve already approved expense",
			expenseID:  1,
			approverID: 3,
			notes:      "Approving again",
			setupMock: func(expRepo *mockExpenseRepo, apprRepo *mockApprovalRepo, auditRepo *mockAuditRepo) {
				expRepo.getByIDFunc = func(ctx context.Context, id int) (*domain.Expense, error) {
					return &domain.Expense{
						ID:        1,
						UserID:    1,
						AmountIDR: 1500000,
						Status:    domain.StatusApproved,
					}, nil
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			paymentChan := make(chan PaymentJob, 10)

			expenseRepo := &mockExpenseRepo{}
			approvalRepo := &mockApprovalRepo{}
			auditRepo := &mockAuditRepo{}
			userRepo := &mockUserRepo{}

			if tt.setupMock != nil {
				tt.setupMock(expenseRepo, approvalRepo, auditRepo)
			}

			uc := NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

			err := uc.Approve(ctx, tt.approverID, tt.expenseID, strPtr(tt.notes))

			if (err != nil) != tt.wantErr {
				t.Errorf("Approve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check payment queue
				select {
				case job := <-paymentChan:
					if job.ExpenseID != tt.expenseID {
						t.Errorf("Payment job expenseID = %v, want %v", job.ExpenseID, tt.expenseID)
					}
				default:
					t.Error("Expected payment job to be queued")
				}
			}
		})
	}
}

func TestExpenseUsecase_Reject(t *testing.T) {
	ctx := context.Background()
	paymentChan := make(chan PaymentJob, 10)

	expenseRepo := &mockExpenseRepo{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Expense, error) {
			return &domain.Expense{
				ID:        1,
				UserID:    1,
				AmountIDR: 1500000,
				Status:    domain.StatusAwaitingApproval,
			}, nil
		},
	}
	approvalRepo := &mockApprovalRepo{}
	auditRepo := &mockAuditRepo{}
	userRepo := &mockUserRepo{}

	uc := NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

	err := uc.Reject(ctx, 3, 1, strPtr("Receipt not clear"))
	if err != nil {
		t.Errorf("Reject() unexpected error = %v", err)
	}

	// Ensure no payment job is queued for rejected expense
	select {
	case <-paymentChan:
		t.Error("Payment job should not be queued for rejected expense")
	default:
		// Expected
	}
}

func TestExpenseUsecase_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		userID     int
		expenseID  int
		isManager  bool
		setupMock  func(*mockExpenseRepo, *mockApprovalRepo)
		wantErr    bool
		wantUserID int
	}{
		{
			name:      "Employee can access own expense",
			userID:    1,
			expenseID: 1,
			isManager: false,
			setupMock: func(expRepo *mockExpenseRepo, apprRepo *mockApprovalRepo) {
				expRepo.getByIDFunc = func(ctx context.Context, id int) (*domain.Expense, error) {
					return &domain.Expense{
						ID:        1,
						UserID:    1,
						AmountIDR: 500000,
						Status:    domain.StatusApproved,
					}, nil
				}
				apprRepo.getByExpenseIDFunc = func(ctx context.Context, expenseID int) (*domain.Approval, error) {
					return &domain.Approval{
						ID:         1,
						ExpenseID:  1,
						ApproverID: 3,
					}, nil
				}
			},
			wantErr:    false,
			wantUserID: 1,
		},
		{
			name:      "Employee cannot access other's expense",
			userID:    1,
			expenseID: 2,
			isManager: false,
			setupMock: func(expRepo *mockExpenseRepo, apprRepo *mockApprovalRepo) {
				expRepo.getByIDFunc = func(ctx context.Context, id int) (*domain.Expense, error) {
					return &domain.Expense{
						ID:        2,
						UserID:    2,
						AmountIDR: 500000,
						Status:    domain.StatusApproved,
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name:      "Manager can access any expense",
			userID:    3,
			expenseID: 1,
			isManager: true,
			setupMock: func(expRepo *mockExpenseRepo, apprRepo *mockApprovalRepo) {
				expRepo.getByIDFunc = func(ctx context.Context, id int) (*domain.Expense, error) {
					return &domain.Expense{
						ID:        1,
						UserID:    1,
						AmountIDR: 500000,
						Status:    domain.StatusCompleted,
					}, nil
				}
				apprRepo.getByExpenseIDFunc = func(ctx context.Context, expenseID int) (*domain.Approval, error) {
					return &domain.Approval{
						ID:         1,
						ExpenseID:  1,
						ApproverID: 3,
					}, nil
				}
			},
			wantErr:    false,
			wantUserID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			paymentChan := make(chan PaymentJob, 10)

			expenseRepo := &mockExpenseRepo{}
			approvalRepo := &mockApprovalRepo{}
			auditRepo := &mockAuditRepo{}
			userRepo := &mockUserRepo{}

			if tt.setupMock != nil {
				tt.setupMock(expenseRepo, approvalRepo)
			}

			uc := NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

			expense, err := uc.GetByID(ctx, tt.userID, tt.expenseID, tt.isManager)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if expense == nil {
					t.Error("Expected expense to be returned")
					return
				}

				if expense.UserID != tt.wantUserID {
					t.Errorf("GetByID() userID = %v, want %v", expense.UserID, tt.wantUserID)
				}
			}
		})
	}
}

func TestExpenseUsecase_GetUserExpenses(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		status    string
		page      int
		limit     int
		isManager bool
		setupMock func(*mockExpenseRepo)
		wantErr   bool
		wantCount int
	}{
		{
			name:      "Employee gets own expenses",
			userID:    1,
			status:    "",
			page:      1,
			limit:     20,
			isManager: false,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getByUserIDFunc = func(ctx context.Context, userID int, status string, limit, offset int) ([]*domain.Expense, int, error) {
					return []*domain.Expense{
						{ID: 1, UserID: 1, AmountIDR: 500000},
						{ID: 2, UserID: 1, AmountIDR: 750000},
					}, 2, nil
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:      "Manager gets all expenses",
			userID:    3,
			status:    domain.StatusApproved,
			page:      1,
			limit:     20,
			isManager: true,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getAllFunc = func(ctx context.Context, status string, limit, offset int) ([]*domain.Expense, int, error) {
					return []*domain.Expense{
						{ID: 1, UserID: 1, AmountIDR: 500000, Status: domain.StatusApproved},
						{ID: 2, UserID: 2, AmountIDR: 750000, Status: domain.StatusApproved},
						{ID: 3, UserID: 1, AmountIDR: 900000, Status: domain.StatusApproved},
					}, 3, nil
				}
			},
			wantErr:   false,
			wantCount: 3,
		},
		{
			name:      "Pagination with page < 1 defaults to 1",
			userID:    1,
			status:    "",
			page:      0,
			limit:     20,
			isManager: false,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getByUserIDFunc = func(ctx context.Context, userID int, status string, limit, offset int) ([]*domain.Expense, int, error) {
					if offset != 0 {
						t.Error("Expected offset to be 0 for page 1")
					}
					return []*domain.Expense{{ID: 1}}, 1, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:      "Limit > 100 defaults to 20",
			userID:    1,
			status:    "",
			page:      1,
			limit:     150,
			isManager: false,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getByUserIDFunc = func(ctx context.Context, userID int, status string, limit, offset int) ([]*domain.Expense, int, error) {
					if limit != 20 {
						t.Errorf("Expected limit to be 20, got %d", limit)
					}
					return []*domain.Expense{{ID: 1}}, 1, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			paymentChan := make(chan PaymentJob, 10)

			expenseRepo := &mockExpenseRepo{}
			approvalRepo := &mockApprovalRepo{}
			auditRepo := &mockAuditRepo{}
			userRepo := &mockUserRepo{}

			if tt.setupMock != nil {
				tt.setupMock(expenseRepo)
			}

			uc := NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

			expenses, count, err := uc.GetUserExpenses(ctx, tt.userID, tt.status, tt.page, tt.limit, tt.isManager)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserExpenses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(expenses) != tt.wantCount {
					t.Errorf("GetUserExpenses() count = %v, want %v", len(expenses), tt.wantCount)
				}

				if count != tt.wantCount {
					t.Errorf("GetUserExpenses() total = %v, want %v", count, tt.wantCount)
				}
			}
		})
	}
}

func TestExpenseUsecase_GetPendingApprovals(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		limit     int
		setupMock func(*mockExpenseRepo)
		wantCount int
	}{
		{
			name:  "Get pending approvals with default pagination",
			page:  1,
			limit: 20,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getPendingApprovals = func(ctx context.Context, limit, offset int) ([]*domain.Expense, int, error) {
					return []*domain.Expense{
						{ID: 1, Status: domain.StatusAwaitingApproval, AmountIDR: 1500000},
						{ID: 2, Status: domain.StatusAwaitingApproval, AmountIDR: 2000000},
					}, 2, nil
				}
			},
			wantCount: 2,
		},
		{
			name:  "Invalid page defaults to 1",
			page:  -1,
			limit: 20,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getPendingApprovals = func(ctx context.Context, limit, offset int) ([]*domain.Expense, int, error) {
					if offset != 0 {
						t.Error("Expected offset to be 0 for page 1")
					}
					return []*domain.Expense{{ID: 1}}, 1, nil
				}
			},
			wantCount: 1,
		},
		{
			name:  "Invalid limit defaults to 20",
			page:  1,
			limit: 0,
			setupMock: func(repo *mockExpenseRepo) {
				repo.getPendingApprovals = func(ctx context.Context, limit, offset int) ([]*domain.Expense, int, error) {
					if limit != 20 {
						t.Errorf("Expected limit to be 20, got %d", limit)
					}
					return []*domain.Expense{{ID: 1}}, 1, nil
				}
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			paymentChan := make(chan PaymentJob, 10)

			expenseRepo := &mockExpenseRepo{}
			approvalRepo := &mockApprovalRepo{}
			auditRepo := &mockAuditRepo{}
			userRepo := &mockUserRepo{}

			if tt.setupMock != nil {
				tt.setupMock(expenseRepo)
			}

			uc := NewExpenseUsecase(expenseRepo, approvalRepo, auditRepo, userRepo, paymentChan)

			expenses, count, err := uc.GetPendingApprovals(ctx, tt.page, tt.limit)
			if err != nil {
				t.Errorf("GetPendingApprovals() unexpected error = %v", err)
				return
			}

			if len(expenses) != tt.wantCount {
				t.Errorf("GetPendingApprovals() count = %v, want %v", len(expenses), tt.wantCount)
			}

			if count != tt.wantCount {
				t.Errorf("GetPendingApprovals() total = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
