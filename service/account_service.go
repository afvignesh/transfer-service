package service

import (
    "context"
    "transfer-service/model"
    "transfer-service/repository"
    "errors"
    "net/http"
    "transfer-service/middleware"
    "go.uber.org/zap"
    "github.com/shopspring/decimal"
)

type AccountService struct {
    repo repository.AccountRepository
}

var ErrAccountExists = errors.New("account already exists")

func NewAccountService(repo repository.AccountRepository) *AccountService {
    return &AccountService{repo: repo}
}

// AccountResult represents the result of an account operation
type AccountResult struct {
    Success bool
    Status  int
    Message string
    Error   string
    Data    interface{}
}

func (s *AccountService) CreateAccount(ctx context.Context, acc model.Account) *AccountResult {
    log := middleware.GetLogger()
    
    // Validate balance precision (5 decimal places)
    if !isValidPrecision(acc.Balance, 5) {
        log.Warn("Account creation failed - invalid balance precision",
            zap.Int("account_id", acc.ID),
            zap.String("balance", acc.Balance.String()),
        )
        return &AccountResult{
            Success: false,
            Status:  http.StatusBadRequest,
            Message: "Balance must have at most 5 decimal places",
            Error:   "invalid precision",
        }
    }
    
    log.Info("Creating account",
        zap.Int("account_id", acc.ID),
        zap.Float64("balance", formatDecimal(acc.Balance)),
    )
    
    err := s.repo.Create(ctx, acc)
    if err != nil {
        if middleware.IsUniqueViolation(err) {
            log.Warn("Account creation failed - duplicate ID",
                zap.Int("account_id", acc.ID),
                zap.Error(err),
            )
            return &AccountResult{
                Success: false,
                Status:  http.StatusConflict,
                Message: "Account already exists",
                Error:   err.Error(),
            }
        }
        log.Error("Account creation failed",
            zap.Int("account_id", acc.ID),
            zap.Error(err),
        )
        return &AccountResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to create account",
            Error:   err.Error(),
        }
    }
    
    log.Info("Account created successfully",
        zap.Int("account_id", acc.ID),
        zap.Float64("balance", formatDecimal(acc.Balance)),
    )
    
    return &AccountResult{
        Success: true,
        Status:  http.StatusCreated,
        Message: "Account created successfully",
        Data:    acc,
    }
}

func (s *AccountService) GetAccount(ctx context.Context, id int) *AccountResult {
    log := middleware.GetLogger()
    
    log.Info("Retrieving account",
        zap.Int("account_id", id),
    )
    
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        log.Warn("Account not found",
            zap.Int("account_id", id),
            zap.Error(err),
        )
        return &AccountResult{
            Success: false,
            Status:  http.StatusNotFound,
            Message: "Account not found",
            Error:   err.Error(),
        }
    }
    
    log.Info("Account retrieved successfully",
        zap.Int("account_id", id),
        zap.Float64("balance", formatDecimal(account.Balance)),
    )
    
    return &AccountResult{
        Success: true,
        Status:  http.StatusOK,
        Message: "Account retrieved successfully",
        Data:    account,
    }
}

// isValidPrecision checks if a decimal has at most the specified number of decimal places
func isValidPrecision(d decimal.Decimal, maxDecimalPlaces int) bool {
    return d.Exponent() >= int32(-maxDecimalPlaces)
}

// formatDecimal formats a decimal to exactly 5 decimal places
func formatDecimal(d decimal.Decimal) float64 {
    return d.Round(5).InexactFloat64()
}
