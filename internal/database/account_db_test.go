package database

import (
	"database/sql"
	"testing"

	"github.com/rodrigobunhak/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db *sql.DB
	accountDB *AccountDB
	client *entity.Client
}

func (suite *AccountDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.Nil(err)
	suite.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance int, created_at date)")
	suite.accountDB = NewAccountDB(db)
	suite.client, _ = entity.NewClient("John Doe", "j@j.com")
}

func (suite *AccountDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE accounts")
	suite.db.Exec("DROP TABLE clients")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (suite *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(suite.client)
	err := suite.accountDB.Save(account)
	suite.Nil(err)
}

func (suite *AccountDBTestSuite) TestGet() {
	suite.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)", suite.client.ID, suite.client.Name, suite.client.Email, suite.client.CreatedAt)
	account := entity.NewAccount(suite.client)
	err := suite.accountDB.Save(account)
	suite.Nil(err)
	accountDB, err := suite.accountDB.Get(account.ID)
	suite.Nil(err)
	suite.Equal(account.ID, accountDB.ID)
	suite.Equal(account.Client.ID, accountDB.Client.ID)
	suite.Equal(account.Balance, accountDB.Balance)
	suite.Equal(account.Client.ID, accountDB.Client.ID)
	suite.Equal(account.Client.Name, accountDB.Client.Name)
	suite.Equal(account.Client.Email, accountDB.Client.Email)
}