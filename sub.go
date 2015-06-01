package main

import (
	"fmt"
	"reflect"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

// Sub is a Thing wrapper for Stripe's Sub type.
type Sub struct {
	s *stripe.Sub
}

func (s Sub) ID() string {
	return s.s.ID
}

func (s Sub) newParams() *stripe.SubParams {
	return &stripe.SubParams{
		// TODO: Add support for discounts/coupons
		Params: stripe.Params{
			Meta: s.s.Meta,
		},
		Customer:   s.s.Customer.ID,
		Plan:       s.s.Plan.ID,
		TrialEnd:   s.s.TrialEnd, // TODO: Investigate if we need to do TrialEndNow when null?
		Quantity:   s.s.Quantity,
		FeePercent: s.s.FeePercent,
		TaxPercent: s.s.TaxPercent,
		EndCancel:  s.s.EndCancel,
		NoProrate:  true,
		// TODO: Add billing_cycle_anchor once https://github.com/stripe/stripe-go/issues/105 is resolved.
	}
}

func (s Sub) New(api *client.API) error {
	_, err := api.Subs.New(s.newParams())
	return err
}

func (s Sub) updateParams() *stripe.SubParams {
	// Seems SubParams for new and update are the same.
	return s.newParams()
}

func (s Sub) Update(api *client.API) error {
	_, err := api.Subs.Update(s.ID(), s.updateParams())
	return err
}

func (s Sub) Compare(sub StripeThing) bool {
	// TODO: Be smarter about comparing attributes which we can't update?
	return reflect.DeepEqual(s.updateParams(), sub.(*Sub).updateParams())
}

func (s Sub) String() string {
	// TODO: Something prettier here?
	return fmt.Sprintf("<Sub %s>", s.s.ID)
}
