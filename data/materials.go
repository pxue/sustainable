package data

import (
	"upper.io/bond"
	"upper.io/db.v3"
)

type Material struct {
	ID    int64  `db:"id,omitempty"`
	Name  string `db:"name"`
	Type  string `db:"type"`
	Score int64  `db:"score"`
}

type MaterialStore struct {
	bond.Store
}

func (*Material) CollectionName() string {
	return `materials`
}

func (store MaterialStore) FindByName(name string) (*Material, error) {
	return store.FindOne(db.Cond{"name": name})
}

func (store MaterialStore) FindOne(cond db.Cond) (*Material, error) {
	var material *Material
	if err := store.Find(cond).One(&material); err != nil {
		return nil, err
	}
	return material, nil
}
