package patagonia

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	ProdClass   string  `json:"prod_class"`
	ProdCat     string  `json:"prod_cat"`
	ProdTeam    string  `json:"prod_team"`
	Price       float64 `json:"price"`
}

// worker to fetch and cache products.
func productWorker(id int, chann chan int, signal chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	log.Printf("starting worker %d", id)
	signal <- 1 // signal we're ready.

	processed := 0
	timeSpent := 0
	for i := range chann {
		start := time.Now()

		offset := i * 20
		resp, err := http.DefaultClient.Get(fmt.Sprintf("https://www.patagonia.com/shop/mens?sz=20&start=%d&format=page-element", offset))
		if err != nil {
			log.Printf("fetch product: %v", err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("fetch product: %v", err)
			continue
		}

		var products []product
		doc.Find("li").Each(func(i int, s *goquery.Selection) {
			if _, exists := s.Attr("id"); !exists {
				return
			}

			dataAttr, _ := s.Find("div.product-tile").Attr("data-tealium")
			var p product
			if err := json.Unmarshal([]byte(dataAttr), &p); err != nil {
				log.Printf("json decode %v", err)
				return
			}
			p.Price, _ = strconv.ParseFloat(strings.TrimPrefix(s.Find("span.price-sales").Text(), "$"), 64)

			products = append(products, p)
		})

		if len(products) == 0 {
			log.Printf("worker %d: page %d resulted in no product.", id, i)
			continue
		}

		out, err := os.Create(fmt.Sprintf("./brands/patagonia/dump/mproducts/page_%d.json", i))
		if err != nil {
			log.Printf("file create: %v", err)
			return
		}
		if err := json.NewEncoder(out).Encode(&products); err != nil {
			log.Printf("json encode %v", err)
		}
		out.Close()

		processed++
		duration := time.Since(start)
		timeSpent += int(duration)
		log.Printf("worker %d processed page %d. Ts: %v", id, i, duration)

		// artificial sleep.
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	if processed > 0 {
		log.Printf("worker %d exiting. processed %d. avgTs: %v", id, processed, time.Duration(int(timeSpent)/processed))
	}
}

func Products() {
	numWorkers := 10

	var wg sync.WaitGroup
	chann := make(chan int, numWorkers)
	signal := make(chan int)
	for i := 1; i <= numWorkers; i++ {
		go productWorker(i, chann, signal, &wg)
		<-signal // wait for worker to launch
	}

	for i := 0; i < 30; i++ {
		chann <- i
	}
	close(chann)

	wg.Wait()
	log.Println("cleaning up.")
}
