package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BCoaracy/pacientezero/internal/model"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound
		msg = "not found"
	case errors.Is(err, model.ErrUnauthorized):
		status = http.StatusUnauthorized
		msg = "unauthorized"
	case errors.Is(err, model.ErrConflict):
		status = http.StatusConflict
		msg = "conflict"
	case errors.Is(err, model.ErrBadRequest):
		status = http.StatusBadRequest
		msg = "bad request"
	}

	writeJSON(w, status, map[string]string{"error": msg})
}
