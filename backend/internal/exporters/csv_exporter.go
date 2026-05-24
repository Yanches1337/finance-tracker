package exporters

import (
	"backend/internal/domain"
	"encoding/csv"
	"fmt"
	"os"
)

func ExportTransactionsToCSV(
	filePath string,
	transactions []domain.Transaction,
) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"ID",
		"Type",
		"Amount",
		"Category_id",
		"Date",
		"Description",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, t := range transactions {
		row := []string{
			fmt.Sprintf("%d", t.ID),
			string(t.Type),
			fmt.Sprintf("%.2f", t.Amount),
			fmt.Sprintf("%d", t.CategoryID),
			t.Date.Format("2006-01-02"),
			t.Description,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
