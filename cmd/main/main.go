package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

var ErrorNoSuchKey = errors.New("error: no such key")

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

func cloudPutHandlerChi(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusCreated)
	return
}

func cloudGetHandlerChi(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	fmt.Println(store)
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

func cloudDeleteHandlerChi(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	err := Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	router := chi.NewRouter()

	router.Put("/v1/{key}", cloudPutHandlerChi)
	router.Get("/v1/{key}", cloudGetHandlerChi)
	router.Delete("/v1/{key}", cloudDeleteHandlerChi)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	log.Fatalln(server.ListenAndServe())
}
