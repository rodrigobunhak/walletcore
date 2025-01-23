package web

import (
	"encoding/json"
	"net/http"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/create_account"
)

type WebAccountHandler struct {
	CreateAccountUseCase create_account.CreateAccountUseCase
}

func NewWebAccountHandler(createAccountUseCase create_account.CreateAccountUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (w WebAccountHandler) CreateAccount(res http.ResponseWriter, req *http.Request) {
	var dto create_account.CreateAccountInputDTO
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := w.CreateAccountUseCase.Execute(dto)
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