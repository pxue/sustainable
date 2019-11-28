package data

import (
	"upper.io/bond"
	"upper.io/db.v3"
)

type Product struct {
	ID        int64           `db:"id,omitempty"`
	Name      string          `db:"name"`
	Code      string          `db:"code"`
	Category  string          `db:"category"`
	Materials ProductMaterial `db:"materials"`
	Price     float64         `db:"price"`
}

type ProductMaterial map[string]float64

type ProductStore struct {
	bond.Store
}

func (*Product) CollectionName() string {
	return `products`
}

func (store ProductStore) FindByCode(code string) (*Product, error) {
	return store.FindOne(db.Cond{"code": code})
}

func (store ProductStore) FindOne(cond db.Cond) (*Product, error) {
	var product *Product
	if err := store.Find(cond).One(&product); err != nil {
		return nil, err
	}
	return product, nil
}
