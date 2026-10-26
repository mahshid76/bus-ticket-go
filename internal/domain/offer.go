package domain

import (
	"context"
	"time"
)

// Query is one search: every trip from any stop of From to any stop of To
// whose departure falls on Date.
type Query struct {
	From string
	To   string
	Date time.Time
}

// Price is optional. Providers without a price source leave Offer.Price nil
// and send the visitor to their own site to see it.
type Price struct {
	Amount   float64
	Currency string
}

// Offer is one bookable trip from one provider.
type Offer struct {
	Provider  string
	Route     string
	Departure time.Time // in the origin stop's time zone
	Arrival   time.Time // in the destination stop's time zone
	FromStop  string
	ToStop    string
	Price     *Price
	BookURL   string
}

func (o Offer) Duration() time.Duration { return o.Arrival.Sub(o.Departure) }

// Provider is the extension point: one implementation per bus company or
// price source. Search must respect ctx cancellation.
type Provider interface {
	Name() string
	Search(ctx context.Context, q Query) ([]Offer, error)
}
