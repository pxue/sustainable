package hm

import (
	"context"
	"net/http"

	internal "github.com/pxue/sustainable/brands/hm/internal"
)

type Provider struct {
	*internal.Client
}

func (p *Provider) ShopByProduct(ctx context.Context, pageSize, offset int) ([]*internal.ProductListItem, *http.Response, error) {
	params := &internal.ProductListParams{
		Sort:      "stock",
		ImageSize: "small",
		Image:     "model",
		Offset:    offset,
		PageSize:  pageSize,
	}
	return p.ProductList.ShopByProduct(ctx, params)
}

func (p *Provider) GetProduct(ctx context.Context, code string) (*internal.Product, *http.Response, error) {
	return p.Product.Get(ctx, code)
}

func (p *Provider) GetSupplier(ctx context.Context, code string) (*internal.SupplierWrapper, *http.Response, error) {
	return p.Supplier.Get(ctx, code)
}

func New(debug bool) (*Provider, error) {
	client, err := internal.NewClient(
		nil,
		internal.Country("hm-canada", "en_ca"),
		internal.Debug(debug),
	)
	if err != nil {
		return nil, err
	}
	return &Provider{client}, nil
}
