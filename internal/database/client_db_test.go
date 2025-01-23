package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db *sql.DB
	ClientDB *ClientDB
}

func (suite *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.Nil(err)
	suite.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at datetime)")
	suite.ClientDB = NewClientDB(db)
}

func (suite *ClientDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (suite *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID: "1",
		Name: "John Doe",
		Email: "j@j.com",
	}
	err := suite.ClientDB.Save(client)
	suite.Nil(err)
}

func (suite *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "j@j.com")
	suite.ClientDB.Save(client)
	ClientDB, err := suite.ClientDB.Get(client.ID)
	suite.Nil(err)
	suite.Equal(client.ID, ClientDB.ID)
	suite.Equal(client.Name, ClientDB.Name)
	suite.Equal(client.Email, ClientDB.Email)
}