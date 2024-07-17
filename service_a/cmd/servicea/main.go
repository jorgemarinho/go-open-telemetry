package main

import (
	"net/http"

	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/infra/web"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to service a"))
	})
	http.HandleFunc("/busca/cidade", web.BuscaCepHandler)
	http.ListenAndServe(":8082", nil)
}
