package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/pxue/sustainable/data"
	"github.com/pxue/sustainable/lib/geocode"
	"github.com/pxue/sustainable/web"
	"upper.io/db.v3"
)

func GetFactoryCoords() {
	geo, err := geocode.New("")
	if err != nil {
		log.Fatal(err)
	}

	var factories []*data.Factory
	data.DB.Factory.Find(db.Cond{"lat": 0}).
		OrderBy("id").All(&factories)

	for _, f := range factories {
		log.Printf("processing: %s", f.Address)

		coord, err := geo.AddressToCoord(f.Address)
		if err != nil {
			if err == geocode.ErrNoResult {
				log.Println("no result")
				continue
			}
			log.Panic(err)
		}

		f.Lat = coord[0]
		f.Lon = coord[1]

		if err := data.DB.Factory.Save(f); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dbConf := data.DBConf{
		Database:        "sustainable",
		Hosts:           []string{"localhost"},
		Username:        "sustainable",
		ApplicationName: "test",
		DebugQueries:    true,
	}
	if _, err := data.NewDB(dbConf); err != nil {
		log.Fatal(errors.Wrap(err, "database: connection failed"))
	}

	log.Panic(web.New())
}
