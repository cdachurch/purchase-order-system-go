package main

import "time"

type PurchaseOrder struct {
	PoID        string    `json:"po_id" datastore:"po_id"`
	PrettyPoID  int       `json:"pretty_po_id" datastore:"pretty_po_id"`
	Purchaser   string    `json:"purchaser" datastore:"purchaser"`
	Supplier    string    `json:"supplier" datastore:"supplier"`
	Product     string    `json:"product" datastore:"product"`
	Price       float32   `json:"price" datastore:"price"`
	ApprovedBy  string    `json:"approved_by" datastore:"approved_by"`
	AccountCode string    `json:"account_code" datastore:"account_code"`
	IsApproved  bool      `json:"is_approved" datastore:"is_approved"`
	IsDenied    bool      `json:"is_denied" datastore:"is_denied"`
	IsInvoiced  bool      `json:"is_invoiced" datastore:"is_invoiced"`
	IsCancelled bool      `json:"is_cancelled" datastore:"is_cancelled"`
	IsAddressed bool      `json:"is_addressed" datastore:"is_addressed"`
	Updated     time.Time `json:"last_updated" datastore:"updated"`
	Created     time.Time `json:"created_date" datastore:"created"`
	Deleted     time.Time `json:"deleted_date" datastore:"deleted"`
}

// CalculateIsAddressed will just set IsAddressed to what it already is if it is set, or it will set it by computation
func (po *PurchaseOrder) CalculateIsAddressed() {
	isAddressed := po.IsAddressed || (po.IsApproved || po.IsCancelled || po.IsDenied)
	po.IsAddressed = isAddressed
}
