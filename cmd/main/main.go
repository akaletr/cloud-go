package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"cmd/main/main.go/internal/transactionLogger"

	"github.com/go-chi/chi/v5"
)

var ErrorNoSuchKey = errors.New("error: no such key")
var logger transactionLogger.TransactionLogger

var store = struct {
	sync.RWMutex
	storage map[string]string
}{storage: make(map[string]string)}

func Put(key, value string) error {
	store.Lock()
	defer store.Unlock()

	store.storage[key] = value

	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	defer store.RUnlock()
	value, ok := store.storage[key]
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.storage, key)
	return nil
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = Put(key, string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WritePut(key, string(body))

	w.WriteHeader(http.StatusCreated)
	return
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value, err := Get(key)
	if err != nil {
		if errors.Is(err, ErrorNoSuchKey) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	err := Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.WriteDelete(key)

	w.WriteHeader(http.StatusOK)
	return
}

func initializeLogger() error {
	var err error
	logger, err = transactionLogger.NewFile("test.txt")
	if err != nil {
		fmt.Println(err)
	}

	events, errs := logger.ReadEvents()
	e, ok := transactionLogger.Event{}, true
	for ok && err == nil {
		select {
		case e, ok = <-events:
			switch e.Method {
			case transactionLogger.EventPut:
				err = Put(e.Key, e.Value)
				if err != nil {
					return err
				}
			case transactionLogger.EventDelete:
				err = Delete(e.Key)
				if err != nil {
					return err
				}
			}
		case err, ok = <-errs:
			fmt.Println(err)
		}
	}

	logger.Run()
	return err
}

func main() {
	router := chi.NewRouter()

	go func() {
		err := initializeLogger()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	router.Put("/v1/{key}", putHandler)
	router.Get("/v1/{key}", getHandler)
	router.Delete("/v1/{key}", deleteHandler)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Fatalln(server.ListenAndServe())
}
