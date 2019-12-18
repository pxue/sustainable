package data

import (
	"fmt"
	"strings"

	"upper.io/bond"
	"upper.io/db.v3/postgresql"
)

type Database struct {
	bond.Session

	Brand    BrandStore
	Product  ProductStore
	Supplier SupplierStore
	Factory  FactoryStore
	Material MaterialStore
	BCI      BCIStore
}

// configuration for db (postgres)
type DBConf struct {
	Database        string   `toml:"database"`
	Hosts           []string `toml:"hosts"`
	Username        string   `toml:"username"`
	Password        string   `toml:"password"`
	DebugQueries    bool     `toml:"debug_queries"`
	ApplicationName string   `toml:"application_name"`
	MaxConnection   int      `tomp:"max_connection"`
}

func (cf *DBConf) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		cf.Username, cf.Password, strings.Join(cf.Hosts, ","), cf.Database)
}

// instance
var DB *Database

func NewDB(conf DBConf) (*Database, error) {
	connURL, err := postgresql.ParseURL(conf.String())
	if err != nil {
		return nil, err
	}
	// extra options
	connURL.Options = map[string]string{
		"application_name": conf.ApplicationName,
	}

	db := &Database{}
	db.Session, err = bond.Open(postgresql.Adapter, connURL)
	if conf.DebugQueries {
		db.Session.SetLogging(true)
	}
	if conf.MaxConnection > 0 {
		db.Session.SetMaxIdleConns(conf.MaxConnection)
	}
	if err != nil {
		return nil, err
	}

	db.Product = ProductStore{db.Store(&Product{})}
	db.Supplier = SupplierStore{db.Store(&Supplier{})}
	db.Factory = FactoryStore{db.Store(&Factory{})}
	db.Material = MaterialStore{db.Store(&Material{})}
	db.BCI = BCIStore{db.Store(&BCI{})}
	db.Brand = BrandStore{db.Store(&Brand{})}

	DB = db
	return db, nil
}
