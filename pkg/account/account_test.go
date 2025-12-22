package account

import (
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rohankarmacharya/TigIntegration/pkg/client"
)

func init() {
	_ = godotenv.Load("../../.env")
}

func getTiggConfig(t *testing.T) client.Config {
	baseURL := os.Getenv("TIGG_API_URL")
	clientKey := os.Getenv("TIGG_CLIENT_KEY")
	secretKey := os.Getenv("TIGG_SECRET_KEY")
	namespace := os.Getenv("TIGG_NAMESPACE")

	if baseURL == "" || clientKey == "" || secretKey == "" || namespace == "" {
		t.Skip("Tigg integration not configured; set TIGG_API_URL, TIGG_CLIENT_KEY, TIGG_SECRET_KEY, TIGG_NAMESPACE to run this test")
	}

	return client.Config{
		ClientKey: clientKey,
		SecretKey: secretKey,
		Namespace: namespace,
		BaseURL:   baseURL,
	}
}

func TestCreateAccount(t *testing.T) {
	cfg := getTiggConfig(t)
	cl := client.New(cfg)
	svc := NewService(cl)

	acc := Account{
		Code:        "TEST-001",
		Name:        "Test Account",
		Description: "This is test account",
	}

	created, err := svc.CreateAccount(acc)
	if err != nil {
		if strings.Contains(err.Error(), "namespace is not registered") {
			t.Skip("Tigg namespace is not registered in this environment; skipping integration test")
		}
		t.Fatalf("CreateAccount failed: %v", err)
	}

	if created.ID == "" {
		t.Error("expected non-empty ID for created account")
	}
}
