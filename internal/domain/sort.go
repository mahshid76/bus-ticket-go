package domain

import "sort"

type SortKey string

const (
	SortEarliest SortKey = "dep"
	SortFastest  SortKey = "dur"
	SortCheapest SortKey = "price"
)

// SortOffers orders offers in place. Offers without a price sort after
// priced ones when sorting by price.
func SortOffers(offers []Offer, by SortKey) {
	sort.SliceStable(offers, func(i, j int) bool {
		a, b := offers[i], offers[j]
		switch by {
		case SortCheapest:
			if (a.Price == nil) != (b.Price == nil) {
				return a.Price != nil
			}
			if a.Price != nil && a.Price.Amount != b.Price.Amount {
				return a.Price.Amount < b.Price.Amount
			}
		case SortFastest:
			if a.Duration() != b.Duration() {
				return a.Duration() < b.Duration()
			}
		}
		return a.Departure.Before(b.Departure)
	})
}

// HasPrices reports whether any offer carries a price.
func HasPrices(offers []Offer) bool {
	for _, o := range offers {
		if o.Price != nil {
			return true
		}
	}
	return false
}
