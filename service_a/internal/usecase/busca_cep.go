package usecase

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/jorgemarinho/go-open-telemetry/service_a/internal/dto"
)

const (
	serviceBCepURL = "http://localhost:8080/clima?cep=%s"
)

type BuscaCepUseCase struct {
	BuscaCepInputDTO dto.BuscaCepInputDTO
}

func NewBuscaCepUseCase(buscaCepInputDTO dto.BuscaCepInputDTO) *BuscaCepUseCase {
	return &BuscaCepUseCase{
		BuscaCepInputDTO: buscaCepInputDTO,
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)

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
