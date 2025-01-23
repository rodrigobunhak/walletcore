package web

import (
	"encoding/json"
	"net/http"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/create_client"
)

type WebClientHandler struct {
	CreateClientUseCase create_client.CreateClientUseCase
}

func NewWebClientHandler(createClientUseCase create_client.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase: createClientUseCase,
	}
}

func (w WebClientHandler) CreateClient(res http.ResponseWriter, req *http.Request) {
	var dto create_client.CreateClientInputDTO
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := w.CreateClientUseCase.Execute(dto)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(output)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
}