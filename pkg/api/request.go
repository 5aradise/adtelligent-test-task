package api

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON[T any](r *http.Request, dst *T) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
