package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/dto"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var (
	serviceBCepURL string
)

func init() {
	serviceBCepURL = viper.GetString("URL_SERVICE_B") + "/clima?cep=%s"
}

type BuscaCepUseCase struct {
	BuscaCepInputDTO dto.BuscaCepInputDTO
	Ctx              context.Context
}

func NewBuscaCepUseCase(buscaCepInputDTO dto.BuscaCepInputDTO, ctx context.Context) *BuscaCepUseCase {
	return &BuscaCepUseCase{
		BuscaCepInputDTO: buscaCepInputDTO,
		Ctx:              ctx,
	}
}

func (b BuscaCepUseCase) Execute() (dto.BuscaCepOutputDTO, error) {

	if !isValidCep(b.BuscaCepInputDTO.Cep) {
		return dto.BuscaCepOutputDTO{}, fmt.Errorf("CEP must have 8 digits and only contain numbers")
	}

	cepDto, err := b.BuscaCep(b.BuscaCepInputDTO.Cep)
	if err != nil {
		return dto.BuscaCepOutputDTO{}, err
	}

	return cepDto, nil
}

func (b BuscaCepUseCase) BuscaCep(cep string) (dto.BuscaCepOutputDTO, error) {
	url := fmt.Sprintf(serviceBCepURL, cep)
	return b.makeHTTPRequestCep(url)
}

func (b BuscaCepUseCase) makeHTTPRequestCep(url string) (dto.BuscaCepOutputDTO, error) {

	req, err := http.NewRequestWithContext(b.Ctx, "GET", url, nil)

	if err != nil {
		return dto.BuscaCepOutputDTO{}, fmt.Errorf("error making HTTP request: %w", err)
	}

	// Injetando o header do request id. Necessário para realizar o tracker
	otel.GetTextMapPropagator().Inject(b.Ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return dto.BuscaCepOutputDTO{}, fmt.Errorf("error making HTTP request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return dto.BuscaCepOutputDTO{}, fmt.Errorf("error reading response body: %w", err)
	}

	var result dto.BuscaCepOutputDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return dto.BuscaCepOutputDTO{}, fmt.Errorf("error unmarshalling cep response: %w", err)
	}

	return result, nil
}

func isValidCep(cep string) bool {
	return regexp.MustCompile(`^\d{8}$`).MatchString(cep)
}
