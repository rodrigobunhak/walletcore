package create_transaction

import (
	"context"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/entity"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/gateway"
	"github.com/rodrigobunhak/fc-ms-wallet/pkg/events"
	"github.com/rodrigobunhak/fc-ms-wallet/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo string `json:"account_id_to"`
	Amount float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID string `json:"id"`
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo string `json:"account_id_to"`
	Amount float64 `json:"amount"`
}

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	Uow uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated events.EventInterface
}

func NewCreateTransactionUseCase(
	Uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
	) *CreateTransactionUseCase {
		return &CreateTransactionUseCase{
			Uow: Uow,
			EventDispatcher: eventDispatcher,
			TransactionCreated: transactionCreated,
			BalanceUpdated: balanceUpdated,
		}
}

func (useCase *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	balanceUpdatedOutput := &BalanceUpdatedOutputDTO{}
	err := useCase.Uow.Do(ctx, func(_ *uow.Uow) error {
		
		accountRepository := useCase.getAccountRepository(ctx)
		transactionRepository := useCase.getTransactionRepository(ctx)
		// Get account from
		accountFrom, err := accountRepository.Get(input.AccountIDFrom)
		if err != nil {
			return err
		}
		// Get account to
		accountTo, err := accountRepository.Get(input.AccountIDTo)
		if err != nil {
			return err
		}
		// Create transaction
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}
		// Update balance account from
		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		// Update balance account to
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}
		
		// Save transaction
		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}
		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount

		balanceUpdatedOutput.AccountIDFrom = input.AccountIDFrom
		balanceUpdatedOutput.AccountIDTo = input.AccountIDTo
		balanceUpdatedOutput.BalanceAccountIDFrom = accountFrom.Balance
		balanceUpdatedOutput.BalanceAccountIDTo = accountTo.Balance
		return nil
	})
	if err != nil {
		return nil, err
	}
	useCase.TransactionCreated.SetPayload(output)
	useCase.EventDispatcher.Dispatch(useCase.TransactionCreated)

	useCase.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	useCase.EventDispatcher.Dispatch(useCase.BalanceUpdated)
	return output, nil
}

func (usecase *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := usecase.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (usecase *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := usecase.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}