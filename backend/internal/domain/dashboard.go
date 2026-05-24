package domain

type Dashboard struct {
	TotalIncome        float64         `json:"total_income"`
	TotalExpense       float64         `json:"total_expense"`
	Balance            float64         `json:"balance"`
	ExpensesByCategory []CategoryStats `json:"expenses_by_category"`
	IncomeByCategory   []CategoryStats `json:"income_by_category"`
}

type CategoryStats struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
}
