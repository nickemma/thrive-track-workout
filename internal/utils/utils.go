package utils

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func ReadJSON(r *http.Request) (int64, error) {
	idParams := chi.URLParam(r, "id")
	if idParams == "" {
		return 0, errors.New("invalid id parameter")
	}
	id, err := strconv.ParseInt(idParams, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}
