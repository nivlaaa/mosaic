package main

import (
	"log"
	"net/http"

	"github.com/alvinfeng/mosaic/imgserver"
)

func main() {
	r, err := imgserver.New()
	if err != nil {
		log.Fatal("Error initializing image server: ", err)
	}

	log.Fatal(http.ListenAndServe(":8080", r.Router))
}
