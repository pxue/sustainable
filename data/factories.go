package data

import (
	"upper.io/bond"
	"upper.io/db.v3"
)

type Factory struct {
	ID         int64   `db:"id,omitempty"`
	SupplierID *int64  `db:"supplier_id,omitempty"`
	Code       string  `db:"code"`
	Name       string  `db:"name"`
	Address    string  `db:"address"`
	Country    string  `db:"country"`
	Lat        float64 `db:"lat"`
	Lon        float64 `db:"lon"`
}

type FactoryStore struct {
	bond.Store
}

func (*Factory) CollectionName() string {
	return `factories`
}

func (store FactoryStore) FindByCode(code string) (*Factory, error) {
	return store.FindOne(db.Cond{"code": code})
}

func (store FactoryStore) FindOne(cond ...interface{}) (*Factory, error) {
	var factory *Factory
	if err := store.Find(cond...).One(&factory); err != nil {
		return nil, err
	}
	return factory, nil
}

// we use postgres built-in full text search here to approximate matching
// factory name based on some search query.

type FactoryWithRank struct {
	*Factory
	Rank float64 `db:"rank"`
}

func (store FactoryStore) SearchByName(name string) ([]*FactoryWithRank, error) {
	withQuery := db.Raw("*, to_tsvector(name) txtsearch, plainto_tsquery(?) query", name)
	innerSelect := DB.Select(withQuery).From("factories")

	withRank := db.Raw("t.*, ts_rank(txtsearch, query) as rank")
	query := DB.Select(withRank).From(innerSelect).As("t").OrderBy("rank desc").Limit(10)

	var rankings []*FactoryWithRank
	if err := query.All(&rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}
