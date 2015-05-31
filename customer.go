package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var errSyncCustomer = errors.New("cannot sync customers")

// Customer is a Thing wrapper for Stripe's Customer type.
type Customer struct {
	s *stripe.Customer
}

func (c Customer) ID() string {
	return c.s.ID
}

func (c Customer) New(api *client.API) error {
	return errSyncCustomer
}

func (c Customer) compareParams() *stripe.CustomerParams {
	defaultSource := stripe.SourceParams{}
	if c.s.DefaultSource != nil {
		// This might not be technically correct, as Token != ID necessarily?
		// But good enough for comparison purposes (probably shouldn't use for update params).
		defaultSource.Token = c.s.DefaultSource.ID
	}

	return &stripe.CustomerParams{
		Params: stripe.Params{
			Meta: c.s.Meta,
		},
		Desc:   c.s.Desc,
		Email:  c.s.Email,
		Source: &defaultSource,
	}
}

func (c Customer) Update(api *client.API) error {
	// TODO: Implement updates for what we can update (desc, meta?)
	return errSyncCustomer
}

func (c Customer) Compare(customer StripeThing) bool {
	return reflect.DeepEqual(c.compareParams(), customer.(*Customer).compareParams())
}

func (c Customer) String() string {
	// TODO: Something more useful here?
	return fmt.Sprintf("<Customer %s>", c.s.ID)
}
