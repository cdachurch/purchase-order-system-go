package main

import "time"

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

func (po *PurchaseOrder) FormatDates() {
	if time.Time.IsZero(po.Updated) {
		po.UpdatedStr = ""
	}
	po.UpdatedStr = po.Updated.Format("2018-05-23")
	if time.Time.IsZero(po.Created) {
		po.CreatedStr = ""
	}
	po.CreatedStr = po.Updated.Format("2018-05-23")
	if time.Time.IsZero(po.Deleted) {
		po.DeletedStr = ""
	}
	po.DeletedStr = po.Deleted.Format("2018-05-23")
}
