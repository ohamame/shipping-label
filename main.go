package main

import (
	"log"
	"os"

	"github.com/dsychin/ohamame-shipping-label/label"
	"github.com/signintech/gopdf"
)

func main() {
	f, err := os.Create("output.pdf")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer f.Close()

	l := label.NewLabel(2, 6, *gopdf.PageSizeA4, 10)
	err = l.CreateShippingLabelPdf(f)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
