package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rohankarmacharya/TigIntegration/pkg/client"
	"github.com/rohankarmacharya/TigIntegration/pkg/errors"
)

type Service struct {
	client *client.TiggClient
}

func NewService(c *client.TiggClient) *Service {
	return &Service{client: c}
}

// CreateAccount that sends POST /accounts request to Tigg
func (s *Service) CreateAccount(acc Account) (*Account, error) {
	url := fmt.Sprintf("%s/accounts", s.client.BaseURL)

	bodyBytes, err := json.Marshal(acc)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)
	s.client.SignRequest(req, bodyBytes)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var created Account
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateAccount sends POST /accounts/{id} request to Tigg
func (s *Service) UpdateAccount(id string, acc Account) (*Account, error) {
	url := fmt.Sprintf("%s/accounts/%s", s.client.BaseURL, id)

	bodyBytes, err := json.Marshal(acc)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)
	s.client.SignRequest(req, bodyBytes)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var updated Account
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// GetAccountByID sends GET /accounts/{id} request to Tigg
func (s *Service) GetAccountByID(id string) (*Account, error) {
	url := fmt.Sprintf("%s/accounts/%s", s.client.BaseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var acc Account
	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}
	return &acc, nil
}

// ListAccounts sends GET /accounts request to Tigg
func (s *Service) ListAccounts() ([]Account, error) {
	url := fmt.Sprintf("%s/accounts", s.client.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var accounts []Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// ActivateAccount sends PATCH /accounts/{id}/active request to Tigg
func (s *Service) ActivateAccount(id string) error {
	url := fmt.Sprintf("%s/accounts/%s/active", s.client.BaseURL, id)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return err
	}

	s.client.AddHeaders(req)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.NewTiggError(resp)
	}

	return nil
}

// DeactivateAccount sends PATCH /accounts/{id}/inactive request to Tigg
func (s *Service) DeactivateAccount(id string) error {
	url := fmt.Sprintf("%s/accounts/%s/inactive", s.client.BaseURL, id)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return err
	}

	s.client.AddHeaders(req)

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.NewTiggError(resp)
	}

	return nil
}
