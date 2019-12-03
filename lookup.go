package lookup

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

type Lookup interface {
	Open(context.Context, string) error	
	Find(string) (interface{}, bool)
}

type LookupWriter interface {
	Add(interface{}) error
}

var lookup_providers roster.Roster

func ensureRoster() error {

	if lookup_providers == nil {
		
		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		lookup_providers = r
	}

	return nil
}

func RegisterLookup(ctx context.Context, name string, lu Lookup) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return lookup_providers.Register(ctx, name, lu)
}

func NewLookup(ctx context.Context, uri string) (Lookup, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := lookup_providers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	lu := i.(Lookup)
	
	err = lu.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return lu, nil
}
