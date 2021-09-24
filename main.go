package main

import (
	"log"
	"os"

	"github.com/dsychin/ohamame-shipping-label/label"
)

func main() {
	err := label.CreateShippingLabelPdf()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
