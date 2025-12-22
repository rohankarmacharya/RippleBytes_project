package accountgroup

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

// Service wraps the Tigg client
type Service struct {
	client *client.TiggClient
}

// NewService constructor makes it reusable
func NewService(c *client.TiggClient) *Service {
	return &Service{client: c}
}

// CreateAccountGroup
type createAccountGroupResponse struct {
	Data AccountGroup `json:"data"`
}

func (s *Service) CreateAccountGroup(reqBody CreateAccountGroupRequest) (*AccountGroup, error) {
	url := fmt.Sprintf("%s/account-groups", s.client.BaseURL)

	// Generate timestamp (ms) and nonce as required by Tigg signature spec
	timestampMs := time.Now().UnixMilli()
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())

	// Build unsigned payload with timestamp and nonce
	type unsignedPayload struct {
		CreateAccountGroupRequest
		Timestamp int64  `json:"timestamp"`
		Nonce     string `json:"nonce"`
	}

	up := unsignedPayload{
		CreateAccountGroupRequest: reqBody,
		Timestamp:                 timestampMs,
		Nonce:                     nonce,
	}

	unsignedJSON, err := json.Marshal(up)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	// Base64-encode the JSON string
	payloadString := base64.StdEncoding.EncodeToString(unsignedJSON)

	// Sign the base64-encoded payload with HMAC-SHA256 using the client secret
	mac := hmac.New(sha256.New, []byte(s.client.SecretKey))
	mac.Write([]byte(payloadString))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Final payload includes the signature alongside timestamp and nonce
	type signedPayload struct {
		CreateAccountGroupRequest
		Timestamp int64  `json:"timestamp"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}

	sp := signedPayload{
		CreateAccountGroupRequest: reqBody,
		Timestamp:                 timestampMs,
		Nonce:                     nonce,
		Signature:                 signature,
	}

	payloadBytes, err := json.Marshal(sp)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestampMs))

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var res createAccountGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// GetAccountGroupByID
type getAccountGroupResponse struct {
	Data AccountGroup `json:"data"`
}

func (s *Service) GetAccountGroupByID(id string) (*AccountGroup, error) {
	url := fmt.Sprintf("%s/account-groups/%s", s.client.BaseURL, id)

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

	var res getAccountGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// ListAccountGroups
type accountGroupListResponse struct {
	Data []AccountGroup `json:"data"`
}

func (s *Service) ListAccountGroups() ([]AccountGroup, error) {
	url := fmt.Sprintf("%s/account-groups", s.client.BaseURL)

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

	var res accountGroupListResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// GetAccountGroupByName performs a ListAccountGroups call and searches by exact name.
func (s *Service) GetAccountGroupByName(name string) (*AccountGroup, error) {
	groups, err := s.ListAccountGroups()
	if err != nil {
		return nil, err
	}

	for i := range groups {
		if groups[i].Name == name {
			return &groups[i], nil
		}
	}

	return nil, fmt.Errorf("account group with name %q not found", name)
}

// UpdateAccountGroup
func (s *Service) UpdateAccountGroup(id string, reqBody UpdateAccountGroupRequest) (*AccountGroup, error) {
	// SAFEGUARD: The 'id' is strictly required. Providing an ID identifies this as an UPDATE operation.
	// If the ID is missing or empty in the payload, Tigg will treat this as a CREATE operation and generate a new record!
	if id == "" {
		return nil, fmt.Errorf("id is required for update to prevent duplicate creation")
	}

	// FORCE the ID from the function argument into the payload to guarantee it's present.
	// This ensures Tigg sees the 'id' field in the JSON body.
	reqBody.ID = id

	url := fmt.Sprintf("%s/account-groups/%s", s.client.BaseURL, id)

	// Generate timestamp (ms) and nonce as required by Tigg signature spec
	timestampMs := time.Now().UnixMilli()
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())

	// Build unsigned payload with timestamp and nonce
	type unsignedUpdatePayload struct {
		UpdateAccountGroupRequest
		Timestamp int64  `json:"timestamp"`
		Nonce     string `json:"nonce"`
	}

	up := unsignedUpdatePayload{
		UpdateAccountGroupRequest: reqBody,
		Timestamp:                 timestampMs,
		Nonce:                     nonce,
	}

	unsignedJSON, err := json.Marshal(up)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	// Base64-encode the JSON string
	payloadString := base64.StdEncoding.EncodeToString(unsignedJSON)

	// Sign the base64-encoded payload with HMAC-SHA256 using the client secret
	mac := hmac.New(sha256.New, []byte(s.client.SecretKey))
	mac.Write([]byte(payloadString))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Final payload includes the signature alongside timestamp and nonce
	type signedUpdatePayload struct {
		UpdateAccountGroupRequest
		Timestamp int64  `json:"timestamp"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}

	sp := signedUpdatePayload{
		UpdateAccountGroupRequest: reqBody,
		Timestamp:                 timestampMs,
		Nonce:                     nonce,
		Signature:                 signature,
	}

	payloadBytes, err := json.Marshal(sp)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	s.client.AddHeaders(req)
	req.Header.Set("X-Nonce", nonce)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestampMs))

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.NewTiggError(resp)
	}

	var res createAccountGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// Activate/Deactivate AccountGroup
func (s *Service) ActivateAccountGroup(id string) error {
	url := fmt.Sprintf("%s/account-groups/%s/active", s.client.BaseURL, id)
	req, _ := http.NewRequest("PATCH", url, nil)
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

func (s *Service) DeactivateAccountGroup(id string) error {
	url := fmt.Sprintf("%s/account-groups/%s/inactive", s.client.BaseURL, id)
	req, _ := http.NewRequest("PATCH", url, nil)
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
