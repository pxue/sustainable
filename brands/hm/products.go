package hm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// worker to fetch and cache products.
func productWorker(id int, chann chan string, signal chan int, wg *sync.WaitGroup) {
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
		item, _, err := provider.GetProduct(ctx, code)
		if err != nil {
			log.Printf("fetch product: %v", err)
			continue
		}
		f, err := os.Create(fmt.Sprintf("./brands/hm/dump/hmmens/product_%s.json", code))
		if err != nil {
			log.Printf("file create: %v", err)
			continue
		}
		if err := json.NewEncoder(f).Encode(item); err != nil {
			log.Printf("json encode %v", err)
		}
		f.Close()

		processed++
		duration := time.Since(start)
		timeSpent += int(duration)
		log.Printf("worker %d processed item %s. Ts: %v", id, code, duration)

		// artificial sleep.
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	if processed > 0 {
		log.Printf("worker %d exiting. processed %d. avgTs: %v", id, processed, time.Duration(int(timeSpent)/processed))
	}
}

func Products() {
	numWorkers := 20

	var wg sync.WaitGroup
	chann := make(chan string, numWorkers)
	signal := make(chan int)
	for i := 1; i <= numWorkers; i++ {
		go productWorker(i, chann, signal, &wg)
		<-signal // wait for worker to launch
	}

	// read in the dump
	f, _ := os.Open("./brands/hm/dump/hm_mens.json")

	type article struct {
		Code string `json:"articleCode"`
	}
	var products []*article
	if err := json.NewDecoder(f).Decode(&products); err != nil {
		log.Panic(err)
	}
	f.Close()

	for _, p := range products {
		chann <- p.Code
	}
	close(chann)

	wg.Wait()
	log.Println("cleaning up.")
}
