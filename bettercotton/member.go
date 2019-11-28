package bettercotton

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/pxue/sustainable/data"
)

type member struct {
	*data.BCI
	Since string `json:"since"`
}

func Load() {
	file, err := os.Open("./bettercotton/members.json")
	if err != nil {
		log.Fatal(err)
	}

	var members map[string]*member
	if err := json.NewDecoder(file).Decode(&members); err != nil {
		log.Fatal(err)
	}

	for _, v := range members {
		v.BCI.Since, _ = time.Parse("January 2006", v.Since)
		if err := data.DB.BCI.Save(v); err != nil {
			log.Fatal(err)
		}
	}
}
