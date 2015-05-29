package main

import (
	"errors"
	"reflect"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

const envStripeSource = "STRIPE_SOURCE"
const envStripeTarget = "STRIPE_TARGET"

var errMissingKey = errors.New("Missing API key.")

type StripeThing interface {
	ID() string
	UpdateParams() interface{}
	NewParams() interface{}
	Compare(StripeThing) bool
}

type stripeAPI struct {
	source    *client.API
	target    *client.API
	pretend   bool
	stopAfter int
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

func (s stripeAPI) SyncPlans() error {
	targetPlans := map[string]*stripe.Plan{}
	missingPlans := map[string]*stripe.Plan{}
	changedPlans := map[string]*stripe.Plan{}

	params := &stripe.PlanListParams{}

	logger.Debug("Loading target plans...")
	iter := s.target.Plans.List(params)
	for iter.Next() {
		p := iter.Plan()
		targetPlans[p.ID] = p
	}
	if err := iter.Err(); err != nil {
		return err
	}

	logger.Debugf("Loaded %d To plans. Loading source plans...", len(targetPlans))
	iter = s.source.Plans.List(params)
	for iter.Next() {
		p := iter.Plan()
		p2, has := targetPlans[p.ID]
		if has && reflect.DeepEqual(p, p2) {
			logger.Debugf("Matched plan: %s", p.ID)
			continue
		}
		if has {
			logger.Debugf("Changed plan: %s\n%v\n%v", p.ID, p, p2)
			missingPlans[p.ID] = p
			continue
		}
		logger.Debugf("Missing plan: %s", p.ID)
		missingPlans[p.ID] = p
	}
	if err := iter.Err(); err != nil {
		return err
	}

	logger.Infof("Plans: %d loaded, %d missing, %d changed.", len(targetPlans), len(missingPlans), len(changedPlans))

	if s.pretend {
		logger.Debug("Pretend mode: Stopping early.")
		return nil
	}

	count := 0

	// Add missing plans
	for _, plan := range missingPlans {
		if s.stopAfter > 0 && count >= s.stopAfter {
			logger.Debugf("Stopping early: Completed %d stopAfter operations.", s.stopAfter)
			return nil
		}
		params := stripe.PlanParams{
			Params: stripe.Params{
				Meta: plan.Meta,
			},
			ID:            plan.ID,
			Name:          plan.Name,
			Currency:      plan.Currency,
			Amount:        plan.Amount,
			Interval:      plan.Interval,
			IntervalCount: plan.IntervalCount,
			TrialPeriod:   plan.TrialPeriod,
			Statement:     plan.Statement,
		}
		newPlan, err := s.target.Plans.New(&params)
		if err != nil {
			return err
		}
		logger.Debugf("Target: Plan created: %s", newPlan.ID)
		count++
	}

	for _, plan := range changedPlans {
		if s.stopAfter > 0 && count >= s.stopAfter {
			logger.Debugf("Stopping early: Completed %d stopAfter operations.", s.stopAfter)
			return nil
		}
		params := stripe.PlanParams{
			Params: stripe.Params{
				Meta: plan.Meta,
			},
			Name:      plan.Name,
			Statement: plan.Statement,
		}
		newPlan, err := s.target.Plans.Update(plan.ID, &params)
		if err != nil {
			return err
		}
		logger.Debugf("Target: Plan updated: %s", newPlan.ID)
		count++
	}

	logger.Infof("Created and updated %s plans.", count)

	return nil
}
