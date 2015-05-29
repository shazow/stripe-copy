package main

import (
	"reflect"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

// Plan is a Thing wrapper for Stripe's Plan type.
type Plan struct {
	s *stripe.Plan
}

func (p Plan) ID() string {
	return p.s.ID
}

func (p Plan) NewParams() *stripe.PlanParams {
	return &stripe.PlanParams{
		Params: stripe.Params{
			Meta: p.s.Meta,
		},
		ID:            p.s.ID,
		Name:          p.s.Name,
		Currency:      p.s.Currency,
		Amount:        p.s.Amount,
		Interval:      p.s.Interval,
		IntervalCount: p.s.IntervalCount,
		TrialPeriod:   p.s.TrialPeriod,
		Statement:     p.s.Statement,
	}
}

func (p Plan) New(api *client.API) error {
	_, err := api.Plans.New(p.NewParams())
	return err
}

func (p Plan) UpdateParams() *stripe.PlanParams {
	return &stripe.PlanParams{
		Params: stripe.Params{
			Meta: p.s.Meta,
		},
		Name:      p.s.Name,
		Statement: p.s.Statement,
	}
}

func (p Plan) Update(api *client.API) error {
	_, err := api.Plans.Update(p.ID(), p.UpdateParams())
	return err
}

func (p *Plan) Compare(plan StripeThing) bool {
	return reflect.DeepEqual(p.UpdateParams(), plan.(*Plan).UpdateParams())
}
