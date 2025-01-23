package database

import (
	"database/sql"
	"testing"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db *sql.DB
	client *entity.Client
	client2 *entity.Client
	accountFrom *entity.Account
	accountTo *entity.Account
	transactionDB *TransactionDB
}

func (suite *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.Nil(err)
	suite.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	db.Exec("CREATE TABLE transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")
	client, err := entity.NewClient("John Doe", "j@j.com")
	suite.Nil(err)
	suite.client = client
	client2, err := entity.NewClient("Jane Doe2", "j@j2.com")
	suite.Nil(err)
	suite.client2 = client2
	// creating accounts
	accountFrom := entity.NewAccount(suite.client)
	accountFrom.Balance = 1000
	suite.accountFrom = accountFrom
	accountTo := entity.NewAccount(suite.client2)
	accountTo.Balance = 1000
	suite.accountTo = accountTo	
	suite.transactionDB = NewTransactionDB(db)
}

func (suite *TransactionDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE transactions")
	suite.db.Exec("DROP TABLE accounts")
	suite.db.Exec("DROP TABLE clients")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (suite *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(suite.accountFrom, suite.accountTo, 100)
	suite.Nil(err)
	err = suite.transactionDB.Create(transaction)
	suite.Nil(err)
}