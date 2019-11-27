package data

import "upper.io/bond"

type Product struct {
	ID        int64           `db:"id,omitempty"`
	Name      string          `db:"name"`
	Code      string          `db:"code"`
	Category  string          `db:"category"`
	Materials ProductMaterial `db:"materials"`
}

type ProductMaterial map[string]float64

type ProductStore struct {
	bond.Store
}

func (*Product) CollectionName() string {
	return `products`
}
