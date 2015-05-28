package main

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

const EnvStripeFrom = "STRIPE_FROM"
const EnvStripeTo = "STRIPE_TO"

type stripeAPI struct {
	from    *client.API
	to      *client.API
	pretend bool
}

func newStripeAPI(fromKey string, toKey string) *stripeAPI {
	s := stripeAPI{}
	if fromKey != "" {
		s.from = &client.API{}
		s.from.Init(fromKey, nil)
	}
	if toKey != "" {
		s.to = &client.API{}
		s.to.Init(toKey, nil)
	}
	return &s
}

func (s stripeAPI) SyncPlans() error {
	toPlans := map[string]*stripe.Plan{}
	missingPlans := map[string]*stripe.Plan{}

	params := &stripe.PlanListParams{}

	logger.Debug("Loading To plans...")
	iter := s.to.Plans.List(params)
	for iter.Next() {
		p := iter.Plan()
		toPlans[p.ID] = p
	}
	if err := iter.Err(); err != nil {
		return err
	}

	logger.Debugf("Loaded %d To plans. Loading From plans...", len(toPlans))
	iter = s.from.Plans.List(params)
	for iter.Next() {
		p := iter.Plan()
		_, has := toPlans[p.ID]
		if has {
			logger.Debugf("Matched plan: %s", p.ID)
			// TODO: Compare if matched plan is the same.
			continue
		}
		logger.Debugf("Missing plan: %s", p.ID)
		missingPlans[p.ID] = p
	}
	if err := iter.Err(); err != nil {
		return err
	}

	logger.Debugf("Found %d missing plans.", len(missingPlans))

	// TODO: Add missing plans.

	return nil
}
