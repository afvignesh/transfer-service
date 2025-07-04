package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"transfer-service/model"
	svc "transfer-service/service"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// MockAccountRepository implements repository.AccountRepository for testing
type MockAccountRepository struct {
	accounts map[int]*model.Account
	createError error
	getError error
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		accounts: make(map[int]*model.Account),
	}
}

func (m *MockAccountRepository) Create(ctx context.Context, account model.Account) error {
	if m.createError != nil {
		return m.createError
	}
	if _, exists := m.accounts[account.ID]; exists {
		return &pq.Error{
			Code: "23505", // PostgreSQL unique violation error code
		}
	}
	m.accounts[account.ID] = &account
	return nil
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id int) (*model.Account, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	if account, exists := m.accounts[id]; exists {
		return account, nil
	}
	return nil, sql.ErrNoRows
}

func (m *MockAccountRepository) GetByIDWithLock(ctx context.Context, tx *sql.Tx, id int) (*model.Account, error) {
	return m.GetByID(ctx, id)
}

func (m *MockAccountRepository) UpdateBalance(ctx context.Context, id int, newBalance decimal.Decimal) error {
	return nil
}

func (m *MockAccountRepository) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int, newBalance decimal.Decimal) error {
	return nil
}

func (m *MockAccountRepository) DeleteByID(ctx context.Context, id int) error {
	return nil
}

func (m *MockAccountRepository) GetDB() *sql.DB {
	return nil
}

func TestCreateAccount_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockAccountRepository()
	service := svc.NewAccountService(mockRepo)
	ctx := context.Background()
	account := model.Account{ID: 1, Balance: decimal.NewFromFloat(1000.0)}

	// Act
	result := service.CreateAccount(ctx, account)

	// Assert
	if !result.Success {
		t.Errorf("Expected success, got failure: %s", result.Message)
	}
	if result.Status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, result.Status)
	}
	if result.Message != "Account created successfully" {
		t.Errorf("Expected message 'Account created successfully', got '%s'", result.Message)
	}
}

func TestCreateAccount_DuplicateID(t *testing.T) {
	// Arrange
	mockRepo := NewMockAccountRepository()
	service := svc.NewAccountService(mockRepo)
	ctx := context.Background()
	
	// Create first account
	account1 := model.Account{ID: 1, Balance: decimal.NewFromFloat(1000.0)}
	service.CreateAccount(ctx, account1)
	
	// Try to create second account with same ID
	account2 := model.Account{ID: 1, Balance: decimal.NewFromFloat(2000.0)}

	// Act
	result := service.CreateAccount(ctx, account2)

	// Assert
	if result.Success {
		t.Error("Expected failure for duplicate ID, got success")
	}
	if result.Status != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, result.Status)
	}
	if result.Message != "Account already exists" {
		t.Errorf("Expected message 'Account already exists', got '%s'", result.Message)
	}
}

func TestCreateAccount_InvalidRequest(t *testing.T) {
	// Arrange
	mockRepo := NewMockAccountRepository()
	mockRepo.createError = errors.New("invalid input syntax for integer")
	service := svc.NewAccountService(mockRepo)
	ctx := context.Background()
	account := model.Account{ID: 1, Balance: decimal.NewFromFloat(1000.0)}

	// Act
	result := service.CreateAccount(ctx, account)

	// Assert
	if result.Success {
		t.Error("Expected failure for invalid request, got success")
	}
	if result.Status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, result.Status)
	}
	if result.Message != "Failed to create account" {
		t.Errorf("Expected message 'Failed to create account', got '%s'", result.Message)
	}
}

func TestGetAccount_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockAccountRepository()
	service := svc.NewAccountService(mockRepo)
	ctx := context.Background()
	
	// Create account first
	account := model.Account{ID: 1, Balance: decimal.NewFromFloat(1000.0)}
	service.CreateAccount(ctx, account)

	// Act
	result := service.GetAccount(ctx, 1)

	// Assert
	if !result.Success {
		t.Errorf("Expected success, got failure: %s", result.Message)
	}
	if result.Status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, result.Status)
	}
	if result.Message != "Account retrieved successfully" {
		t.Errorf("Expected message 'Account retrieved successfully', got '%s'", result.Message)
	}
	
	// Check if returned data is correct
	if result.Data == nil {
		t.Error("Expected account data, got nil")
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockAccountRepository()
	service := svc.NewAccountService(mockRepo)
	ctx := context.Background()

	// Act
	result := service.GetAccount(ctx, 999)

	// Assert
	if result.Success {
		t.Error("Expected failure for non-existent account, got success")
	}
	if result.Status != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, result.Status)
	}
	if result.Message != "Account not found" {
		t.Errorf("Expected message 'Account not found', got '%s'", result.Message)
	}
} 