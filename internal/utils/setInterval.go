package utils

import (
	"time"
)

func AddInterval(start time.Time, unit int, length int) time.Time {
	// based on intervals table in database
	// 1 = days
	// 2 = week
	// 3 = month
	// 4 = year
	switch unit {
	case 1:
		// days
		return start.AddDate(0, 0, length)
	case 2:
		// week
		return start.AddDate(0, 0, length*7)
	case 3:
		// month
		return start.AddDate(0, length, 0)
	case 4:
		// year
		return start.AddDate(length, 0, 0)
	default:
		return start
	}
}

func CalculateReminderDates(
	startDate time.Time,
	unit int,
	length int,
) (next, warning, remove time.Time) {

	dueDate := AddInterval(startDate, unit, length)
	remove = dueDate

	// CASE 1: DAY-BASED ≤ 7
	if unit == 1 && length <= 7 {
		switch length {
		case 1:
			next = startDate.AddDate(0, 0, 1)
			warning = next
			remove = next

		case 2:
			next = startDate.AddDate(0, 0, 1)
			warning = next

		case 3:
			next = startDate.AddDate(0, 0, 2)
			warning = next

		case 4:
			next = startDate.AddDate(0, 0, 2)
			warning = next

		case 5:
			next = startDate.AddDate(0, 0, 3)
			warning = next

		case 6:
			next = startDate.AddDate(0, 0, 3)
			warning = next

		case 7:
			next = startDate.AddDate(0, 0, 1)
			warning = startDate.AddDate(0, 0, 4)
		}

		return
	}

	// CASE 2: ≥ 1 WEEK / MONTH / YEAR
	warning = remove.AddDate(0, 0, -3)
	next = warning.AddDate(0, 0, -3)

	return
}
