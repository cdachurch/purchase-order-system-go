package po

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/sync/errgroup"
)

// PurchaseOrderGetter defines something that can get purchase orders
type PurchaseOrderGetter interface {
	GetPurchaseOrders(ctx context.Context, email string) ([]PurchaseOrder, error)
	ListPurchaseOrders(ctx context.Context, email string, start, length int) (*PagedResponse, error)
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

	return p.getAllPurchaseOrders(ctx, email)
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

// PagedResponse represents a response from the server - it includes data as well as total rows, etc.
type PagedResponse struct {
	POs   []PurchaseOrder
	Total int
}

func (p *poService) ListPurchaseOrders(ctx context.Context, email string, start, length int) (*PagedResponse, error) {
	q := datastore.NewQuery("PurchaseOrder")
	if !shouldAttachEmail(email) {
		email = ""
	}

	if email != "" {
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}

	eg, egCtx := errgroup.WithContext(ctx)
	resp := &PagedResponse{}
	// Get the "total" count
	eg.Go(func() error {
		var q2 *datastore.Query
		if email == "" {
			// If they're getting all POs, look up the first pretty po id and set to that. Close enough!
			q2 = datastore.NewQuery("PurchaseOrder").Limit(1).Order("-pretty_po_id")
			po, err := p.getPOsFromQuery(egCtx, q2)
			if err != nil {
				return err
			}
			resp.Total = po[0].PrettyPoID
			return nil
		}
		// Otherwise, get the count of their total without the limit or offset to get all.
		q2 = q
		total, err := p.dsClient.Count(egCtx, q2)
		if err != nil {
			return err
		}
		resp.Total = total
		return nil
	})

	q = q.Limit(length).Offset(start).Order("-pretty_po_id")

	eg.Go(func() error {
		pos, err := p.getPOsFromQuery(egCtx, q)
		if err != nil {
			return err
		}
		resp.POs = pos
		return nil
	})

	eg.Wait()

	return resp, nil
}
