package domain

import "sort"

func SortByPrice(trips []Trip) {
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].PriceEUR < trips[j].PriceEUR
	})
}
