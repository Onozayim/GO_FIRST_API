package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ValidateBody(body any, w http.ResponseWriter, r *http.Request) error {
	if r.Body == http.NoBody {
		return fmt.Errorf("body is empty")
	}

	return json.NewDecoder(r.Body).Decode(&body)
}
