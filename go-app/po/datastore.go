package po

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func (p *poService) getAllPurchaseOrders(ctx context.Context, email string) ([]PurchaseOrder, error) {
	q := datastore.NewQuery("PurchaseOrder")
	if email != "" {
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}
	q = q.Limit(5000)

	log.Printf("About to get POs")
	pos, err := p.getPOsFromQuery(ctx, q)
	log.Printf("Done getting POs, there are %d of them", len(pos))
	return pos, err
}

// Exhaust a query and return all the purchase orders attained.
func (p *poService) getPOsFromQuery(ctx context.Context, q *datastore.Query) ([]PurchaseOrder, error) {
	var pos []PurchaseOrder
	for t := p.dsClient.Run(ctx, q); ; {
		var po PurchaseOrder
		_, err := t.Next(&po)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error somehow. Skip it maybe?
			log.Printf("Error received: %s", err.Error())
			return nil, err
		}
		// "Calculated fields"
		po.calculateIsAddressed()
		po.formatDates()
		pos = append(pos, po)
	}
	return pos, nil
}
