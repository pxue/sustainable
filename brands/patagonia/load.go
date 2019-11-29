package patagonia

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pxue/sustainable/data"
	"upper.io/db.v3"
)

func loadProducts(cache map[string][]string) {
	file, err := os.Open("./brands/patagonia/dump/mens_products.json")
	if err != nil {
		log.Fatal(err)
	}

	var products []product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		log.Fatal(err)
	}
	file.Close()

	for _, p := range products {
		dbProduct := &data.Product{
			Code:     p.ID,
			Name:     p.Name,
			Category: fmt.Sprintf("%s %s %s", p.Category, p.ProdClass, p.ProdCat),
			Price:    p.Price,
			Brand:    "Patagonia",
		}
		if err := data.DB.Product.Save(dbProduct); err != nil {
			log.Fatal(err)
		}

		for _, code := range cache[p.ID] {
			factory, err := data.DB.Factory.FindByCode(code)
			if err != nil {
				log.Printf("factory code not found %s", code)
				continue
			}

			_, err = data.DB.InsertInto("product_suppliers").
				Columns("product_id", "factory_id").
				Values(dbProduct.ID, factory.ID).
				Exec()
			if err != nil {
				log.Fatalf("failed to insert: %v", err)
			}
		}
	}

}

func loadFactories() {
	// load factories
	file, err := os.Open("./brands/patagonia/dump/factories.json")
	if err != nil {
		log.Fatal(err)
	}

	var factories []*factory
	if err := json.NewDecoder(file).Decode(&factories); err != nil {
		log.Fatal(err)
	}
	file.Close()

	for _, fact := range factories {
		name := strings.ReplaceAll(strings.ToLower(fact.Name), "co.", "")
		name = strings.ReplaceAll(strings.ToLower(name), "s.a. de c.v.", "")
		name = strings.ReplaceAll(strings.ToLower(name), "ltd", "")
		name = strings.ReplaceAll(strings.ToLower(name), "joint stock company", "")

		//hasFound := false

		factory, _ := data.DB.Factory.FindOne(db.Cond{"code": fact.Code})
		if factory != nil {
			continue
		}
		//for _, s := range factories {
		//if s.Rank > 0.2 {
		//log.Printf("%s maybe found! %s (%.2f)", name, s.Name, s.Rank)
		//hasFound = true
		//}
		//}

		//suppliers, _ := data.DB.Supplier.SearchByName(name)
		//for _, s := range suppliers {
		//if s.Rank > 0.2 {
		//log.Printf("%s maybe found! %s (%.2f)", name, s.Name, s.Rank)
		//hasFound = true
		//}
		//}

		//if hasFound {
		//// manually add
		//continue
		//}

		toSave := &data.Factory{
			Code:    fact.Code,
			Name:    fact.Name,
			Address: fact.Address,
			Country: fact.Country,
		}
		if err := data.DB.Factory.Save(toSave); err != nil {
			log.Println(err)
		}
	}
}

func Load() {
	pfFile, err := os.Open("./brands/patagonia/dump/product_factories.json")
	if err != nil {
		log.Fatal(err)
	}

	var pfs []map[string]string
	if err := json.NewDecoder(pfFile).Decode(&pfs); err != nil {
		log.Fatal(err)
	}
	pfFile.Close()

	pfsLookup := map[string][]string{}
	for _, v := range pfs {
		pfsLookup[v["p"]] = append(pfsLookup[v["p"]], v["f"])
	}

	loadProducts(pfsLookup)
}
