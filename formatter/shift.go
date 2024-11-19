package formatter

import "api-kasirapp/models"

type ShiftFormatter struct {
	ID           int     `json:"id"`
	UserID       int     `json:"name"`
	StartBalance float64 `json:"start_balance"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Status       string  `json:"status"`
	TotalSales   float64 `json:"total_sales"`
	Expenses     float64 `json:"expenses"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func FormatShift(shift models.Shift) ShiftFormatter {
	formatter := ShiftFormatter{
		ID:           shift.ID,
		UserID:       shift.UserID.ID,
		StartBalance: shift.StartBalance,
		StartTime:    shift.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:      shift.EndTime.Format("2006-01-02 15:04:05"),
		Status:       shift.Status,
		TotalSales:   shift.TotalSales,
		Expenses:     shift.Expenses,
		CreatedAt:    shift.CreatedAt,
		UpdatedAt:    shift.UpdatedAt,
	}
	return formatter
}

func FormatShifts(shifts []models.Shift) []ShiftFormatter {
	var formatter []ShiftFormatter
	for _, shift := range shifts {
		formatter = append(formatter, FormatShift(shift))
	}
	return formatter
}
