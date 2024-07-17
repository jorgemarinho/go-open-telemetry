package main

import (
	"net/http"

	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/infra/web"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	http.HandleFunc("/cep", web.BuscaCepHandler)
	http.ListenAndServe(":8082", nil)
}
