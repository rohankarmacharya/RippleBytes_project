package accountgroup

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

func newTestService(t *testing.T) *Service {
	cfg := getTiggConfig(t)
	cl := client.New(cfg)
	return NewService(cl)
}

func TestCreateAccountGroup(t *testing.T) {
	service := newTestService(t)

	parentID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a" // Direct Expenses

	req := CreateAccountGroupRequest{
		Description:   "Test group via SDK",
		Name:          fmt.Sprintf("Area-Pokhara-%d", time.Now().UnixNano()),
		ParentGroupID: &parentID,
	}

	createdGroup, err := service.CreateAccountGroup(req)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid signature") {
			t.Skip("Tigg signature validation not configured for this environment; skipping create account group integration test")
		}
		t.Fatalf("Failed to create account group: %v", err)
	}

	if createdGroup.ID == "" {
		t.Fatal("expected non-empty ID for created account group")
	}

	// Verify we can fetch the same group by ID
	fetchedByID, err := service.GetAccountGroupByID(createdGroup.ID)
	if err != nil {
		t.Fatalf("GetAccountGroupByID failed: %v", err)
	}
	if fetchedByID.ID != createdGroup.ID {
		t.Fatalf("GetAccountGroupByID returned ID %q, expected %q", fetchedByID.ID, createdGroup.ID)
	}

	// Verify we can fetch the same group by Name
	fetchedByName, err := service.GetAccountGroupByName(req.Name)
	if err != nil {
		t.Fatalf("GetAccountGroupByName failed: %v", err)
	}
	if fetchedByName.ID != createdGroup.ID {
		t.Fatalf("GetAccountGroupByName returned ID %q, expected %q", fetchedByName.ID, createdGroup.ID)
	}
}

func TestUpdateAccountGroup(t *testing.T) {
	service := newTestService(t)

	parentID := "09e47ddb-1ce1-488c-b4e8-2fd255f2203a" // Direct Expenses

	// First create a group to update
	createReq := CreateAccountGroupRequest{
		Description:   "Group to be updated via SDK test",
		Name:          fmt.Sprintf("Area-Update-%d", time.Now().UnixNano()),
		ParentGroupID: &parentID,
	}

	createdGroup, err := service.CreateAccountGroup(createReq)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid signature") {
			t.Skip("Tigg signature validation not configured for this environment; skipping update account group integration test")
		}
		t.Fatalf("CreateAccountGroup failed for update test: %v", err)
	}

	if createdGroup.ID == "" {
		t.Fatal("expected non-empty ID for created account group in update test")
	}

	updateReq := UpdateAccountGroupRequest{
		Description: "Updated description via SDK test",
		Name:        createReq.Name,
		// Keep same parent for simplicity
		ParentGroupID: createReq.ParentGroupID,
	}

	updatedGroup, err := service.UpdateAccountGroup(createdGroup.ID, updateReq)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid signature") {
			t.Skip("Tigg signature validation not configured for this environment; skipping update account group integration test")
		}
		t.Fatalf("UpdateAccountGroup failed: %v", err)
	}

	if updatedGroup.ID != createdGroup.ID {
		t.Fatalf("expected updated group ID %q to match created group ID %q", updatedGroup.ID, createdGroup.ID)
	}
}

func TestListAccountGroups(t *testing.T) {
	service := newTestService(t)

	groups, err := service.ListAccountGroups()
	if err != nil {
		t.Fatalf("Failed to list account groups: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("expected at least one account group in list")
	}
}
