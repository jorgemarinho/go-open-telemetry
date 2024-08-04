package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jorgemarinho/go-open-telemetry/service_b/internal/shared"
)

type Webserver struct {
	TemplateData *shared.TemplateData
}

// NewServer creates a new server instance
func NewServer(templateData *shared.TemplateData) *Webserver {
	return &Webserver{
		TemplateData: templateData,
	}
}

func (we *Webserver) CreateServer() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/clima", func(w http.ResponseWriter, r *http.Request) {
		BuscaCepHandler(w, r, *we.TemplateData)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to service b"))

	})

	return router
}
