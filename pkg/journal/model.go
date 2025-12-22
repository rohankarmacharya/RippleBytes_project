package journal

type VoucherStatus string

const (
	statusDraft  VoucherStatus = "DRAFT"
	statusPosted VoucherStatus = "POSTED"
	statusVoided VoucherStatus = "VOIDED"
)

type JournalVoucher struct {
	ID            string               `json:"id,omitempty"`
	Code          string               `json:"code"`
	Date          string               `json:"date"`
	CurrencyCode  string               `json:"currency_code"`
	VoucherStatus VoucherStatus        `json:"status,omitempty"`
	Narration     string               `json:"narration,omitempty"`
	Items         []JournalVoucherItem `json:"items"`
	CreatedAt     string               `json:"created_at,omitempty"`
	UpdatedAt     string               `json:"updated_at,omitempty"`
}

type JournalVoucherItem struct {
	AccountID   string `json:"account_id,omitempty"`
	AccountCode string `json:"account_code,omitempty"`
	Amount      string `json:"amount"`
	TxnType     string `json:"txn_type"`
	Narration   string `json:"narration,omitempty"`
}
