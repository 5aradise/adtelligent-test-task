package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteText(w http.ResponseWriter, statusCode int, v string) error {
	const op = "api.WtireText"

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(v))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteHTML(w http.ResponseWriter, statusCode int, v string) error {
	const op = "api.WtireHTML"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(v))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	const op = "api.WtireJSON"

	data, marshalErr := json.Marshal(v)
	if marshalErr != nil {
		writeErr := WriteError(w, http.StatusInternalServerError, marshalErr.Error())
		if writeErr != nil {
			return fmt.Errorf("%s: %w, %w", op, marshalErr, writeErr)
		}
		return fmt.Errorf("%s: %w", op, marshalErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteM3U8(w http.ResponseWriter, statusCode int, data []byte) error {
	const op = "api.WriteM3U8"

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.WriteHeader(statusCode)
	_, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func WriteError(w http.ResponseWriter, statusCode int, msg string) error {
	const op = "api.WtireError"

	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_, err := fmt.Fprintln(w, msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
