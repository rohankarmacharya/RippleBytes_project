package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type TiggClient struct {
	ClientKey  string
	SecretKey  string
	Namespace  string
	BaseURL    string
	HTTPClient *http.Client
}

// Config struct to initialize client
type Config struct {
	ClientKey string
	SecretKey string
	Namespace string
	BaseURL   string
}

// New creates a new Tigg client using Config
func New(cfg Config) *TiggClient {
	return &TiggClient{
		ClientKey:  cfg.ClientKey,
		SecretKey:  cfg.SecretKey,
		Namespace:  cfg.Namespace,
		BaseURL:    cfg.BaseURL,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// AddHeaders adds required Tigg headers to any request
func (c *TiggClient) AddHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.ClientKey) // was X-Client-Key
	req.Header.Set("X-Nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	req.Header.Set("namespace", c.Namespace) // was Namespace
}

// SignRequest adds HMAC-SHA256 signature to request
func (c *TiggClient) SignRequest(req *http.Request, body []byte) {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Signature", signature)
}
