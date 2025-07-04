package model

import (
    "encoding/json"
    "github.com/shopspring/decimal"
)

type Account struct {
    ID      int             `json:"account_id"`
    Balance decimal.Decimal `json:"balance"`
}

// MarshalJSON customizes JSON marshaling to format balance with 5 decimal places
func (a Account) MarshalJSON() ([]byte, error) {
    type Alias Account
    return json.Marshal(&struct {
        *Alias
        Balance float64 `json:"balance"`
    }{
        Alias:   (*Alias)(&a),
        Balance: a.Balance.Round(5).InexactFloat64(),
    })
}
