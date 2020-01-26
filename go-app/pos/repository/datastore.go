package porepository

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"po/go-app/pos"
)

// New returns a purchase order getter
func New(dsClient *datastore.Client) pos.PurchaseOrderGetter {
	return &poRepository{dsClient: dsClient}
}

type poRepository struct {
	dsClient *datastore.Client
}

// Exhaust a query and return all the purchase orders attained.
func (p *poRepository) getPOsFromQuery(ctx context.Context, q *datastore.Query) ([]pos.PurchaseOrder, error) {
	var purchaseOrders []pos.PurchaseOrder
	for t := p.dsClient.Run(ctx, q); ; {
		var po pos.PurchaseOrder
		_, err := t.Next(&po)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("error from querying datastore: %s", err.Error())
			return nil, err
		}
		// "Calculated fields"
		po.CalculateIsAddressed()
		po.FormatDates()
		purchaseOrders = append(purchaseOrders, po)
	}
	return purchaseOrders, nil
}
