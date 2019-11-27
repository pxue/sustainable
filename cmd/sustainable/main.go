package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	dbConf := data.DBConf{
		Database:        "sustainable",
		Hosts:           []string{"localhost"},
		Username:        "sustainable",
		ApplicationName: "test",
		DebugQueries:    false,
	}
	if _, err := data.NewDB(dbConf); err != nil {
		log.Fatal(errors.Wrap(err, "database: connection failed"))
	}
}
