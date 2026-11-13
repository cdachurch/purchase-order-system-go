package porepository

import (
	"context"
	"fmt"
	"po/go-app/pos"
)

// ListPurchaseOrders retrieves a paginated list of purchase orders for the given email.
// It's used on the "list POs" page that most people use and never actually gets paginated :)
func (p *poRepository) ListPurchaseOrders(ctx context.Context, email string, start, length int) (*pos.PagedResponse, error) {
	q := buildPOQuery(email)

	purchaseOrders, err := p.getPOsFromQuery(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error getting purchaseOrders for individual: %v", err)
	}
	resp := &pos.PagedResponse{
		POs:   purchaseOrders,
		Total: len(purchaseOrders),
	}

	return resp, nil
}
