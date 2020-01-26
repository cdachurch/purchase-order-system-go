package poservice

import (
	"context"
	"errors"
	"fmt"

	"po/go-app/pos"
)

// New returns a PurchaseOrderGetter
func New(poRepo pos.PurchaseOrderGetter) pos.PurchaseOrderGetter {
	// dsClient *datastore.Client
	return &poService{poRepo: poRepo}
}

type poService struct {
	poRepo pos.PurchaseOrderGetter
}

func (p *poService) GetPurchaseOrders(ctx context.Context, email string) ([]pos.PurchaseOrder, error) {
	if email == "" {
		return nil, errors.New("email cannot be blank")
	}
	// If they're an admin just blank out the email so that it won't be used for filtering.
	if !pos.ShouldAttachEmail(email) {
		email = ""
	}

	return p.poRepo.GetPurchaseOrders(ctx, email)
}

func (p *poService) ListPurchaseOrders(ctx context.Context, email string, start, length int) (*pos.PagedResponse, error) {
	if email == "" {
		return nil, fmt.Errorf("email must be specified")
	}
	if !pos.ShouldAttachEmail(email) {
		email = ""
	}

	resp, err := p.poRepo.ListPurchaseOrders(ctx, email, start, length)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
