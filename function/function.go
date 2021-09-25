package function

import (
	"html/template"
	"log"
	"net/http"
)

func ShippingLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("/ handler")

		tmpl := template.Must(template.ParseFiles("./template/index.gohtml"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}
