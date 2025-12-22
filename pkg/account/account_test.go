package account

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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

// Create Account
// func TestCreateAccount(t *testing.T) {
// 	cfg := getTiggConfig(t)
// 	cl := client.New(cfg)
// 	svc := NewService(cl)

// 	timestamp := time.Now().UnixNano()
// 	parentGroupID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a"
// 	acc := Account{
// 		Code:          fmt.Sprintf("TEST-%d", timestamp),
// 		Name:          fmt.Sprintf("Test Account %d", timestamp),
// 		Description:   "This is test account",
// 		ParentGroupID: &parentGroupID,
// 	}

// 	created, err := svc.CreateAccount(acc)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "namespace is not registered") {
// 			t.Skip("Tigg namespace is not registered in this environment; skipping integration test")
// 		}
// 		t.Fatalf("CreateAccount failed: %v", err)
// 	}

// 	if created.ID == "" {
// 		t.Error("expected non-empty ID for created account")
// 	}
// }

// Test account retrieval
// func TestAccountRetrieval(t *testing.T) {
// 	cfg := getTiggConfig(t)
// 	cl := client.New(cfg)
// 	svc := NewService(cl)

// 	timestamp := time.Now().UnixNano()
// 	parentGroupID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a"
// 	accCode := fmt.Sprintf("RET-%d", timestamp)

// 	acc := Account{
// 		Code:          accCode,
// 		Name:          fmt.Sprintf("Retrieval Test Account %d", timestamp),
// 		Description:   "Test account for retrieval methods",
// 		ParentGroupID: &parentGroupID,
// 	}

// 	created, err := svc.CreateAccount(acc)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "namespace is not registered") {
// 			t.Skip("Tigg namespace is not registered; skipping")
// 		}
// 		t.Fatalf("Setup: CreateAccount failed: %v", err)
// 	}
// 	t.Logf("Created Account: ID=%s Code=%s Name=%s", created.ID, created.Code, created.Name)

// 	// 1. Test ListAccounts
// 	t.Run("ListAccounts", func(t *testing.T) {
// 		accounts, err := svc.ListAccounts()
// 		if err != nil {
// 			t.Fatalf("ListAccounts failed: %v", err)
// 		}
// 		found := false
// 		for _, a := range accounts {
// 			if a.ID == created.ID {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			t.Errorf("ListAccounts did not return the created account ID %s", created.ID)
// 		}
// 	})

// 	// 2. Test GetAccountByID
// 	t.Run("GetAccountByID", func(t *testing.T) {
// 		fetched, err := svc.GetAccountByID(created.ID)
// 		if err != nil {
// 			t.Fatalf("GetAccountByID failed: %v", err)
// 		}
// 		if fetched.ID != created.ID {
// 			t.Errorf("GetAccountByID returned ID %s, expected %s", fetched.ID, created.ID)
// 		}
// 		if fetched.Code != created.Code {
// 			t.Errorf("GetAccountByID returned code %s, expected %s", fetched.Code, created.Code)
// 		}
// 	})

// 	// // 3. Test GetAccountByCode
// 	t.Run("GetAccountByCode", func(t *testing.T) {
// 		fetched, err := svc.GetAccountByCode(created.Code)
// 		if err != nil {
// 			t.Fatalf("GetAccountByCode failed: %v", err)
// 		}
// 		if fetched.ID != created.ID {
// 			t.Errorf("GetAccountByCode returned ID %s, expected %s", fetched.ID, created.ID)
// 		}
// 	})
// }

// //Manual check

// func TestManualRetrieval(t *testing.T) {
// 	cfg := getTiggConfig(t)
// 	cl := client.New(cfg)
// 	svc := NewService(cl)

// 	targetID := "4429e03b-9dce-4dd9-9f25-d0a07dff520b"
// 	targetCode := "DE0008"

// 	t.Run("GetByID", func(t *testing.T) {
// 		acc, err := svc.GetAccountByID(targetID)
// 		if err != nil {
// 			t.Fatalf("Failed to get account by ID %s: %v", targetID, err)
// 		}
// 		t.Logf("Fetched by ID success: Code=%s Name=%s", acc.Code, acc.Name)
// 	})

// 	t.Run("GetByCode", func(t *testing.T) {
// 		acc, err := svc.GetAccountByCode(targetCode)
// 		if err != nil {
// 			t.Fatalf("Failed to get account by code %s: %v", targetCode, err)
// 		}
// 		t.Logf("Fetched by Code success: ID=%s Name=%s", acc.ID, acc.Name)
// 	})
// }

// // Test activate and deactivate account
// func TestActivateDeactivateAccount(t *testing.T) {
// 	cfg := getTiggConfig(t)
// 	cl := client.New(cfg)
// 	svc := NewService(cl)

