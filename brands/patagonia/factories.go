package patagonia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type factory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Address     string `json:"address"`
	Country     string `json:"country"`
	Description string `json:"description"`
	Since       string `json:"since"`
	Type        string `json:"type"`
	ProductType string `json:"product_type"`
	GenderMix   struct {
		Male   string `json:"male"`
		Female string `json:"female"`
	} `json:"genderMix"`
	NumWorkers string `json:"num_workers"`
	Coord      struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	} `json:"coord"`
}

// worker to fetch and cache products.
func factoryWorker(id int, chann chan string, signal chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	log.Printf("starting worker %d", id)
	signal <- 1 // signal we're ready.

	processed := 0
	timeSpent := 0
	for ID := range chann {
		start := time.Now()

		resp, err := http.DefaultClient.Get(fmt.Sprintf("https://www.patagonia.com/on/demandware.store/Sites-patagonia-us-Site/en_US/WhereToGetIt-GetFactories?pid=%s", ID))
		if err != nil {
			log.Printf("fetch factory: %v", err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("fetch factory: %v", err)
			continue
		}

		var factories []factory
		doc.Find("div.factory").Each(func(i int, s *goquery.Selection) {
			var f factory
			f.ID, _ = s.Attr("data-clientkey")
			f.Name = s.Find("span").Text()
			f.Description = s.Find("p").Text()
			factories = append(factories, f)
		})

		if len(factories) == 0 {
			log.Printf("worker %d: product %s resulted in no factories.", id, ID)
			continue
		}

		out, err := os.Create(fmt.Sprintf("./brands/patagonia/dump/test/product_%s.json", ID))
		if err != nil {
			log.Printf("file create: %v", err)
			return
		}
		if err := json.NewEncoder(out).Encode(&factories); err != nil {
			log.Printf("json encode %v", err)
		}
		out.Close()

		processed++
		duration := time.Since(start)
		timeSpent += int(duration)
		log.Printf("worker %d processed product %s. Ts: %v", id, ID, duration)

		// artificial sleep.
		time.Sleep(time.Duration(500+rand.Intn(1500)) * time.Millisecond)
	}

	if processed > 0 {
		log.Printf("worker %d exiting. processed %d. avgTs: %v", id, processed, time.Duration(int(timeSpent)/processed))
	}
}

func getParsed() map[string]bool {
	ret := map[string]bool{}
	files, err := ioutil.ReadDir("./brands/patagonia/dump/test/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filename := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "product_"), ".json")
		ret[filename] = true
	}
	return ret
}

func Factories() {
	numWorkers := 10

	var wg sync.WaitGroup
	chann := make(chan string, numWorkers)
	signal := make(chan int)
	for i := 1; i <= numWorkers; i++ {
		go factoryWorker(i, chann, signal, &wg)
		<-signal // wait for worker to launch
	}

	file, err := os.Open("./brands/patagonia/dump/mens_products.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var products []product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		log.Panic(err)
	}

	exists := getParsed()
	for _, p := range products {
		if exists[p.ID] {
			continue
		}
		chann <- p.ID
	}
	close(chann)

	wg.Wait()
	log.Println("cleaning up.")
}
