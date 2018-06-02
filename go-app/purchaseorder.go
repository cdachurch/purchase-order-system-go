package main

type PurchaseOrder struct {
	poID        string
	prettyPoID  int
	purchaser   string
	supplier    string
	product     string
	price       float32
	approvedBy  string
	accountCode string
	isApproved  bool
	isDenied    bool
	isInvoiced  bool
	isCancelled bool
}
