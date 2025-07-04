package repository

import (
    "context"
    "database/sql"
    "transfer-service/model"
    "github.com/shopspring/decimal"
)

type AccountRepository interface {
    Create(ctx context.Context, account model.Account) error
    GetByID(ctx context.Context, id int) (*model.Account, error)
    GetByIDWithLock(ctx context.Context, tx *sql.Tx, id int) (*model.Account, error)
    UpdateBalance(ctx context.Context, id int, newBalance decimal.Decimal) error
    UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int, newBalance decimal.Decimal) error
    DeleteByID(ctx context.Context, id int) error
    GetDB() *sql.DB
}

type accountRepo struct {
    db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
    return &accountRepo{db: db}
}

func (r *accountRepo) Create(ctx context.Context, a model.Account) error {
    _, err := r.db.ExecContext(ctx, "INSERT INTO accounts (id, balance) VALUES ($1, $2)", a.ID, a.Balance)
    return err
}

func (r *accountRepo) GetByID(ctx context.Context, id int) (*model.Account, error) {
    var a model.Account
    err := r.db.QueryRowContext(ctx, "SELECT id, balance FROM accounts WHERE id = $1", id).Scan(&a.ID, &a.Balance)
    if err != nil {
        return nil, err
    }
    return &a, nil
}

// GetByIDWithLock uses SELECT FOR UPDATE to lock the row for update
func (r *accountRepo) GetByIDWithLock(ctx context.Context, tx *sql.Tx, id int) (*model.Account, error) {
    var a model.Account
    err := tx.QueryRowContext(ctx, "SELECT id, balance FROM accounts WHERE id = $1 FOR UPDATE", id).Scan(&a.ID, &a.Balance)
    if err != nil {
        return nil, err
    }
    return &a, nil
}

func (r *accountRepo) UpdateBalance(ctx context.Context, id int, newBalance decimal.Decimal) error {
    _, err := r.db.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", newBalance, id)
    return err
}

// UpdateBalanceWithTx updates balance within a transaction
func (r *accountRepo) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int, newBalance decimal.Decimal) error {
    _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = $1 WHERE id = $2", newBalance, id)
    return err
}

func (r *accountRepo) DeleteByID(ctx context.Context, id int) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM accounts WHERE id = $1", id)
    return err
}

func (r *accountRepo) GetDB() *sql.DB {
    return r.db
}
