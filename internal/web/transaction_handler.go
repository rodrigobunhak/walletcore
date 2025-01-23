package web

import (
	"encoding/json"
	"net/http"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransactionUseCase create_transaction.CreateTransactionUseCase) *WebTransactionHandler {
	return &WebTransactionHandler{
		CreateTransactionUseCase: createTransactionUseCase,
	}
}

func (w WebTransactionHandler) CreateTransaction(res http.ResponseWriter, req *http.Request) {
	var dto create_transaction.CreateTransactionInputDTO
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	output, err := w.CreateTransactionUseCase.Execute(ctx, dto)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
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