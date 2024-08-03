package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/jorgemarinho/go-open-telemetry/service_b/internal/dto"
	"github.com/jorgemarinho/go-open-telemetry/service_b/internal/entity"
	"github.com/jorgemarinho/go-open-telemetry/service_b/internal/errors"
	"github.com/jorgemarinho/go-open-telemetry/service_b/internal/infra/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	viaCepURL     = "https://viacep.com.br/ws/%s/json/"
	weatherAPI    = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s"
	weatherAPIKey = "8887ae192b2343f9a32114928240104"
)

type BuscaCepUseCase struct {
	BuscaCepInputDTO dto.BuscaCepInputDTO
	Ctx              context.Context
	h                web.TemplateData
}

func NewBuscaCepUseCase(buscaCepInputDTO dto.BuscaCepInputDTO, ctx context.Context, h web.TemplateData) *BuscaCepUseCase {
	return &BuscaCepUseCase{
		BuscaCepInputDTO: buscaCepInputDTO,
		Ctx:              ctx,
		h:                h,
	}
}

func (b BuscaCepUseCase) Execute() (dto.BuscaCepOutputDTO, error) {
	if !isValidCep(b.BuscaCepInputDTO.Cep) {
		return dto.BuscaCepOutputDTO{}, &errors.HTTPError{Code: http.StatusUnprocessableEntity, Message: "CEP must have 8 digits and only contain numbers"}
	}

	ctx, spanCEP := b.h.OTELTracer.Start(b.Ctx, "Service B: Consulta cep")

	cep, err := b.BuscaCep(b.BuscaCepInputDTO.Cep, ctx)
	if err != nil {
		return dto.BuscaCepOutputDTO{}, err
	}

	spanCEP.End()

	ctx, spanTemp := b.h.OTELTracer.Start(b.Ctx, "Service B: Consulta temperatura")

	temperatura, err := b.BuscaTemperatura(cep.Localidade, ctx)
	if err != nil {
		return dto.BuscaCepOutputDTO{}, err
	}

	spanTemp.End()

	return dto.BuscaCepOutputDTO{
		City:  cep.Localidade,
		TempC: temperatura.TempC,
		TempF: getTemperatureFahrenheit(temperatura.TempF),
		TempK: getTemperatureKelvin(temperatura.TempK),
	}, nil
}

func (b BuscaCepUseCase) BuscaCep(cep string, ctx context.Context) (*entity.Cep, error) {
	url := fmt.Sprintf(viaCepURL, cep)
	return b.makeHTTPRequestCep(url, ctx)
}

func (b BuscaCepUseCase) BuscaTemperatura(nomeCidade string, ctx context.Context) (*entity.Temperatura, error) {
	encodedNomeCidade := url.QueryEscape(nomeCidade)
	url := fmt.Sprintf(weatherAPI, weatherAPIKey, encodedNomeCidade)
	return b.makeHTTPRequestTemperatura(url, ctx)
}

func (b BuscaCepUseCase) makeHTTPRequestCep(url string, ctx context.Context) (*entity.Cep, error) {

	req, err := http.NewRequestWithContext(b.Ctx, "GET", url, nil)

	if err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error making HTTP request"}
	}

	// Injetando o header do request id. Necessário para realizar o tracker
	otel.GetTextMapPropagator().Inject(b.Ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error making HTTP request"}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error reading response body"}
	}

	var tempResult map[string]interface{}
	if err := json.Unmarshal(body, &tempResult); err == nil {
		if errVal, ok := tempResult["erro"]; ok && errVal == "true" {
			return nil, &errors.HTTPError{Code: http.StatusNotFound, Message: "can not find zipcode"}
		}
	}

	var result entity.Cep
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error unmarshalling cep response"}
	}

	return &result, nil
}

func (b BuscaCepUseCase) makeHTTPRequestTemperatura(url string, ctx context.Context) (*entity.Temperatura, error) {

	req, err := http.NewRequestWithContext(b.Ctx, "GET", url, nil)

	if err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error making HTTP request"}
	}

	// Injetando o header do request id. Necessário para realizar o tracker
	otel.GetTextMapPropagator().Inject(b.Ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error making HTTP request"}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error reading response body"}
	}

	var response struct {
		Current entity.Temperatura `json:"current"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, &errors.HTTPError{Code: http.StatusInternalServerError, Message: "error unmarshalling temperatura response"}
	}

	return &response.Current, nil
}

func isValidCep(cep string) bool {
	return regexp.MustCompile(`^\d{8}$`).MatchString(cep)
}

func getTemperatureFahrenheit(celsius float64) float64 {
	return (celsius * 1.8) + 32
}

func getTemperatureKelvin(celsius float64) float64 {
	return celsius + 273.15
}
