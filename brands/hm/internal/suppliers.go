package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type SupplierService service

type SupplierWrapper struct {
	ResponseStatusCode string             `json:"responseStatusCode"`
	Countries          []*SupplierCountry `json:"countries"`
}

type Supplier struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Factories []*SupplierFactory `json:"factories"`
}

type SupplierFactory struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	WorkersNumber string `json:"workersNumber"`
}

type SupplierCountry struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Suppliers []*Supplier `json:"suppliers"`
}

func (s *SupplierService) Get(ctx context.Context, productID string) (*SupplierWrapper, *http.Response, error) {
	path, _ := url.Parse(fmt.Sprintf("/%s/supplierDetails/articles/%s", s.client.opts.countryCode, productID))
	baseUrl, _ := url.Parse("https://www2.hm.com")
	u := baseUrl.ResolveReference(path)

	req, err := s.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var wrapper *SupplierWrapper
	resp, err := s.client.Do(ctx, req, &wrapper)
	if err != nil {
		return nil, resp, err
	}

	return wrapper, resp, nil
}
