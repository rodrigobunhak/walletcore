package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/database"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/event"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/create_account"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/create_client"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/web"
	"github.com/rodrigobunhak/fc-ms-wallet/internal/web/webserver"
	"github.com/rodrigobunhak/fc-ms-wallet/pkg/events"
	"github.com/rodrigobunhak/fc-ms-wallet/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()

	// eventDispatcher.Register("TransactionCreated", handler)
	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func (tx *sql.Tx) interface{} {
		return accountDb // TODO: Ver se posso passar apenas accountDb
	})

	uow.Register("TransactionDB", func (tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

	webServer := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webServer.Start()
}