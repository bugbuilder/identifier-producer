package server

import (
	"bennu.cl/identifier-producer/api/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func ToJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func FromJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		if err == io.EOF {
			return errors.New("got EOF while reading request body")
		}
		return err
	}
	return nil
}

func getHTTPErrorStatusCode(err error) int {
	return http.StatusInternalServerError
}

func MakeErrorHandler(statusCode int, err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := &types.ErrorResponse{
			Message: err.Error(),
		}
		_ = ToJSON(w, statusCode, response)
	}
}
