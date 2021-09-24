package main

import (
	"log"
	"os"

	"github.com/dsychin/ohamame-shipping-label/label"
)

func main() {
	f, err := os.Create("output.pdf")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer f.Close()

	err = label.CreateShippingLabelPdf(f)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
