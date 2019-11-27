package data

import (
	"upper.io/bond"
	"upper.io/db.v3"
)

type Factory struct {
	ID         int64  `db:"id,omitempty"`
	SupplierID int64  `db:"supplier_id"`
	HmID       string `db:"hm_id"`
	Name       string `db:"name"`
	Address    string `db:"address"`
	Country    string `db:"country"`
}

type FactoryStore struct {
	bond.Store
}

func (*Factory) CollectionName() string {
	return `factories`
}

func (store FactoryStore) FindByHMID(hmID string) (*Factory, error) {
	return store.FindOne(db.Cond{"hm_id": hmID})
}

func (store FactoryStore) FindOne(cond db.Cond) (*Factory, error) {
	var factory *Factory
	if err := store.Find(cond).One(&factory); err != nil {
		return nil, err
	}
	return factory, nil
}
