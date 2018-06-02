// Main API handler for this Golang App Engine module

package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func main() {
	http.HandleFunc("/api/v1", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	po := PurchaseOrder{
		poID:       "Whatever",
		prettyPoID: 1,
		purchaser:  "Graham",
		supplier:   "Crystal",
		product:    "Rings",
		price:      10.40,
	}
	ds, err := datastore.NewClient(ctx, "cdac-purchase-order-demo")
	if err != nil {
		// Blah blah error handling
		log.Errorf(ctx, "Something went wrong!")
	}

	key := datastore.IncompleteKey("PurchaseOrder", nil)
	if _, err := ds.Put(ctx, key, &po); err != nil {
		log.Errorf(ctx, "Could not put new PurchaseOrder.")
	}
	fmt.Fprintln(w, "Hello, Crystal!")
}
