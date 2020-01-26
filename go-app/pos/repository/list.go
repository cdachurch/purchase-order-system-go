package porepository

import (
	"context"
	"fmt"
	"strings"
	"po/go-app/pos"
	"cloud.google.com/go/datastore"
	"golang.org/x/sync/errgroup"
)

func (p *poRepository) ListPurchaseOrders(ctx context.Context, email string, start, length int) (*pos.PagedResponse, error) {
	q := datastore.NewQuery("PurchaseOrder")
	if email != "" {
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}

	eg, egCtx := errgroup.WithContext(ctx)
	resp := &pos.PagedResponse{}
	// Get the "total" count
	eg.Go(func() error {
		var q2 *datastore.Query
		if email == "" {
			// If they're getting all POs, look up the first pretty po id and set to that. Close enough!
			q2 = datastore.NewQuery("PurchaseOrder").Limit(1).Order("-pretty_po_id")
			po, err := p.getPOsFromQuery(egCtx, q2)
			if err != nil {
				return fmt.Errorf("error getting count of pos: %v", err)
			}
			resp.Total = po[0].PrettyPoID
			return nil
		}
		// Otherwise, get the count of their total without the limit or offset to get all.
		q2 = q
		total, err := p.dsClient.Count(egCtx, q2)
		if err != nil {
			return fmt.Errorf("error getting count of pos for individual: %v", err)
		}
		resp.Total = total
		return nil
	})

	q = q.Limit(length).Offset(start).Order("-pretty_po_id")

	eg.Go(func() error {
		purchaseOrders, err := p.getPOsFromQuery(egCtx, q)
		if err != nil {
			return fmt.Errorf("error getting purchaseOrders for individual: %v", err)
		}
		resp.POs = purchaseOrders
		return nil
	})

	err := eg.Wait()

	if err != nil {
		return nil, err
	}

	return resp, nil
}
