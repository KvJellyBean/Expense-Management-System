package handler

import (
	"encoding/json"
	"expense-management-system/internal/domain"
	"expense-management-system/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ExpenseHandler struct {
	expenseUsecase domain.ExpenseUsecase
}

func NewExpenseHandler(expenseUsecase domain.ExpenseUsecase) *ExpenseHandler {
	return &ExpenseHandler{expenseUsecase: expenseUsecase}
}

type SubmitExpenseRequest struct {
	AmountIDR   int     `json:"amount_idr"`
	Description string  `json:"description"`
	ReceiptURL  *string `json:"receipt_url,omitempty"`
}

type SubmitExpenseResponse struct {
	ID               int     `json:"id"`
	AmountIDR        int     `json:"amount_idr"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
	RequiresApproval bool    `json:"requires_approval"`
	AutoApproved     bool    `json:"auto_approved"`
	CreatedAt        string  `json:"created_at"`
	ReceiptURL       *string `json:"receipt_url,omitempty"`
}

func (h *ExpenseHandler) Submit(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SubmitExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expense, err := h.expenseUsecase.Submit(r.Context(), user.ID, req.AmountIDR, req.Description, req.ReceiptURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := SubmitExpenseResponse{
		ID:               expense.ID,
		AmountIDR:        expense.AmountIDR,
		Description:      expense.Description,
		Status:           expense.Status,
		RequiresApproval: expense.AmountIDR >= domain.ApprovalThreshold,
		AutoApproved:     expense.AutoApproved,
		CreatedAt:        expense.SubmittedAt.Format("2006-01-02T15:04:05Z"),
		ReceiptURL:       expense.ReceiptURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

type ListExpensesResponse struct {
	Expenses []*domain.Expense `json:"expenses"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	isManager := user.Role == "manager"
	expenses, total, err := h.expenseUsecase.GetUserExpenses(r.Context(), user.ID, status, page, limit, isManager)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ListExpensesResponse{
		Expenses: expenses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	isManager := user.Role == domain.RoleManager
	expense, err := h.expenseUsecase.GetByID(r.Context(), user.ID, expenseID, isManager)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *ExpenseHandler) GetPendingApprovals(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	expenses, total, err := h.expenseUsecase.GetPendingApprovals(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ListExpensesResponse{
		Expenses: expenses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ApprovalRequest struct {
	Notes *string `json:"notes,omitempty"`
}

func (h *ExpenseHandler) Approve(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var req ApprovalRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.expenseUsecase.Approve(r.Context(), user.ID, expenseID, req.Notes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense approved successfully"})
}

func (h *ExpenseHandler) Reject(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var req ApprovalRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.expenseUsecase.Reject(r.Context(), user.ID, expenseID, req.Notes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense rejected successfully"})
}
