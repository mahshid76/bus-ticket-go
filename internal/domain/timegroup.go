package domain

type TimeGroup string

const (
	Morning   TimeGroup = "morning"   // before 12:00
	Afternoon TimeGroup = "afternoon" // 12:00 to 17:59
	Evening   TimeGroup = "evening"   // 18:00 to 23:59, including night departures
)

func GroupOf(o Offer) TimeGroup {
	switch h := o.Departure.Hour(); {
	case h < 12:
		return Morning
	case h < 18:
		return Afternoon
	default:
		return Evening
	}
}

// FilterByTimeGroup keeps offers departing in g. An empty g keeps all.
func FilterByTimeGroup(offers []Offer, g TimeGroup) []Offer {
	if g == "" {
		return offers
	}
	out := make([]Offer, 0, len(offers))
	for _, o := range offers {
		if GroupOf(o) == g {
			out = append(out, o)
		}
	}
	return out
}
