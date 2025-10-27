package date_utils

import (
	"fmt"
	"time"
)

// CalculateDaysRemaining menghitung hari tersisa dari sekarang ke target date
func CalculateDaysRemaining(endDateStr string) int {
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return 0
	}

	now := time.Now()
	// Set ke midnight untuk perhitungan yang akurat
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDateMidnight := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	diff := endDateMidnight.Sub(today)
	days := int(diff.Hours() / 24)

	return days
}

// GetDaysRemainingStatus mengembalikan status text berdasarkan hari tersisa
func GetDaysRemainingStatus(daysRemaining int) string {
	switch {
	case daysRemaining < 0:
		if daysRemaining == -1 {
			return "Overdue by 1 day"
		}
		return fmt.Sprintf("Overdue by %d days", -daysRemaining)
	case daysRemaining == 0:
		return "Due today"
	case daysRemaining == 1:
		return "1 day left"
	case daysRemaining <= 7:
		return fmt.Sprintf("%d days left (due this week)", daysRemaining)
	case daysRemaining <= 30:
		return fmt.Sprintf("%d days left (due this month)", daysRemaining)
	default:
		return fmt.Sprintf("%d days left", daysRemaining)
	}
}
