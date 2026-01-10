package usecase

import (
	"context"
	"errors"
	"expense-management-system/internal/domain"
	"expense-management-system/pkg/config"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecase_Login(t *testing.T) {
	// Pre-hash password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		email     string
		password  string
		setupMock func(*mockUserRepo)
		wantErr   bool
		wantRole  string
	}{
		{
			name:     "Valid employee login",
			email:    "employee1@example.com",
			password: "password123",
			setupMock: func(repo *mockUserRepo) {
				repo.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return &domain.User{
						ID:           1,
						Email:        "employee1@example.com",
						Name:         "Employee One",
						Role:         domain.RoleEmployee,
						PasswordHash: string(hashedPassword),
					}, nil
				}
			},
			wantErr:  false,
			wantRole: domain.RoleEmployee,
		},
		{
			name:     "Valid manager login",
			email:    "manager@example.com",
			password: "password123",
			setupMock: func(repo *mockUserRepo) {
				repo.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return &domain.User{
						ID:           3,
						Email:        "manager@example.com",
						Name:         "Manager Name",
						Role:         domain.RoleManager,
						PasswordHash: string(hashedPassword),
					}, nil
				}
			},
			wantErr:  false,
			wantRole: domain.RoleManager,
		},
		{
			name:     "Invalid password",
			email:    "employee1@example.com",
			password: "wrongpassword",
			setupMock: func(repo *mockUserRepo) {
				repo.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return &domain.User{
						ID:           1,
						Email:        "employee1@example.com",
						PasswordHash: string(hashedPassword),
					}, nil
				}
			},
			wantErr: true,
		},
		{
			name:     "Non-existent user",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMock: func(repo *mockUserRepo) {
				repo.getByEmailFunc = func(ctx context.Context, email string) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userRepo := &mockUserRepo{}

			if tt.setupMock != nil {
				tt.setupMock(userRepo)
			}

			cfg := &config.Config{
				JWTSecret:      "test-secret-key",
				JWTExpiryHours: 24,
			}

			uc := NewAuthUsecase(userRepo, cfg)

			token, user, err := uc.Login(ctx, tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("Expected token to be generated")
				}

				if user == nil {
					t.Error("Expected user to be returned")
					return
				}

				if user.Email != tt.email {
					t.Errorf("Login() email = %v, want %v", user.Email, tt.email)
				}

				if user.Role != tt.wantRole {
					t.Errorf("Login() role = %v, want %v", user.Role, tt.wantRole)
				}
			}
		})
	}
}

func TestAuthUsecase_ValidateToken(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		JWTSecret:      "test-secret-key",
		JWTExpiryHours: 24,
	}

	userRepo := &mockUserRepo{
		getByIDFunc: func(ctx context.Context, id int) (*domain.User, error) {
			return &domain.User{
				ID:    id,
				Email: "test@example.com",
				Role:  domain.RoleEmployee,
			}, nil
		},
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
			return &domain.User{
				ID:           1,
				Email:        email,
				Role:         domain.RoleEmployee,
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	uc := NewAuthUsecase(userRepo, cfg)

	// Generate valid token
	token, _, err := uc.Login(ctx, "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		setupRepo  func(*mockUserRepo)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "Valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:       "Invalid token format",
			token:      "invalid.token.string",
			wantErr:    true,
			wantErrMsg: "invalid token",
		},
		{
			name:       "Empty token",
			token:      "",
			wantErr:    true,
			wantErrMsg: "invalid token",
		},
		{
			name:       "Token with wrong signature",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoiZW1wbG95ZWUiLCJleHAiOjk5OTk5OTk5OTl9.invalid_signature",
			wantErr:    true,
			wantErrMsg: "invalid token",
		},
		{
			name:  "User not found in database",
			token: token,
			setupRepo: func(repo *mockUserRepo) {
				repo.getByIDFunc = func(ctx context.Context, id int) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr:    true,
			wantErrMsg: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRepo := &mockUserRepo{
				getByIDFunc:    userRepo.getByIDFunc,
				getByEmailFunc: userRepo.getByEmailFunc,
			}

			if tt.setupRepo != nil {
				tt.setupRepo(testRepo)
			}

			testUc := NewAuthUsecase(testRepo, cfg)
			user, err := testUc.ValidateToken(ctx, tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrMsg != "" {
				if err == nil || err.Error() != tt.wantErrMsg {
					t.Errorf("ValidateToken() error message = %v, want %v", err, tt.wantErrMsg)
				}
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("Expected user to be returned")
					return
				}

				if user.Email != "test@example.com" {
					t.Errorf("ValidateToken() email = %v, want test@example.com", user.Email)
				}
			}
		})
	}
}

func TestAuthUsecase_HashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf("HashPassword() unexpected error = %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	// Verify the hash can be compared
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Error("HashPassword() generated hash that doesn't match original password")
	}
}

func TestAuthUsecase_GenerateToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:      "test-secret-key",
		JWTExpiryHours: 24,
	}

	uc := &authUsecase{cfg: cfg}

	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleEmployee,
	}

	token, err := uc.generateToken(user)
	if err != nil {
		t.Errorf("generateToken() unexpected error = %v", err)
	}

	if token == "" {
		t.Error("generateToken() returned empty token")
	}
}
