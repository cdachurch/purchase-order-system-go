package po

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func (p *poService) getAllPurchaseOrders(ctx context.Context, q *datastore.Query) []PurchaseOrder {
	var pos []PurchaseOrder
	log.Printf("About to get POs")
	for t := p.dsClient.Run(ctx, q); ; {
		var po PurchaseOrder
		_, err := t.Next(&po)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error somehow. Skip it maybe?
			log.Printf("Error received: %s", err.Error())
			break
		}
		// "Calculated fields"
		po.calculateIsAddressed()
		po.formatDates()
		pos = append(pos, po)
	}
	log.Printf("Done getting POs, there are %d of them", len(pos))
	return pos
}
