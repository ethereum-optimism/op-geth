package builder

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type httpErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handleError(w http.ResponseWriter, err error) {
	var errorMsg string
	var status int
	switch {
	case errors.Is(err, ErrIncorrectSlot):
		errorMsg = err.Error()
		status = http.StatusBadRequest
	case errors.Is(err, ErrNoPayloads):
		errorMsg = err.Error()
		status = http.StatusNotFound
	case errors.Is(err, ErrSlotFromPayload):
		errorMsg = err.Error()
		status = http.StatusInternalServerError
	case errors.Is(err, ErrSlotMismatch):
		errorMsg = err.Error()
		status = http.StatusBadRequest
	case errors.Is(err, ErrParentHashFromPayload):
		errorMsg = err.Error()
		status = http.StatusInternalServerError
	case errors.Is(err, ErrParentHashMismatch):
		errorMsg = err.Error()
		status = http.StatusBadRequest
	default:
		errorMsg = "error processing request"
		status = http.StatusInternalServerError
	}
	
	respondError(w, status, errorMsg)
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(httpErrorResp{code, message}); err != nil {
		http.Error(w, message, code)
	}
}

// runRetryLoop calls retry periodically with the provided interval respecting context cancellation
func runRetryLoop(ctx context.Context, interval time.Duration, retry func()) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			retry()
		}
	}
}
