package client

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("../../.env")
}

func TestDebugAPI(t *testing.T) {
	cfg := Config{
		BaseURL:   os.Getenv("TIGG_API_URL"),
		ClientKey: os.Getenv("TIGG_CLIENT_KEY"),
		SecretKey: os.Getenv("TIGG_SECRET_KEY"),
		Namespace: os.Getenv("TIGG_NAMESPACE"),
	}

	if cfg.BaseURL == "" {
		t.Skip("Skipping debug test: env vars not set")
	}

	c := New(cfg)

	// Helper to run requests
	runCase := func(name string, setupReq func(*http.Request)) {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", cfg.BaseURL+"/account-groups", nil)
			c.AddHeaders(req)
			setupReq(req)
			doRequest(t, c, req)
		})
	}

	// 1. Standard (As Implemented)
	runCase("Standard", func(req *http.Request) {
		c.SignRequest(req, []byte(""))
	})

	// 2. Invalid Namespace
	runCase("Invalid Namespace", func(req *http.Request) {
		req.Header.Set("Namespace", "invalid-ns")
	})

	// 3. Invalid Client Key
	runCase("Invalid Client Key", func(req *http.Request) {
		req.Header.Set("X-Client-Key", "invalid-key")
	})

	// 4. Missing Client Key
	runCase("Missing Client Key", func(req *http.Request) {
		req.Header.Del("X-Client-Key")
	})

	// 5. Missing Signature
	runCase("Missing Signature", func(req *http.Request) {
		req.Header.Del("X-Signature")
	})

	// 6. Missing Nonce
	runCase("Missing Nonce", func(req *http.Request) {
		req.Header.Del("X-Nonce")
	})

	// 7. Missing Timestamp
	runCase("Missing Timestamp", func(req *http.Request) {
		req.Header.Del("X-Timestamp")
	})
}

func doRequest(t *testing.T, c *TiggClient, req *http.Request) {
	t.Logf("Request: %s %s", req.Method, req.URL)
	for k, v := range req.Header {
		t.Logf("Header: %s=%v", k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response Status: %d", resp.StatusCode)
	bodyStr := string(body)
	if len(bodyStr) > 500 {
		bodyStr = bodyStr[:500] + "..."
	}
	t.Logf("Response Body: %s", bodyStr)
}
