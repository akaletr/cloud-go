package app

import (
	"cmd/main/main.go/internal/mw"
	"fmt"
	"net/http"

	"cmd/main/main.go/internal/logger"
	"cmd/main/main.go/internal/storage"
	"cmd/main/main.go/internal/transactionLogger"

	"github.com/go-chi/chi/v5"
)

type app struct {
	storage     storage.Storage
	transaction transactionLogger.TransactionLogger
	server      *http.Server
	logger      logger.Logger
}

func New() (App, error) {
	st := storage.New()
	tr, err := transactionLogger.NewFile("test.log")
	if err != nil {
		return nil, err
	}

	srv := http.Server{
		Addr: ":8000",
	}

	lg := logger.New()

	return &app{
		storage:     st,
		transaction: tr,
		server:      &srv,
		logger:      lg,
	}, nil
}

func (app *app) Init() error {
	router := chi.NewRouter()

	router.Use(mw.Limiter)

	router.Put("/v1/{key}", app.putHandler)
	router.Get("/v1/{key}", app.getHandler)
	router.Delete("/v1/{key}", app.deleteHandler)

	app.server.Handler = router

	go app.initTransactionLogger()
	app.transaction.Run()

	return nil
}

func (app *app) initTransactionLogger() {
	var err error

	events, errs := app.transaction.ReadEvents()
	e, ok := transactionLogger.Event{}, true

	for ok && err == nil {
		select {
		case e, ok = <-events:
			switch e.Method {
			case transactionLogger.EventPut:
				err = app.storage.Put(e.Key, e.Value)
				if err != nil {
					app.logger.Debug(err)
				}
			case transactionLogger.EventDelete:
				err = app.storage.Delete(e.Key)
				if err != nil {
					app.logger.Debug(err)
				}
			}
		case err, ok = <-errs:
			fmt.Println(err)
		}
	}
}

func (app *app) Start() error {
	err := app.Init()
	if err != nil {
		return err
	}

	return app.server.ListenAndServe()
}

func (app *app) Stop() error {
	return app.transaction.Close()
}
