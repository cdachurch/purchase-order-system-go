package pos

import (
	"time"
)

var (
	// TODO: This should come from an env var someday (then we can change these without deploying code)
	usersThatCanSeeAllPOs = []string{
		"smyhre",
		// "gholtslander",
		// "gdholtslander",
		"cbayles",
		"dwiebe",
		// "test@example.com",
	}
)

// ShouldAttachEmail will be true if the email passed in is not an admin's email.
func ShouldAttachEmail(email string) bool {
	var userIsAdmin bool
	for _, name := range usersThatCanSeeAllPOs {
		if name == email {
			userIsAdmin = true
			break
		}
	}

	return email != "" && !userIsAdmin
}

// PurchaseOrder represents a purchase order in the system.
// These are made by users and approved by admins.
type PurchaseOrder struct {
	PoID           string    `json:"po_id" datastore:"po_id"`
	PrettyPoID     int       `json:"pretty_po_id" datastore:"pretty_po_id"`
	Purchaser      string    `json:"purchaser" datastore:"purchaser"`
	Supplier       string    `json:"supplier" datastore:"supplier"`
	Product        string    `json:"product" datastore:"product"`
	Price          float32   `json:"price" datastore:"price"`
	ApprovedBy     string    `json:"approved_by" datastore:"approved_by"`
	AccountCode    string    `json:"account_code" datastore:"account_code"`
	AccountCodeStr string    `json:"account_code_str" datastore:"account_code_str"`
	IsApproved     bool      `json:"is_approved" datastore:"is_approved"`
	IsDenied       bool      `json:"is_denied" datastore:"is_denied"`
	IsInvoiced     bool      `json:"is_invoiced" datastore:"is_invoiced"`
	IsCancelled    bool      `json:"is_cancelled" datastore:"is_cancelled"`
	IsAddressed    bool      `json:"is_addressed" datastore:"is_addressed"`
	Updated        time.Time `datastore:"updated"`
	UpdatedStr     string    `json:"last_updated"`
	Created        time.Time `datastore:"created"`
	CreatedStr     string    `json:"created_date"`
	Deleted        time.Time `datastore:"deleted"`
	DeletedStr     string    `json:"deleted_date"`
}

// CalculateIsAddressed will just set IsAddressed to what it already is if it is set, or it will set it by computation
func (po *PurchaseOrder) CalculateIsAddressed() {
	isAddressed := po.IsAddressed || (po.IsApproved || po.IsCancelled || po.IsDenied)
	po.IsAddressed = isAddressed
}

// FormatDates takes the time.Times on the model and formats them in 2006-01-02 format.
func (po *PurchaseOrder) FormatDates() {
	if time.Time.IsZero(po.Updated) {
		po.UpdatedStr = ""
	}
	po.UpdatedStr = po.Updated.Format("2006-01-02")
	if po.Created.IsZero() {
		po.CreatedStr = ""
	}
	po.CreatedStr = po.Created.Format("2006-01-02")
	if time.Time.IsZero(po.Deleted) {
		po.DeletedStr = ""
	}
	po.DeletedStr = po.Deleted.Format("2006-01-02")
}
