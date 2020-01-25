package po

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/datastore"
)

// PurchaseOrderGetter defines something that can get purchase orders
type PurchaseOrderGetter interface {
	GetPurchaseOrders(ctx context.Context, email string) []PurchaseOrder
}

// NewPurchaseOrderGetter returns a PurchaseOrderGetter
func NewPurchaseOrderGetter(dsClient *datastore.Client) PurchaseOrderGetter {
	return &poService{dsClient: dsClient}
}

type poService struct {
	dsClient *datastore.Client
}

func (p *poService) GetPurchaseOrders(ctx context.Context, email string) []PurchaseOrder {
	log.Printf("Setting up a new query for POs")
	q := datastore.NewQuery("PurchaseOrder")
	if shouldAttachEmail(email) {
		log.Printf("Using poID: %s", email)
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}
	q = q.Limit(5000)
	log.Printf("Executing query")
	return p.getAllPurchaseOrders(ctx, q)
}

func shouldAttachEmail(email string) bool {
	var userIsAdmin bool
	for _, name := range usersThatCanSeeAllPOs {
		if name == email {
			userIsAdmin = true
			break
		}
	}

	return email != "" && !userIsAdmin
}
