package main

import (
	"net/http"

	"github.com/jorgemarinho/go-open-telemetry/service_b/internal/infra/web"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to service b"))
	})
	http.HandleFunc("/clima", web.BuscaCepHandler)
	http.ListenAndServe(":8080", nil)
}
