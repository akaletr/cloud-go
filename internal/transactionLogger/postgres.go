package transactionLogger

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresTransactionLogger struct {
	events  chan<- Event
	errors  <-chan error
	eventID uint64
	conn    *gorm.DB
}

func NewPostgres() (TransactionLogger, error) {
	dsn := "host=localhost user=postgres password=adelaida2011 dbname=cloud port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &postgresTransactionLogger{conn: db}, nil
}

func (logger *postgresTransactionLogger) WritePut(key, value string) {
	logger.events <- Event{
		Method: EventPut,
		Key:    key,
		Value:  value,
	}
}

func (logger *postgresTransactionLogger) WriteDelete(key string) {
	logger.events <- Event{
		Method: EventDelete,
		Key:    key,
	}
}

func (logger *postgresTransactionLogger) Err() <-chan error {
	return logger.errors
}

func (logger *postgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	logger.events = events
	errors := make(chan error, 1)
	logger.errors = errors

	go func() {
		for event := range events {
			logger.eventID++
			event.ID = logger.eventID
			logger.conn.Create(&event)
		}
	}()
}

func (logger *postgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		var events []Event
		logger.conn.Find(&events)

		for _, event := range events {
			logger.eventID = event.ID
			outEvent <- event
		}
	}()
	return outEvent, outError
}

func (logger *postgresTransactionLogger) Close() error {
	return nil
}
