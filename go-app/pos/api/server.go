package poapis

import (
	"net/http"
	"po/go-app/pos"
)

type poAPIServer struct {
	poService pos.PurchaseOrderGetter
}

// Server defines an interface through which an http server can handle requests for POs
type Server interface {
	GetPurchaseOrders(w http.ResponseWriter, r *http.Request)
	ListPurchaseOrders(w http.ResponseWriter, r *http.Request)
}

// NewServer returns a Server for handling requests.
func NewServer(getter pos.PurchaseOrderGetter) Server {
	return &poAPIServer{poService: getter}
}
