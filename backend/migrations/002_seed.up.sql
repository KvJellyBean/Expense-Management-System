-- Seed users (password for all: password123)
INSERT INTO users (email, password_hash, name, role) VALUES
('employee1@example.com', '$2a$10$ArfoA5Y.NYwKkh/e61P5kutQB7u0zC2coCvmTD7qv9kwJ.GhgHZ1y', 'Kevin Natanael', 'employee'),
('employee2@example.com', '$2a$10$ArfoA5Y.NYwKkh/e61P5kutQB7u0zC2coCvmTD7qv9kwJ.GhgHZ1y', 'Natanael Kevin', 'employee'),
('manager@example.com', '$2a$10$ArfoA5Y.NYwKkh/e61P5kutQB7u0zC2coCvmTD7qv9kwJ.GhgHZ1y', 'Manager K', 'manager');

-- Seed expenses dengan berbagai status dan amounts
INSERT INTO expenses (user_id, amount_idr, description, receipt_url, status, auto_approved, submitted_at, payment_external_id) VALUES
-- Auto-approved expenses (< 1.000.000) - Employee 1
(1, 250000, 'Transportation to client meeting', '/mock-receipt.jpg', 'completed', true, '2026-01-05 10:00:00', 'exp-auto-1'),
(1, 500000, 'Office supplies purchase', '/mock-receipt.jpg', 'completed', true, '2026-01-06 14:30:00', 'exp-auto-2'),
-- Auto-approved expenses - Employee 2
(2, 750000, 'Team lunch at Warteg Bahari', '/mock-receipt.jpg', 'completed', true, '2026-01-07 12:00:00', 'exp-auto-3'),

-- Pending approval (>= 1.000.000) - Employee 1
(1, 1500000, 'Client meeting lunch at Plaza Indonesia', '/mock-receipt.jpg', 'awaiting_approval', false, '2026-01-08 09:00:00', 'exp-pending-1'),
-- Pending approval - Employee 2
(2, 2000000, 'Conference attendance fee', '/mock-receipt.jpg', 'awaiting_approval', false, '2026-01-08 10:30:00', 'exp-pending-2'),

-- Approved expenses (completed after manager approval) - Employee 1
(1, 3000000, 'Training workshop in Bali', '/mock-receipt.jpg', 'completed', false, '2026-01-03 08:00:00', 'exp-approved-1'),
-- Approved expenses - Employee 2
(2, 1200000, 'Marketing materials printing', '/mock-receipt.jpg', 'completed', false, '2026-01-04 11:00:00', 'exp-approved-2'),

-- Rejected expense - Employee 1
(1, 5000000, 'Luxury car rental for team outing', '/mock-receipt.jpg', 'rejected', false, '2026-01-02 16:00:00', 'exp-rejected-1');

-- Seed approvals for approved/rejected expenses
INSERT INTO approvals (expense_id, approver_id, status, notes, created_at) VALUES
(6, 3, 'approved', 'Approved for professional development', '2026-01-03 09:00:00'),
(7, 3, 'approved', 'Marketing budget approved', '2026-01-04 12:00:00'),
(8, 3, 'rejected', 'Excessive amount for team outing, please submit reasonable alternative', '2026-01-02 17:00:00');

-- Seed audit logs
INSERT INTO audit_logs (expense_id, user_id, action, old_status, new_status, metadata, created_at) VALUES
(1, 1, 'submit', NULL, 'approved', '{"amount_idr": 250000, "auto_approved": true}', '2026-01-05 10:00:00'),
(1, NULL, 'complete', 'approved', 'completed', '{"payment_id": "pay-1", "amount": 250000}', '2026-01-05 10:01:00'),
(2, 1, 'submit', NULL, 'approved', '{"amount_idr": 500000, "auto_approved": true}', '2026-01-06 14:30:00'),
(2, NULL, 'complete', 'approved', 'completed', '{"payment_id": "pay-2", "amount": 500000}', '2026-01-06 14:31:00'),
(3, 2, 'submit', NULL, 'approved', '{"amount_idr": 750000, "auto_approved": true}', '2026-01-07 12:00:00'),
(4, 1, 'submit', NULL, 'awaiting_approval', '{"amount_idr": 1500000, "auto_approved": false}', '2026-01-08 09:00:00'),
(5, 2, 'submit', NULL, 'awaiting_approval', '{"amount_idr": 2000000, "auto_approved": false}', '2026-01-08 10:30:00'),
(6, 1, 'submit', NULL, 'awaiting_approval', '{"amount_idr": 3000000, "auto_approved": false}', '2026-01-03 08:00:00'),
(6, 3, 'approve', 'awaiting_approval', 'approved', '{"notes": "Approved for professional development"}', '2026-01-03 09:00:00'),
(6, NULL, 'complete', 'approved', 'completed', '{"payment_id": "pay-3", "amount": 3000000}', '2026-01-03 09:05:00'),
(7, 2, 'submit', NULL, 'awaiting_approval', '{"amount_idr": 1200000, "auto_approved": false}', '2026-01-04 11:00:00'),
(7, 3, 'approve', 'awaiting_approval', 'approved', '{"notes": "Marketing budget approved"}', '2026-01-04 12:00:00'),
(8, 1, 'submit', NULL, 'awaiting_approval', '{"amount_idr": 5000000, "auto_approved": false}', '2026-01-02 16:00:00'),
(8, 3, 'reject', 'awaiting_approval', 'rejected', '{"notes": "Excessive amount for team outing"}', '2026-01-02 17:00:00');
