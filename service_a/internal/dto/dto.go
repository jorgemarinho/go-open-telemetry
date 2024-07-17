package dto

type BuscaCepInputDTO struct {
	Cep string `json:"cep"`
}

type BuscaCepOutputDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
