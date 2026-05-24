package domain

type Category struct {
	ID     int64           `json:"id"`
	UserID int64           `json:"user_id"`
	Name   string          `json:"name"`
	Type   TransactionType `json:"type"`
	Color  string          `json:"color,omitempty"`
	Icon   string          `json:"icon,omitempty"`
}

type CreateCategoryRequest struct {
	Name  string          `json:"name" validate:"required,min=1,max=100"`
	Type  TransactionType `json:"type" validate:"required"`
	Color string          `json:"color,omitempty"`
	Icon  string          `json:"icon,omitempty"`
}

type UpdateCategoryRequest struct {
	ID    int64           `json:"id" validate:"required"`
	Name  string          `json:"name" validate:"required,min=1,max=100"`
	Type  TransactionType `json:"type" validate:"required"`
	Color string          `json:"color,omitempty"`
	Icon  string          `json:"icon,omitempty"`
}
