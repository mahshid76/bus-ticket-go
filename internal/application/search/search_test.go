package search

import (
	"context"
	"errors"
	"testing"
	"time"

	"bus-ticket-go/internal/domain"
)

type fake struct {
	name   string
	offers []domain.Offer
	err    error
}

func (f fake) Name() string { return f.name }
func (f fake) Search(context.Context, domain.Query) ([]domain.Offer, error) {
	return f.offers, f.err
}

func at(h int) time.Time { return time.Date(2026, 11, 14, h, 0, 0, 0, time.UTC) }

func TestMergesProvidersAndSorts(t *testing.T) {
	a := fake{name: "A", offers: []domain.Offer{{Provider: "A", Departure: at(10), Arrival: at(14)}}}
	b := fake{name: "B", offers: []domain.Offer{{Provider: "B", Departure: at(8), Arrival: at(13)}}}
	res, err := New(time.Second, a, b).Execute(context.Background(), domain.Query{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Offers) != 2 || res.Offers[0].Provider != "B" {
		t.Errorf("want B (08:00) first, got %+v", res.Offers)
	}
}

func TestOneProviderFailingKeepsTheRest(t *testing.T) {
	ok := fake{name: "ok", offers: []domain.Offer{{Provider: "ok", Departure: at(9), Arrival: at(12)}}}
	bad := fake{name: "bad", err: errors.New("boom")}
	res, err := New(time.Second, ok, bad).Execute(context.Background(), domain.Query{}, Options{})
	if err != nil || len(res.Offers) != 1 || res.Failures["bad"] == nil {
		t.Errorf("err=%v offers=%d failures=%v", err, len(res.Offers), res.Failures)
	}
}

func TestAllProvidersFailingIsAnError(t *testing.T) {
	_, err := New(time.Second, fake{name: "x", err: errors.New("boom")}).Execute(context.Background(), domain.Query{}, Options{})
	if err == nil {
		t.Error("want error when every provider fails")
	}
}
