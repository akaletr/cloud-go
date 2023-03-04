package transactionLogger

import (
	"bufio"
	"fmt"
	"os"
)

type EventType uint64

const (
	_                  = iota
	EventPut EventType = iota
	EventDelete
)

type Event struct {
	ID     uint64
	Method EventType
	Key    string
	Value  string
}

type fileTransactionLogger struct {
	events  chan<- Event
	errors  <-chan error
	eventID uint64
	file    *os.File
}

func NewFile(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log fileTransactionLogger: %w", err)
	}

	return &fileTransactionLogger{file: file}, nil
}

func (logger *fileTransactionLogger) Run() {
	events := make(chan Event, 16)
	logger.events = events
	errors := make(chan error, 1)
	logger.errors = errors

	go func() {
		for e := range events {
			logger.eventID++
			_, err := fmt.Fprintf(logger.file, "%d\t%d\t%s\t%s\n", logger.eventID, e.Method, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (logger *fileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(logger.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()
			_, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.ID, &e.Method, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}
			if logger.eventID >= e.ID {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}
			logger.eventID = e.ID
			outEvent <- e
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvent, outError
}

func (logger *fileTransactionLogger) WritePut(key, value string) {
	logger.events <- Event{
		Method: EventPut,
		Key:    key,
		Value:  value,
	}
}

func (logger *fileTransactionLogger) WriteDelete(key string) {
	logger.events <- Event{
		Method: EventDelete,
		Key:    key,
	}
}

func (logger *fileTransactionLogger) Err() <-chan error {
	return logger.errors
}
