package poapis

import (
	"net/http"
	"po/go-app/po"
)

type poAPIServer struct {
	poGetter po.PurchaseOrderGetter
}

// APIServer defines an interface through which an http server can handle requests for POs
type APIServer interface {
	GetPurchaseOrders(w http.ResponseWriter, r *http.Request)
}

func NewServer(getter po.PurchaseOrderGetter) APIServer {
	return &poAPIServer{poGetter: getter}
}
