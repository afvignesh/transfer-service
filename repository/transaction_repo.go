package repository

import (
    "context"
    "database/sql"
    "transfer-service/model"
)

type TransactionRepository interface {
    Create(ctx context.Context, tx model.Transaction) (*model.Transaction, error)
    GetByID(ctx context.Context, id int) (*model.Transaction, error)
    GetByAccountID(ctx context.Context, accountID int) ([]*model.Transaction, error)
    GetAll(ctx context.Context) ([]*model.Transaction, error)
}

type transactionRepo struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
    return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(ctx context.Context, t model.Transaction) (*model.Transaction, error) {
    var id int
    err := r.db.QueryRowContext(ctx, 
        "INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3) RETURNING id",
        t.SourceAccountID, t.DestinationAccountID, t.Amount,
    ).Scan(&id)
    
    if err != nil {
        return nil, err
    }
    
    // Get the created transaction with timestamp
    return r.GetByID(ctx, id)
}

func (r *transactionRepo) GetByID(ctx context.Context, id int) (*model.Transaction, error) {
    var t model.Transaction
    err := r.db.QueryRowContext(ctx, 
        "SELECT id, source_account_id, destination_account_id, amount, created_at FROM transactions WHERE id = $1",
        id,
    ).Scan(&t.ID, &t.SourceAccountID, &t.DestinationAccountID, &t.Amount, &t.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &t, nil
}

func (r *transactionRepo) GetByAccountID(ctx context.Context, accountID int) ([]*model.Transaction, error) {
    rows, err := r.db.QueryContext(ctx, 
        "SELECT id, source_account_id, destination_account_id, amount, created_at FROM transactions WHERE source_account_id = $1 OR destination_account_id = $1 ORDER BY created_at DESC",
        accountID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var transactions []*model.Transaction
    for rows.Next() {
        var t model.Transaction
        err := rows.Scan(&t.ID, &t.SourceAccountID, &t.DestinationAccountID, &t.Amount, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, &t)
    }
    
    return transactions, nil
}

func (r *transactionRepo) GetAll(ctx context.Context) ([]*model.Transaction, error) {
    rows, err := r.db.QueryContext(ctx, 
        "SELECT id, source_account_id, destination_account_id, amount, created_at FROM transactions ORDER BY created_at DESC",
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var transactions []*model.Transaction
    for rows.Next() {
        var t model.Transaction
        err := rows.Scan(&t.ID, &t.SourceAccountID, &t.DestinationAccountID, &t.Amount, &t.CreatedAt)
        if err != nil {
            return nil, err
        }
        transactions = append(transactions, &t)
    }
    
    return transactions, nil
} 