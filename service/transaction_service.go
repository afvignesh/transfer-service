package service

import (
    "context"
    "errors"
    "transfer-service/model"
    "transfer-service/repository"
    "database/sql"
    "net/http"
    "transfer-service/middleware"
    "go.uber.org/zap"
)

type TransactionService struct {
    accountRepo     repository.AccountRepository
    transactionRepo repository.TransactionRepository
}

var ErrInsufficientBalance = errors.New("insufficient balance")
var ErrSourceAccountNotFound = errors.New("source account not found")
var ErrDestinationAccountNotFound = errors.New("destination account not found")

func NewTransactionService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) *TransactionService {
    return &TransactionService{
        accountRepo:     accountRepo,
        transactionRepo: transactionRepo,
    }
}

// TransferResult represents the result of a transfer operation
type TransferResult struct {
    Success bool
    Status  int
    Message string
    Error   string
    Data    interface{}
}

func (s *TransactionService) Transfer(ctx context.Context, t model.Transaction) *TransferResult {
    log := middleware.GetLogger()
    
    // Validate amount precision (5 decimal places)
    if !isValidPrecision(t.Amount, 5) {
        log.Warn("Transfer failed - invalid amount precision",
            zap.String("amount", t.Amount.String()),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusBadRequest,
            Message: "Amount must have at most 5 decimal places",
            Error:   "invalid precision",
        }
    }
    
    log.Info("Starting transfer",
        zap.Int("source_account_id", t.SourceAccountID),
        zap.Int("destination_account_id", t.DestinationAccountID),
        zap.Float64("amount", t.Amount.Round(5).InexactFloat64()),
    )
    
    // Check if source and destination accounts are the same
    if t.SourceAccountID == t.DestinationAccountID {
        log.Warn("Transfer rejected - same source and destination accounts",
            zap.Int("account_id", t.SourceAccountID),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusBadRequest,
            Message: "Source and destination accounts are the same",
            Error:   "same accounts",
        }
    }

    // Start a database transaction
    tx, err := s.accountRepo.GetDB().BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable, // Highest isolation level for financial transactions
    })
    if err != nil {
        log.Error("Failed to start database transaction",
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to start transaction",
            Error:   err.Error(),
        }
    }
    
    // Ensure transaction is rolled back on error
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // Lock and get source account with FOR UPDATE
    from, err := s.accountRepo.GetByIDWithLock(ctx, tx, t.SourceAccountID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Warn("Transfer failed - source account not found",
                zap.Int("source_account_id", t.SourceAccountID),
            )
            return &TransferResult{
                Success: false,
                Status:  http.StatusNotFound,
                Message: "Source account not found",
                Error:   err.Error(),
            }
        }
        log.Error("Failed to get source account",
            zap.Int("source_account_id", t.SourceAccountID),
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to get source account",
            Error:   err.Error(),
        }
    }

    // Lock and get destination account with FOR UPDATE
    to, err := s.accountRepo.GetByIDWithLock(ctx, tx, t.DestinationAccountID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Warn("Transfer failed - destination account not found",
                zap.Int("destination_account_id", t.DestinationAccountID),
            )
            return &TransferResult{
                Success: false,
                Status:  http.StatusNotFound,
                Message: "Destination account not found",
                Error:   err.Error(),
            }
        }
        log.Error("Failed to get destination account",
            zap.Int("destination_account_id", t.DestinationAccountID),
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to get destination account",
            Error:   err.Error(),
        }
    }

    // Check balance with locked data
    if from.Balance.LessThan(t.Amount) {
        log.Warn("Transfer failed - insufficient balance",
            zap.Int("source_account_id", t.SourceAccountID),
            zap.Float64("current_balance", from.Balance.Round(5).InexactFloat64()),
            zap.Float64("requested_amount", t.Amount.Round(5).InexactFloat64()),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusBadRequest,
            Message: "Insufficient balance",
            Error:   "insufficient balance",
        }
    }

    // Calculate new balances
    from.Balance = from.Balance.Sub(t.Amount)
    to.Balance = to.Balance.Add(t.Amount)

    log.Info("Updating account balances",
        zap.Int("source_account_id", from.ID),
        zap.Float64("source_new_balance", from.Balance.Round(5).InexactFloat64()),
        zap.Int("destination_account_id", to.ID),
        zap.Float64("destination_new_balance", to.Balance.Round(5).InexactFloat64()),
    )

    // Update both accounts within the transaction
    if err := s.accountRepo.UpdateBalanceWithTx(ctx, tx, from.ID, from.Balance); err != nil {
        log.Error("Failed to update source account balance",
            zap.Int("source_account_id", from.ID),
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to update source account",
            Error:   err.Error(),
        }
    }

    if err := s.accountRepo.UpdateBalanceWithTx(ctx, tx, to.ID, to.Balance); err != nil {
        log.Error("Failed to update destination account balance",
            zap.Int("destination_account_id", to.ID),
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to update destination account",
            Error:   err.Error(),
        }
    }

    // Log the transaction (within the same database transaction)
    loggedTx, err := s.transactionRepo.Create(ctx, t)
    if err != nil {
        log.Error("Failed to log transaction",
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to log transaction",
            Error:   err.Error(),
        }
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        log.Error("Failed to commit database transaction",
            zap.Error(err),
        )
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to commit transaction",
            Error:   err.Error(),
        }
    }

    log.Info("Transfer completed successfully",
        zap.Int("transaction_id", loggedTx.ID),
        zap.Int("source_account_id", t.SourceAccountID),
        zap.Int("destination_account_id", t.DestinationAccountID),
        zap.Float64("amount", t.Amount.Round(5).InexactFloat64()),
    )

    return &TransferResult{
        Success: true,
        Status:  http.StatusOK,
        Message: "Transfer completed successfully",
        Data: map[string]interface{}{
            "message": "Transfer completed successfully",
            "transaction": loggedTx,
        },
    }
}

func (s *TransactionService) GetTransactionHistory(ctx context.Context) *TransferResult {
    transactions, err := s.transactionRepo.GetAll(ctx)
    if err != nil {
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to retrieve transaction history",
            Error:   err.Error(),
        }
    }
    
    return &TransferResult{
        Success: true,
        Status:  http.StatusOK,
        Message: "Transaction history retrieved successfully",
        Data:    transactions,
    }
}

func (s *TransactionService) GetAccountTransactionHistory(ctx context.Context, accountID int) *TransferResult {
    transactions, err := s.transactionRepo.GetByAccountID(ctx, accountID)
    if err != nil {
        return &TransferResult{
            Success: false,
            Status:  http.StatusInternalServerError,
            Message: "Failed to retrieve account transaction history",
            Error:   err.Error(),
        }
    }
    
    return &TransferResult{
        Success: true,
        Status:  http.StatusOK,
        Message: "Account transaction history retrieved successfully",
        Data:    transactions,
    }
}
