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

	if config.Insecure {
		log.Println("Warning you are running an insecure image server")
		log.Fatal(http.ListenAndServe(config.Addr, r.Router))
	} else {
		log.Fatal(http.ListenAndServeTLS(config.Addr, config.CertPath, config.KeyPath, r.Router))
	}
}
