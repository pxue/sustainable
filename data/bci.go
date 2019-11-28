package data

import (
	"time"

	"upper.io/bond"
	"upper.io/db.v3"
)

type BCI struct {
	ID       int64     `db:"id,omitempty" json:"id"`
	Name     string    `db:"name" json:"name"`
	Category string    `db:"category" json:"category"`
	Country  string    `db:"country" json:"country"`
	Website  string    `db:"website" json:"website"`
	Since    time.Time `db:"since" json:"since"`
}

type BCIStore struct {
	bond.Store
}

func (*BCI) CollectionName() string {
	return `bci_members`
}

func (store BCIStore) FindByName(name string) (*BCI, error) {
	return store.FindOne(db.Cond{"name": name})
}

func (store BCIStore) FindOne(cond db.Cond) (*BCI, error) {
	var bci *BCI
	if err := store.Find(cond).One(&bci); err != nil {
		return nil, err
	}
	return bci, nil
}
