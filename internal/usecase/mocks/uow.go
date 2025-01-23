package mocks

import (
	"context"

	"github.com/rodrigobunhak/fc-ms-wallet/pkg/uow"
	"github.com/stretchr/testify/mock"
)

type UowMock struct {
	mock.Mock
}

// todo: Register(name string, fc RepositoryFactory)
// todo: GetRepository(ctx context.Context, name string) (interface{}, error)
// todo: Do(ctx context.Context, fn func(Uow *Uow) error) error
// todo: CommitOrRollback() error
// todo: Rollback() error
// todo: UnRegister(name string)

func (m *UowMock) Register(name string, fc uow.RepositoryFactory) {
	m.Called(name, fc)
}

func (m *UowMock) GetRepository(ctx context.Context, name string) (interface{}, error) {
	args := m.Called(ctx, name)
	return args.Get(0), args.Error(1)
}

func (m *UowMock) Do(ctx context.Context, fn func(Uow *uow.Uow) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *UowMock) CommitOrRollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) UnRegister(name string) {
	m.Called(name)
}