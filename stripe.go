package main

import (
	"errors"

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
}

func newStripeAPI(sourceKey string, targetKey string) (*stripeAPI, error) {
	if sourceKey == "" || targetKey == "" {
		return nil, errMissingKey
	}
	s := stripeAPI{
		source: &client.API{},
		target: &client.API{},
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

	logger.Infof("Plans: %d loaded, %d missing, %d changed.", len(sync.target), len(sync.missing), len(sync.changed))

	if api.pretend {
		logger.Debug("Pretend mode: Stopping early.")
		return nil
	}

	return sync.SyncTarget(api.target)
}
