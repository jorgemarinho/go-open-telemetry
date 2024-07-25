package web

import (
	"encoding/json"
	"net/http"

	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/dto"
	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/errors"
	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func BuscaCepHandler(w http.ResponseWriter, r *http.Request, h *TemplateData) {

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	// Criação de span inicial
	ctx, span := h.OTELTracer.Start(ctx, "Início Processamento "+h.RequestNameOTEL)
	defer span.End()

	cepParam := r.FormValue("cep")

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
