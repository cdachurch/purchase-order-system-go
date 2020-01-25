package po

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/datastore"
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
	Total int64
}

func (p *poService) ListPurchaseOrders(ctx context.Context, email string, start, length int) (*PagedResponse, error) {
	q := datastore.NewQuery("PurchaseOrder").Limit(length).Offset(start)
	if !shouldAttachEmail(email) {
		email = ""
	}

	if email != "" {
		tokens := strings.Split(email, "@")
		q = q.Filter("purchaser =", tokens[0])
	}
	pos, err := p.getPOsFromQuery(ctx, q)
	if err != nil {
		return nil, err
	}

	return &PagedResponse{
		POs:   pos,
		Total: 27,
	}, nil
}
