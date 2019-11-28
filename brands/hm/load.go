package hm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pxue/sustainable/data"
	"upper.io/db.v3"
)

// product
// articlesList -> compositions -> materials
//  -> mainCategory.name
// suppliers

type product struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	MainCategory struct {
		Name string `json:"name"`
	} `json:"mainCategory"`
	ArticlesList []*variant `json:"articlesList"`
	WhitePrice   *price     `json:"whitePrice"`
}

type price struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type variant struct {
	Compositions []*composition `json:"compositions"`
}

type composition struct {
	Materials []struct {
		Name       string `json:"name"`
		Percentage string `json:"percentage"`
	} `json:"materials"`
	CompositionType string `json:"compositionType,omitempty"`
}

type productSuppliers struct {
	Countries []*country `json:"countries"`
}

type country struct {
	Name      string      `json:"name"`
	Suppliers []*supplier `json:"suppliers"`
}

type supplier struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Factories []*supplierFactory `json:"factories"`
}

type supplierFactory struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	WorkersNumber string `json:"workersNumber"`
}

func DoSupplier() {
	// parse all the products.
	// read in the dump
	files, err := ioutil.ReadDir("./dump/hmsuppliers")
	if err != nil {
		log.Fatal(err)
	}
	for _, fnfo := range files {
		f, _ := os.Open("./dump/hmsuppliers/" + fnfo.Name())
		code := strings.TrimSuffix(fnfo.Name(), ".json")

		var product *data.Product
		if err := data.DB.Product.Find(db.Cond{"code": code}).One(&product); err != nil {
			log.Println(err)
			continue
		}

		var s productSuppliers
		if err := json.NewDecoder(f).Decode(&s); err != nil {
			log.Fatal(err)
		}

		for _, country := range s.Countries {
			for _, sup := range country.Suppliers {

				// check if supplier exists
				supplier, err := data.DB.Supplier.FindByHMID(sup.ID)
				if err != nil {
					if err == db.ErrNoMoreRows {
						supplier = &data.Supplier{
							Name:    sup.Name,
							Country: country.Name,
							HmID:    sup.ID,
						}
						if err := data.DB.Supplier.Save(supplier); err != nil {
							log.Panicf("save supplier %+v error: %v", supplier, err)
							continue
						}
					} else {
						log.Panicf("get supplier error: %v", err)
						continue
					}
				}

				for _, fac := range sup.Factories {
					// check if factory exists
					factory, err := data.DB.Factory.FindByHMID(fac.ID)
					if err != nil {
						if err == db.ErrNoMoreRows {
							factory = &data.Factory{
								SupplierID: supplier.ID,
								HmID:       fac.ID,
								Name:       fac.Name,
								Country:    country.Name,
								Address:    fac.Address,
							}
							if err := data.DB.Factory.Save(factory); err != nil {
								log.Panicf("save factory %+v error: %v", factory, err)
								continue
							}
						} else {
							log.Panicf("get factory error: %v", err)
							continue
						}
					}

					// insert into product suppliers
					_, err = data.DB.InsertInto("product_suppliers").
						Columns("product_id", "supplier_id", "factory_id").
						Values(product.ID, supplier.ID, factory.ID).
						Exec()
					if err != nil {
						log.Printf("failed to insert: %v", err)
					}
				}

			}
		}
	}

}

func DoProduct() {
	// parse all the products.
	// read in the dump
	files, err := ioutil.ReadDir("./brands/hm/dump/hmmens")
	if err != nil {
		log.Fatal(err)
	}
	for _, fnfo := range files {
		f, _ := os.Open("./brands/hm/dump/hmmens/" + fnfo.Name())

		var p product
		if err := json.NewDecoder(f).Decode(&p); err != nil {
			log.Fatal(err)
		}

		product := &data.Product{
			Name:      p.Name,
			Code:      p.Code,
			Category:  p.MainCategory.Name,
			Price:     p.WhitePrice.Price,
			Materials: data.ProductMaterial{},
		}

		for _, variant := range p.ArticlesList {
			for _, compo := range variant.Compositions {
				for _, m := range compo.Materials {
					product.Materials[m.Name], _ = strconv.ParseFloat(m.Percentage, 64)
				}
			}
			break
		}

		if err := data.DB.Product.Save(product); err != nil {
			log.Printf("db error: %v", err)
		}
	}
}
