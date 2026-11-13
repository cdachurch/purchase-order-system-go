package porepository

import (
	"context"
	"log"

	"po/go-app/pos"
)

// GetPurchaseOrders returns allllll of the POs and can be slow.
func (p *poRepository) GetPurchaseOrders(ctx context.Context, email string) ([]pos.PurchaseOrder, error) {
	q := buildPOQuery(email)

	log.Printf("About to get POs")
	purchaseOrders, err := p.getPOsFromQuery(ctx, q)
	log.Printf("Done getting POs, there are %d of them", len(purchaseOrders))
	return purchaseOrders, err
}
