package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alvinfeng/mosaic/config"
	"github.com/alvinfeng/mosaic/imgserver"
)

func main() {
	yaml := "sample-config.yaml"

	if len(os.Args) > 1 {
		yaml = os.Args[1]
	}

	config, err := config.LoadConfig(yaml)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	r, err := imgserver.New(config)
	if err != nil {
		log.Fatal("Error initializing image server: ", err)
	}

	log.Fatal(http.ListenAndServe(":8081", r.Router))
}
