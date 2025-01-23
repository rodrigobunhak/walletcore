package create_transaction

import (
	"context"
	"testing"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/entity"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/event"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/mocks"
	"github.com/rodrigobunhak/fc-ms-wallet/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (mock *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := mock.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (mock *AccountGatewayMock) Save(account *entity.Account) error {
	args := mock.Called(account)
	return args.Error(0)
}

func (mock *AccountGatewayMock) Get(id string) (*entity.Account, error) {
	args := mock.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "j@j.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)
	client2, _ := entity.NewClient("John Doe", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)	

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo: account2.ID,
		Amount: 100,
	}
	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	ctx := context.Background()

	useCase := NewCreateTransactionUseCase(mockUow, dispatcher, event)
	output, err := useCase.Execute(ctx, inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output)

	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}