package accountgroup

// AccountGroup represents the response model returned by the Tigg API.
type AccountGroup struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	NameLower string `json:"name_lower"`

	AccountClassID   string `json:"account_class_id"`
	AccountClassName string `json:"account_class_name"`

	PrimaryGroupID   string `json:"primary_group_id"`
	PrimaryGroupName string `json:"primary_group_name"`

	ParentGroupID   *string `json:"parent_group_id,omitempty"`
	ParentGroupName *string `json:"parent_group_name,omitempty"`

	Description string `json:"description"`
	Inactive    bool   `json:"inactive"`

	CreatedAt string `json:"created_at"`
}

// CreateAccountGroupRequest is the payload used when creating an account group.
type CreateAccountGroupRequest struct {
	Description     string  `json:"description,omitempty"`
	Name            string  `json:"name"`
	ParentGroupID   *string `json:"parent_group_id,omitempty"`
	ParentGroupName *string `json:"parent_group_name,omitempty"`
}

// UpdateAccountGroupRequest is the payload used when updating an account group.
type UpdateAccountGroupRequest struct {
	ID              string  `json:"id"`
	Description     string  `json:"description,omitempty"`
	Name            string  `json:"name"`
	ParentGroupID   *string `json:"parent_group_id,omitempty"`
	ParentGroupName *string `json:"parent_group_name,omitempty"`
}
