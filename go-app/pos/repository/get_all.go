package porepository

import (
	"context"
	"log"
	"strings"

	"po/go-app/pos"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func (p *poRepository) GetPurchaseOrders(ctx context.Context, email string) ([]pos.PurchaseOrder, error) {
	q := datastore.NewQuery("PurchaseOrder")
	if email != "" {
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}

	log.Printf("About to get POs")
	purchaseOrders, err := p.getPOsFromQuery(ctx, q)
	log.Printf("Done getting POs, there are %d of them", len(purchaseOrders))
	return purchaseOrders, err
}
