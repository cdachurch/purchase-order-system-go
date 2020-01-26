package pos

import "context"

// PurchaseOrderGetter defines something that can get purchase orders
type PurchaseOrderGetter interface {
	GetPurchaseOrders(ctx context.Context, email string) ([]PurchaseOrder, error)
	ListPurchaseOrders(ctx context.Context, email string, start, length int) (*PagedResponse, error)
}

// PagedResponse represents a response from the server - it includes data as well as total rows, etc.
type PagedResponse struct {
	POs   []PurchaseOrder
	Total int
}
