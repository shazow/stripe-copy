package main

import (
	"bytes"
	"errors"
	"io"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

const envStripeSource = "STRIPE_SOURCE"
const envStripeTarget = "STRIPE_TARGET"

var errMissingKey = errors.New("Missing API key.")

type StripeThing interface {
	ID() string
	New(*client.API) error
	Update(*client.API) error
	Compare(StripeThing) bool
}

type stripeAPI struct {
	source  *client.API
	target  *client.API
	pretend bool
	out     io.Writer
}

func newStripeAPI(sourceKey string, targetKey string) (*stripeAPI, error) {
	if sourceKey == "" || targetKey == "" {
		return nil, errMissingKey
	}
	s := stripeAPI{
		source: &client.API{},
		target: &client.API{},
		out:    &bytes.Buffer{},
	}
	s.source.Init(sourceKey, nil)
	s.target.Init(targetKey, nil)
	return &s, nil
}

func (api stripeAPI) SyncPlans() error {
	sync := NewSync()

	logger.Debug("Loading target plans...")
	params := &stripe.PlanListParams{}
	iter := api.target.Plans.List(params)
	for iter.Next() {
		p := iter.Plan()
		sync.AddTarget(&Plan{p})
	}
	if err := iter.Err(); err != nil {
		return err
	}

	logger.Debugf("Loaded %d target plans. Loading source plans...", len(sync.target))
	iter = api.source.Plans.List(params)
	for iter.Next() {
		p := iter.Plan()
		sync.AddSource(&Plan{p})
	}
	if err := iter.Err(); err != nil {
		return err
	}

	sync.Diff(api.out)
	if api.pretend {
		logger.Debug("SyncPlans: Pretend mode, stopping early.")
		return nil
	}

	return sync.SyncTarget(api.target)
}

func (api stripeAPI) CheckCustomers() ([]Customer, error) {
	customers := []Customer{}
	sync := NewSync()

	logger.Debug("Loading target customers...")
	params := &stripe.CustomerListParams{}
	iter := api.target.Customers.List(params)
	for iter.Next() {
		p := iter.Customer()
		sync.AddTarget(&Customer{p})
	}
	if err := iter.Err(); err != nil {
		return customers, err
	}

	logger.Debugf("Loaded %d target customers. Loading source customers...", len(sync.target))
	iter = api.source.Customers.List(params)
	for iter.Next() {
		p := iter.Customer()
		customer := Customer{p}
		sync.AddSource(&customer)
		customers = append(customers, customer)
	}
	if err := iter.Err(); err != nil {
		return customers, err
	}

	sync.Diff(api.out)
	if api.pretend {
		logger.Debug("CheckCustomers: Pretend mode, stopping early.")
		return customers, nil
	}

	return customers, sync.SyncTarget(api.target)
}

func (api stripeAPI) SyncSubs(customers []Customer) error {
	sync := NewSync()

	logger.Debug("Loading target subscriptions...")
	for _, customer := range customers {
		params := &stripe.SubListParams{
			Customer: customer.ID(),
		}
		iter := api.target.Subs.List(params)
		for iter.Next() {
			s := iter.Sub()
			sync.AddTarget(&Sub{s})
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}
	logger.Debugf("Loaded %d target subscriptions. Loading source subscriptions...", len(sync.target))

	for _, customer := range customers {
		params := &stripe.SubListParams{
			Customer: customer.ID(),
		}
		iter := api.source.Subs.List(params)
		for iter.Next() {
			s := iter.Sub()
			sync.AddSource(&Sub{s})
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}

	sync.Diff(api.out)
	if api.pretend {
		logger.Debug("SyncSubs: Pretend mode, stopping early.")
		return nil
	}

	return sync.SyncTarget(api.target)
}
