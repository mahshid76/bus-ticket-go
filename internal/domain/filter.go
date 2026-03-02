package domain

import "time"

func FilterByTimeGroup(trips []Trip, group TimeGroup) []Trip {
	out := []Trip{}

	for _, t := range trips {
		dep, err := time.Parse("15:04", t.DepartureTime)
		if err != nil {
			continue
		}
		if GroupByTime(dep) == group {
			out = append(out, t)
		}
	}
	return out
}
