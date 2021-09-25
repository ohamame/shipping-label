package main

import (
	"os"

	"github.com/dsychin/ohamame-shipping-label/label"
	"github.com/gocarina/gocsv"
	"github.com/signintech/gopdf"
)

func main() {
	f, err := os.Create("output.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// create test data
	contentsFile, err := os.ReadFile("input.csv")
	if err != nil {
		panic(err)
	}
	contents := []label.LabelContent{}

	err = gocsv.UnmarshalBytes(contentsFile, &contents)
	if err != nil {
		panic(err)
	}

	l := label.NewLabel(2, 4, *gopdf.PageSizeA4, 10)
	err = l.CreateShippingLabelPdf(f, contents)
	if err != nil {
		panic(err)
	}
}
