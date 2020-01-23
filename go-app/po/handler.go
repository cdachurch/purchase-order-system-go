package po

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/datastore"
)

// PurchaseOrderHandler defines something that can get purchase orders
type PurchaseOrderHandler interface {
	GetPurchaseOrders(ctx context.Context, email string) []PurchaseOrder
}

// NewPurchaseOrderHandler returns a PurchaseOrderHandler
func NewPurchaseOrderHandler(dsClient *datastore.Client) PurchaseOrderHandler {
	return &poHandler{dsClient: dsClient}
}

type poHandler struct {
	dsClient *datastore.Client
}

func (p *poHandler) GetPurchaseOrders(ctx context.Context, email string) []PurchaseOrder {
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
