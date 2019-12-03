package geocode

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

type Geocode struct {
	*maps.Client
}

var ErrNoResult = errors.New("no result")

func New(apiKey string) (*Geocode, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Geocode{client}, nil
}

func (g *Geocode) AddressToCoord(address string) ([]float64, error) {
	r := &maps.GeocodingRequest{
		Address: address,
	}
	resp, err := g.Geocode(context.Background(), r)
	if err != nil {
		return nil, err
	}

	for _, result := range resp {
		return []float64{
			result.Geometry.Location.Lat,
			result.Geometry.Location.Lng,
		}, nil
	}

	return nil, ErrNoResult
}
