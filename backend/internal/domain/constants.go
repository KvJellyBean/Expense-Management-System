package domain

const (
	MinExpenseAmount  = 10000
	MaxExpenseAmount  = 50000000
	ApprovalThreshold = 1000000
)

const (
	RoleEmployee = "employee"
	RoleManager  = "manager"
)

const (
	StatusAwaitingApproval = "awaiting_approval"
	StatusApproved         = "approved"
	StatusRejected         = "rejected"
	StatusCompleted        = "completed"
)

const (
	ActionSubmit   = "submit"
	ActionApprove  = "approve"
	ActionReject   = "reject"
	ActionComplete = "complete"
)
