package hm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func shopWorker(id int, chann chan int, signal chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	log.Printf("starting worker %d", id)
	signal <- 1 // signal we're ready.

	provider, _ := New(false)

	processed := 0
	timeSpent := 0
	for offset := range chann {
		start := time.Now()

		ctx := context.Background()
		items, _, err := provider.ShopByProduct(ctx, 36, offset)
		if err != nil {
			log.Println(err)
			continue
		}
		f, _ := os.Create(fmt.Sprintf("./hmmens/items_%d.json", offset/36))
		if err := json.NewEncoder(f).Encode(items); err != nil {
			log.Println(err)
		}
		f.Close()

		processed++
		duration := time.Since(start)
		timeSpent += int(duration)
		log.Printf("worker %d processed offset %d. Ts: %v", id, offset, duration)
	}

	if processed > 0 {
		log.Printf("worker %d exiting. processed %d. avgTs: %v", id, processed, time.Duration(int(timeSpent)/processed))
	}
}

func shopAll() {
	numWorkers := 10

	var wg sync.WaitGroup
	chann := make(chan int, numWorkers)
	signal := make(chan int)
	for i := 1; i <= numWorkers; i++ {
		go shopWorker(i, chann, signal, &wg)
		<-signal // wait for worker to launch
	}

	for page := 0; page < 100; page++ {
		chann <- page * 36
	}
	close(chann)

	wg.Wait()
	log.Println("cleaning up.")
}
