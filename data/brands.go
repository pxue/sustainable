package data

import (
	"upper.io/bond"
	"upper.io/db.v3"
)

type Brand struct {
	ID          int64  `db:"id,omitempty" json:"id"`
	Name        string `db:"name" json:"name"`
	FactoryLink string `db:"factory_link" json:"factoryLink"`
}

type BrandStore struct {
	bond.Store
}

func (*Brand) CollectionName() string {
	return `brands`
}

func (store BrandStore) FindByName(name string) (*Brand, error) {
	return store.FindOne(db.Cond{"name": name})
}

func (store BrandStore) FindOne(cond db.Cond) (*Brand, error) {
	var b *Brand
	if err := store.Find(cond).One(&b); err != nil {
		return nil, err
	}
	return b, nil
}
