package account

type Account struct {
	ID               string  `json:"id,omitempty"`
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	NameLower        string  `json:"name_lower"`
	Type             string  `json:"type"`
	AccountClassID   string  `json:"account_class_id"`
	AccountClassName string  `json:"account_class_name"`
	PrimaryGroupID   string  `json:"primary_group_id"`
	PrimaryGroupName string  `json:"primary_group_name"`
	ParentGroupID    *string `json:"parent_group_id,omitempty"`
	ParentGroupName  *string `json:"parent_group_name,omitempty"`
	Description      string  `json:"description"`
	Inactive         bool    `json:"inactive"`
	CreatedAt        string  `json:"created_at"`
}

// CreateAccountRequest is the payload used when creating an account.
type CreateAccountRequest struct {
	Name            string  `json:"name"`
	Code            string  `json:"code"`
	ParentGroupID   *string `json:"parent_group_id,omitempty"`
	ParentGroupName *string `json:"parent_group_name,omitempty"`
	Description     string  `json:"description,omitempty"`
}

// UpdateAccountRequest is the payload used when updating an account.
type UpdateAccountRequest struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Code            string  `json:"code"`
	ParentGroupID   *string `json:"parent_group_id,omitempty"`
	ParentGroupName *string `json:"parent_group_name,omitempty"`
	Description     string  `json:"description,omitempty"`
}