// 	timestamp := time.Now().UnixNano()
// 	parentGroupID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a"
// 	acc := Account{
// 		Code:          fmt.Sprintf("ACT-%d", timestamp),
// 		Name:          fmt.Sprintf("Activation Test %d", timestamp),
// 		ParentGroupID: &parentGroupID,
// 	}

// 	created, err := svc.CreateAccount(acc)
// 	if err != nil {
// 		t.Fatalf("Failed to create account: %v", err)
// 	}

// 	// 1. Deactivate
// 	t.Run("Deactivate", func(t *testing.T) {
// 		updated, err := svc.DeactivateAccount(created.ID)
// 		if err != nil {
// 			t.Fatalf("DeactivateAccount failed: %v", err)
// 		}
// 		if !updated.Inactive {
// 			t.Errorf("Expected account to be inactive")
// 		}
// 	})

// 	// 2. Activate
// 	t.Run("Activate", func(t *testing.T) {
// 		updated, err := svc.ActivateAccount(created.ID)
// 		if err != nil {
// 			t.Fatalf("ActivateAccount failed: %v", err)
// 		}
// 		if updated.Inactive {
// 			t.Errorf("Expected account to be active")
// 		}
// 	})
// }

// Manual testing for Activate account
// func TestActivateAccount(t *testing.T) {
// 	cfg := getTiggConfig(t)
// 	cl := client.New(cfg)
// 	svc := NewService(cl)

// 	targetID := "1ccc5e66-0e0e-4bfd-ba30-867d5ac29a3c"

// 	updated, err := svc.ActivateAccount(targetID)
// 	if err != nil {
// 		t.Fatalf("ActivateAccount failed: %v", err)
// 	}

// 	if updated.Inactive {
// 		t.Fatalf("expected account to be active")
// 	}
// }

// Manual testing for Deactivate account
func TestDeactivateAccount(t *testing.T) {
	cfg := getTiggConfig(t)
	cl := client.New(cfg)
	svc := NewService(cl)

	targetID := "1ccc5e66-0e0e-4bfd-ba30-867d5ac29a3c"

	updated, err := svc.DeactivateAccount(targetID)
	if err != nil {
		t.Fatalf("DeactivateAccount failed: %v", err)
	}

	if !updated.Inactive {
		t.Fatalf("expected account to be inactive")
	}
}

func TestUpdateAccount(t *testing.T) {
	cfg := getTiggConfig(t)
	cl := client.New(cfg)
	svc := NewService(cl)

	// 1. Create a fresh account to update
	timestamp := time.Now().UnixNano()
	parentGroupID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a"
	createReq := Account{
		Code:          fmt.Sprintf("UPD-%d", timestamp),
		Name:          fmt.Sprintf("Update Test %d", timestamp),
		ParentGroupID: &parentGroupID,
	}

	created, err := svc.CreateAccount(createReq)
	if err != nil {
		if strings.Contains(err.Error(), "namespace is not registered") {
			t.Skip("Tigg namespace is not registered; skipping")
		}
		t.Fatalf("Failed to create account for update test: %v", err)
	}
	t.Logf("Created Account for update: %+v", created)

	// 2. Perform update
	updateReq := UpdateAccountRequest{
		Name:          fmt.Sprintf("Updated Name %d", timestamp),
		Code:          created.Code, // Keep same code assigned by API
		Description:   "Updated description",
		ParentGroupID: &parentGroupID,
	}

	updated, err := svc.UpdateAccount(created.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdateAccount failed: %v", err)
	}

	// 3. Assertions
	if updated.ID != created.ID {
		t.Errorf("Update created a new record! Expected ID %s, got %s", created.ID, updated.ID)
	}
	if updated.Name != updateReq.Name {
		t.Errorf("Expected name %s, got %s", updateReq.Name, updated.Name)
	}
	if updated.Description != updateReq.Description {
		t.Errorf("Expected description %s, got %s", updateReq.Description, updated.Description)
	}
}

func TestManualUpdateAccount(t *testing.T) {
	cfg := getTiggConfig(t)
	cl := client.New(cfg)
	svc := NewService(cl)

	targetID := "1ccc5e66-0e0e-4bfd-ba30-867d5ac29a3c"
	timestamp := time.Now().UnixNano()
	parentGroupID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a"

	updateReq := UpdateAccountRequest{
		Name:          fmt.Sprintf("Manual Update %d", timestamp),
		Code:          "DE0060", // Original code for this ID
		Description:   "Manually updated via test",
		ParentGroupID: &parentGroupID,
	}

	updated, err := svc.UpdateAccount(targetID, updateReq)
	if err != nil {
		t.Fatalf("UpdateAccount failed for ID %s: %v", targetID, err)
	}

	if updated.Name != updateReq.Name {
		t.Errorf("Expected name %s, got %s", updateReq.Name, updated.Name)
	}
	t.Logf("Successfully updated account %s to name %s", targetID, updated.Name)
}
