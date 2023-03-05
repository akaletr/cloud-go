package app

import (
	"errors"
	"io"
	"net/http"

	"cmd/main/main.go/internal/storage"

	"github.com/go-chi/chi/v5"
)

func (app *app) putHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = app.storage.Put(key, string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	app.transaction.WritePut(key, string(body))

	w.WriteHeader(http.StatusCreated)
	return
}

func (app *app) getHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value, err := app.storage.Get(key)
	if err != nil {
		if errors.Is(err, storage.ErrorNoSuchKey) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (app *app) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	err := app.storage.Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	app.transaction.WriteDelete(key)

	w.WriteHeader(http.StatusOK)
	return
}
