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

// CreateAccountGroup
type createAccountGroupResponse struct {
	Data AccountGroup `json:"data"`
}

// signPayload handles the Tigg API signature generation and request creation
func (s *Service) signPayload(method, url string, payload interface{}) (*http.Request, error) {
	// Generate timestamp (ms) and nonce
	timestampMs := time.Now().UnixMilli()
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())

	// Create a dynamic struct to hold the payload + signature fields
	// We use a map[string]interface{} to merge the original payload with extra fields
	// This avoids defining specific structs for every request type
	finalPayload := make(map[string]interface{})

	// helper to marshal->unmarshal to map
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &finalPayload); err != nil {
			return nil, err
		}
	}

	// 1. Create Unsigned Payload (Payload + Timestamp + Nonce)
	finalPayload["timestamp"] = timestampMs
	finalPayload["nonce"] = nonce

	unsignedJSON, err := json.Marshal(finalPayload)
	if err != nil {
		return nil, errors.ErrInvalidPayLoad
	}

	// 2. Base64 Encode: converts binary data into an ASCII string format; converts the JSON payload into a string
	payloadString := base64.StdEncoding.EncodeToString(unsignedJSON)

	// 3. Sign: creates a hash of the payload string using HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(s.client.SecretKey))
	mac.Write([]byte(payloadString))
	signature := hex.EncodeToString(mac.Sum(nil))

	// 4. Add Signature to Payload
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

func (s *Service) CreateAccountGroup(reqBody CreateAccountGroupRequest) (*AccountGroup, error) {
	url := fmt.Sprintf("%s/account-groups", s.client.BaseURL)
	req, err := s.signPayload("POST", url, reqBody)
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
	if id == "" {
		return nil, fmt.Errorf("id is required for update to prevent duplicate creation")
	}
	reqBody.ID = id

	url := fmt.Sprintf("%s/account-groups/%s", s.client.BaseURL, id)
	req, err := s.signPayload("POST", url, reqBody) // Update uses POST as per original code
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

	var res createAccountGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

// Activate/Deactivate AccountGroup
// ActivateAccountGroup
func (s *Service) ActivateAccountGroup(id string) (*AccountGroup, error) {
	url := fmt.Sprintf("%s/account-groups/%s/active", s.client.BaseURL, id)
	// Even if there's no body content, we might need the signature wrapper.
	// Passing an empty struct or nil. Let's try passing empty struct to ensure signature fields are added.
	req, err := s.signPayload("PATCH", url, map[string]string{})
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

	// API returns success status but not the object, so we fetch it
	return s.GetAccountGroupByID(id)
}

func (s *Service) DeactivateAccountGroup(id string) (*AccountGroup, error) {
	url := fmt.Sprintf("%s/account-groups/%s/inactive", s.client.BaseURL, id)
	req, err := s.signPayload("PATCH", url, map[string]string{})
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

	// API returns success status but not the object, so we fetch it
	return s.GetAccountGroupByID(id)
}
