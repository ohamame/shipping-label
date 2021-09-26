package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dsychin/ohamame-shipping-label/function"
	"github.com/dsychin/ohamame-shipping-label/label"
	"github.com/gocarina/gocsv"
	"github.com/signintech/gopdf"
)

var cli bool
var port int

func init() {
	flag.BoolVar(&cli, "cli", false, "Run as CLI")
	flag.IntVar(&port, "port", 8080, "Port to run web server at")

}

func main() {
	flag.Parse()

	if cli {
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

		l := label.NewLabel(2, 4, *gopdf.PageSizeA4, 10, false)
		err = l.CreateShippingLabelPdf(f, contents)
		if err != nil {
			panic(err)
		}
	} else {
		http.HandleFunc("/", function.ShippingLabel)

		log.Println(fmt.Sprintf("Listening on port %d", port))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
