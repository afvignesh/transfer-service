package model

import (
    "encoding/json"
    "time"
    "github.com/shopspring/decimal"
)

type Transaction struct {
    ID                   int             `json:"id,omitempty"`
    SourceAccountID      int             `json:"source_account_id"`
    DestinationAccountID int             `json:"destination_account_id"`
    Amount               decimal.Decimal `json:"amount"`
    CreatedAt            time.Time       `json:"created_at,omitempty"`
}

// MarshalJSON customizes JSON marshaling to format amount with 5 decimal places
func (t Transaction) MarshalJSON() ([]byte, error) {
    type Alias Transaction
    return json.Marshal(&struct {
        *Alias
        Amount float64 `json:"amount"`
    }{
        Alias:  (*Alias)(&t),
        Amount: t.Amount.Round(5).InexactFloat64(),
    })
}