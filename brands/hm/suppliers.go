package hm

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func supplierWorker(id int, chann chan string, signal chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	log.Printf("starting worker %d", id)
	signal <- 1 // signal we're ready.

	provider, _ := New(false)

	processed := 0
	timeSpent := 0
	for code := range chann {
		start := time.Now()

		ctx := context.Background()
		items, _, err := provider.GetSupplier(ctx, code)
		if err != nil {
			log.Printf("get supplier errored %s: %v", code, err)
			continue
		}
		f, _ := os.Create(fmt.Sprintf("./hmsuppliers/%s.json", code))
		if err := json.NewEncoder(f).Encode(items); err != nil {
			log.Printf("encode errored %s: %v", code, err)
		}
		f.Close()

		processed++
		duration := time.Since(start)
		timeSpent += int(duration)
		log.Printf("worker %d processed product %s. Ts: %v", id, code, duration)

		// artificial sleep.
		time.Sleep(time.Duration(500+rand.Intn(1500)) * time.Millisecond)
	}

	if processed > 0 {
		log.Printf("worker %d exiting. processed %d. avgTs: %v", id, processed, time.Duration(int(timeSpent)/processed))
	}
}

func getParsed() map[string]bool {
	ret := map[string]bool{}
	files, err := ioutil.ReadDir("./hmsuppliers")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filename := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "product_"), ".json")
		ret[filename] = true
	}
	return ret
}

func Suppliers() {
	numWorkers := 20

	var wg sync.WaitGroup
	chann := make(chan string, numWorkers)
	signal := make(chan int)
	for i := 1; i <= numWorkers; i++ {
		go supplierWorker(i, chann, signal, &wg)
		<-signal // wait for worker to launch
	}

	exists := getParsed()

	files, err := ioutil.ReadDir("./hmmens")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		filename := strings.TrimSuffix(strings.TrimPrefix(f.Name(), "product_"), ".json")
		if exists[filename] {
			continue
		}
		chann <- filename
	}
	close(chann)

	wg.Wait()
	log.Println("cleaning up.")
}
