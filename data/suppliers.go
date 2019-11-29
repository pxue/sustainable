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

func (store SupplierStore) FindOne(cond ...interface{}) (*Supplier, error) {
	var supplier *Supplier
	if err := store.Find(cond...).One(&supplier); err != nil {
		return nil, err
	}
	return supplier, nil
}

type SupplierWithRank struct {
	*Supplier
	Rank float64 `db:"rank"`
}

func (store SupplierStore) SearchByName(name string) ([]*SupplierWithRank, error) {
	withQuery := db.Raw("*, to_tsvector(name) txtsearch, plainto_tsquery(?) query", name)
	innerSelect := DB.Select(withQuery).From("suppliers")

	withRank := db.Raw("t.*, ts_rank(txtsearch, query) as rank")
	query := DB.Select(withRank).From(innerSelect).As("t").OrderBy("rank desc").Limit(10)

	var rankings []*SupplierWithRank
	if err := query.All(&rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}
