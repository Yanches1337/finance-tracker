package domain

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	CategoryID  int64           `json:"category_id,omitempty"`
	Date        time.Time       `json:"date"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}

type CreateTransactionRequest struct {
	Type        TransactionType `json:"type" validate:"required"`
	Amount      float64         `json:"amount" validate:"required,gt=0"`
	CategoryID  int64           `json:"category_id"`
	Date        time.Time       `json:"date" validate:"required"`
	Description string          `json:"description"`
}

type UpdateTransactionRequest struct {
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	CategoryID  int64           `json:"category_id"`
	Date        time.Time       `json:"date"`
	Description string          `json:"description"`
}
