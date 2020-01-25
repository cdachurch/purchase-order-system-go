package po

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
)

// PurchaseOrderGetter defines something that can get purchase orders
type PurchaseOrderGetter interface {
	GetPurchaseOrders(ctx context.Context, email string) ([]PurchaseOrder, error)
}

// NewPurchaseOrderGetter returns a PurchaseOrderGetter
func NewPurchaseOrderGetter(dsClient *datastore.Client) PurchaseOrderGetter {
	return &poService{dsClient: dsClient}
}

type poService struct {
	dsClient *datastore.Client
}

func (p *poService) GetPurchaseOrders(ctx context.Context, email string) ([]PurchaseOrder, error) {
	if email == "" {
		return nil, errors.New("email cannot be blank")
	}
	// If they're an admin just blank out the email so that it won't be used for filtering.
	if !shouldAttachEmail(email) {
		email = ""
	}
	log.Printf("Executing query")
	return p.getAllPurchaseOrders(ctx, email), nil
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
