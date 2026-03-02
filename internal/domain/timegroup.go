package domain

import "time"

type TimeGroup string

const (
	Morning   TimeGroup = "Morning"
	afternoon TimeGroup = "Afternoon"
	Evening   TimeGroup = "Evening"
)

func GroupByTime(t time.Time) TimeGroup {
	h := t.Hour()

	switch {
	case h < 12:
		return Morning
	case h < 18:
		return afternoon
	default:
		return Evening
	}
}
