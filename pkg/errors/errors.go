package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrInvalidPayLoad = fmt.Errorf("invalid payload")
)

type TiggError struct {
	StatusCode int
	Message    string
}

func (e *TiggError) Error() string {
	return fmt.Sprintf("tigg api error: status=%d message=%s", e.StatusCode, e.Message)
}

func NewTiggError(resp *http.Response) error {
	defer resp.Body.Close()
	var body struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return &TiggError{
		StatusCode: resp.StatusCode,
		Message:    body.Message,
	}
}
