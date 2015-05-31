package main

import (
	"fmt"
	"io"

	"github.com/stripe/stripe-go/client"
)

func NewSync() *syncThing {
	return &syncThing{
		target:  map[string]StripeThing{},
		missing: map[string]StripeThing{},
		changed: map[string]StripeThing{},
	}
}

type syncThing struct {
	target  map[string]StripeThing
	missing map[string]StripeThing
	changed map[string]StripeThing
}

func (s *syncThing) AddTarget(t StripeThing) {
	s.target[t.ID()] = t
}

func (s *syncThing) AddSource(t StripeThing) {
	id := t.ID()
	t2, has := s.target[id]
	if has && t2.Compare(t) {
		// Matched
		return
	}

	if has {
		// Changed
		s.changed[id] = t
		return
	}

	// Missing
	s.missing[id] = t
}

func (s *syncThing) SyncTarget(api *client.API) error {
	for _, thing := range s.missing {
		err := thing.New(api)
		if err != nil {
			return err
		}
	}

	for _, thing := range s.changed {
		err := thing.Update(api)
		if err != nil {
			return err
		}
	}
	return nil
}

// Diff writes a diff summary into a writer.
func (s syncThing) Diff(w io.Writer) (err error) {
	if len(s.target) > 0 {
		fmt.Fprintf(w, "%d loaded; %d missing, %d changed:\n", len(s.target), len(s.missing), len(s.changed))
	}
	for _, t := range s.missing {
		fmt.Fprintf(w, "+ %s\n", t)
	}
	for _, t := range s.changed {
		// TODO: Print diff
		// fmt.Fprintf(w, "~ %s\n  %s\n", t, s.target[id])
		fmt.Fprintf(w, "~ %s\n", t)
	}
	return
}
