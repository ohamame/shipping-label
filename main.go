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

	// create test data
	contents := []label.LabelContent{
		{
			Address: "Lorem ipsum dolor sit amet",
		},
		{
			Address: "Lorem ipsum dolor sit amet 2",
		},
		{
			Address: "Lorem ipsum dolor sit amet 3",
		},
		{
			Address: "Lorem ipsum dolor sit amet 4",
		},
		{
			Address: "Lorem ipsum dolor sit amet 5",
		},
		{
			Address: "Lorem ipsum dolor sit amet 6",
		},
	}

	l := label.NewLabel(2, 6, *gopdf.PageSizeA4, 10)
	err = l.CreateShippingLabelPdf(f, contents)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
