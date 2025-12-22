package account

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rohankarmacharya/TigIntegration/pkg/client"
	"github.com/rohankarmacharya/TigIntegration/pkg/errors"
)

type Service struct {
	client *client.TiggClient
}

func NewService(c *client.TiggClient) *Service {
	return &Service{client: c}
}

type accountListResponse struct {
	Data []Account `json:"data"`
}

type accountResponse struct {
	Data Account `json:"data"`
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

	var res accountListResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// CreateAccount that sends POST /accounts request to Tigg

func (s *Service) signPayload(method, url string, payload interface{}) (*http.Request, error) {
	timestampMs := time.Now().UnixMilli()
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())

	finalPayload := make(map[string]interface{})

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &finalPayload); err != nil {
			return nil, err
		}
	}
	finalPayload["timestamp"] = timestampMs
	finalPayload["nonce"] = nonce

	unsignedJSON, err := json.Marshal(finalPayload)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	payloadString := base64.StdEncoding.EncodeToString(unsignedJSON)

	mac := hmac.New(sha256.New, []byte(s.client.SecretKey))
	mac.Write([]byte(payloadString))
	signature := hex.EncodeToString(mac.Sum(nil))

	finalPayload["signature"] = signature

	signedJSON, err := json.Marshal(finalPayload)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(signedJSON))
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestampMs))

	return req, nil
}
func (s *Service) CreateAccount(acc Account) (*Account, error) {
	url := fmt.Sprintf("%s/accounts", s.client.BaseURL)

	req, err := s.signPayload("POST", url, acc)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var res accountResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// UpdateAccount sends POST /accounts/{id} request to Tigg
func (s *Service) UpdateAccount(id string, acc Account) (*Account, error) {
	url := fmt.Sprintf("%s/accounts/%s", s.client.BaseURL, id)

	req, err := s.signPayload("POST", url, acc)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var res accountResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Data, nil
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

	var res accountResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// GetAccountByCode performs a ListAccounts call and searches by exact code.
func (s *Service) GetAccountByCode(code string) (*Account, error) {
	accounts, err := s.ListAccounts()
	if err != nil {
		return nil, err
	}

	for i := range accounts {
		if accounts[i].Code == code {
			return &accounts[i], nil
		}
	}

	return nil, fmt.Errorf("account with code %q not found", code)
}

// ActivateAccount sends PATCH /accounts/{id}/active request to Tigg
func (s *Service) ActivateAccount(id string) error {
	url := fmt.Sprintf("%s/accounts/%s/active", s.client.BaseURL, id)

	req, err := s.signPayload("PATCH", url, map[string]string{})
	if err != nil {
		return err
	}

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

	req, err := s.signPayload("PATCH", url, map[string]string{})
	if err != nil {
		return err
	}

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
