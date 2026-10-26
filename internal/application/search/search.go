// Package search runs a query against every provider concurrently and
// merges the results.
package search

import (
	"context"
	"errors"
	"sync"
	"time"

	"bus-ticket-go/internal/domain"
)

type Options struct {
	Sort domain.SortKey
	Time domain.TimeGroup
}

type Result struct {
	Offers   []domain.Offer
	Failures map[string]error // providers that failed; the rest still answered
}

type UseCase struct {
	providers []domain.Provider
	timeout   time.Duration
}

func New(timeout time.Duration, providers ...domain.Provider) *UseCase {
	return &UseCase{providers: providers, timeout: timeout}
}

func (uc *UseCase) Execute(ctx context.Context, q domain.Query, opt Options) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	type answer struct {
		name   string
		offers []domain.Offer
		err    error
	}
	ch := make(chan answer, len(uc.providers))
	var wg sync.WaitGroup
	for _, p := range uc.providers {
		wg.Add(1)
		go func(p domain.Provider) {
			defer wg.Done()
			o, err := p.Search(ctx, q)
			ch <- answer{p.Name(), o, err}
		}(p)
	}
	wg.Wait()
	close(ch)

	res := Result{Failures: map[string]error{}}
	for a := range ch {
		if a.err != nil {
			res.Failures[a.name] = a.err
			continue
		}
		res.Offers = append(res.Offers, a.offers...)
	}
	if len(res.Failures) == len(uc.providers) && len(uc.providers) > 0 {
		var errs []error
		for _, e := range res.Failures {
			errs = append(errs, e)
		}
		return res, errors.Join(errs...)
	}

	res.Offers = domain.FilterByTimeGroup(res.Offers, opt.Time)
	sortKey := opt.Sort
	if sortKey == "" {
		sortKey = domain.SortEarliest
	}
	domain.SortOffers(res.Offers, sortKey)
	return res, nil
}
