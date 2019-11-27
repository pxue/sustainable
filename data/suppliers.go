package data

import (
	"upper.io/bond"
	"upper.io/db.v3"
)

type Supplier struct {
	ID      int64  `db:"id,omitempty"`
	Name    string `db:"name"`
	Country string `db:"country"`
	HmID    string `db:"hm_id"`
}

type SupplierStore struct {
	bond.Store
}

func (*Supplier) CollectionName() string {
	return `suppliers`
}

func (store SupplierStore) FindByHMID(hmID string) (*Supplier, error) {
	return store.FindOne(db.Cond{"hm_id": hmID})
}

func (store SupplierStore) FindOne(cond db.Cond) (*Supplier, error) {
	var supplier *Supplier
	if err := store.Find(cond).One(&supplier); err != nil {
		return nil, err
	}
	return supplier, nil
}
