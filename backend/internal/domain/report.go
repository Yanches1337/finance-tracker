package domain

import "time"

type ReportFormat string

const (
	ReportFormatCSV  ReportFormat = "csv"
	ReportFormatJSON ReportFormat = "json"
	ReportFormatXLSX ReportFormat = "xlsx"
	ReportFormatPDF  ReportFormat = "pdf"
)

type Report struct {
	ID        int64        `json:"id"`
	UserID    int64        `json:"user_id"`
	Format    ReportFormat `json:"format"`
	FilePath  string       `json:"file_path"`
	FileName  string       `json:"file_name"`
	FromDate  time.Time    `json:"from_date"`
	ToDate    time.Time    `json:"to_date"`
	CreatedAt time.Time    `json:"created_at"`
}

type GenerateReportRequest struct {
	Format   ReportFormat `json:"format"`
	FromDate time.Time    `json:"from_date"`
	ToDate   time.Time    `json:"to_date"`
}
