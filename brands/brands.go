package brands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pxue/sustainable/data"
)

func Load() {
	f, err := os.Open("./dumps/brands.json")
	if err != nil {
		log.Fatal(err)
	}
	var brands []*data.Brand
	if err := json.NewDecoder(f).Decode(&brands); err != nil {
		log.Fatal(err)
	}

	for _, b := range brands {
		if err := data.DB.Brand.Save(b); err != nil {
			log.Println(err)
		}
	}
}
