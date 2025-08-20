package web

import (
	"encoding/json"
	"net/http"

	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/dto"
	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/errors"
	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de requisições HTTP recebidas.",
	},
	[]string{"service", "endpoint", "method"},
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}
func BuscaCepHandler(w http.ResponseWriter, r *http.Request, h *TemplateData) {
	httpRequestsTotal.WithLabelValues("service_a", "/busca/cidade", r.Method).Inc()

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	var cepParam string

	cepParam = r.FormValue("cep")

	if cepParam == "" && r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var requestBody map[string]string
		err := decoder.Decode(&requestBody)
		if err != nil {

			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		cepParam = requestBody["cep"]
	}

	if cepParam == "" {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if len(cepParam) < 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ctx, spanServiceB := h.OTELTracer.Start(ctx, "Consulta service-b")
	defer spanServiceB.End()

	buscaCepInputDTO := dto.BuscaCepInputDTO{Cep: cepParam}

	newBuscaCepUseCase := usecase.NewBuscaCepUseCase(buscaCepInputDTO, ctx)

	cep, err := newBuscaCepUseCase.Execute()

	if err != nil {
		code := http.StatusInternalServerError
		message := err.Error()

		if httpErr, ok := err.(*errors.HTTPError); ok {
			code = httpErr.Code
			message = httpErr.Message
		}

		http.Error(w, message, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cep)
}
