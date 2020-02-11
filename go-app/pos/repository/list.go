package porepository

import (
	"context"
	"fmt"
	"po/go-app/pos"
	"strings"

	"cloud.google.com/go/datastore"
)

func (p *poRepository) ListPurchaseOrders(ctx context.Context, email string, start, length int) (*pos.PagedResponse, error) {
	q := datastore.NewQuery("PurchaseOrder").Limit(length).Offset(start).Order("-pretty_po_id")
	if email != "" {
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}

	purchaseOrders, err := p.getPOsFromQuery(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("error getting purchaseOrders for individual: %v", err)
	}
	resp := &pos.PagedResponse{
		POs: purchaseOrders,
	}

	return resp, nil
}
