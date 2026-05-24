package domain

import "time"

type Goal struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Name          string    `json:"name"`
	TargetAmount  int       `json:"target_amount"`
	CurrentAmount int       `json:"current_amount"`
	TargetDate    time.Time `json:"target_date,omitempty"`
	Description   string    `json:"description,omitempty"`
	IsCompleted   bool      `json:"is_completed"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateGoalRequest struct {
	Name         string    `json:"name" validate:"required"`
	TargetAmount int       `json:"target_amount" validate:"required,gt=0"`
	TargetDate   time.Time `json:"target_date"`
	Description  string    `json:"description"`
}

type UpdateGoalRequest struct {
	Name          string    `json:"name"`
	TargetAmount  int       `json:"target_amount"`
	CurrentAmount int       `json:"current_amount"`
	TargetDate    time.Time `json:"target_date"`
	Description   string    `json:"description"`
	IsCompleted   bool      `json:"is_completed"`
}
