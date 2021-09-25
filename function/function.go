package function

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dsychin/ohamame-shipping-label/label"
	"github.com/gocarina/gocsv"
	"github.com/signintech/gopdf"
)

func ShippingLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("GET /")

		tmpl := template.Must(template.ParseFiles("./template/index.gohtml"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		log.Println("POST /")

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(handler.Filename, ".csv") {
			log.Println("File is not CSV")
			http.Error(w, "Please upload csv file", http.StatusBadRequest)
			return
		}

		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contents := []label.LabelContent{}
		err = gocsv.UnmarshalBytes(fileContent, &contents)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		l := label.NewLabel(2, 4, *gopdf.PageSizeA4, 10)
		err = l.CreateShippingLabelPdf(w, contents)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=shipping_label.pdf")
	} else {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
}
