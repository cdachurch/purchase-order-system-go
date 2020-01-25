// Main API handler for this Golang App Engine module

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	poapis "po/go-app/api"
	"po/go-app/po"
	"strings"

	"cloud.google.com/go/datastore"
)

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	dsClient, err := datastore.NewClient(ctx, getAppIDForDatastore())
	if err != nil {
		log.Printf("Error making datastore client: %v", err)
		return
	}

	handler := po.NewPurchaseOrderGetter(dsClient)
	server := poapis.NewServer(handler)

	mux := http.NewServeMux()
	mux.HandleFunc("/goapi/v1/po/list/", addContext(ctx, server.ListPurchaseOrders))
	mux.HandleFunc("/goapi/v1/po/", addContext(ctx, server.GetPurchaseOrders))
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

// Should be cdac-demo-purchaseorder for demo and cdac-purchaseorder for production
func getAppIDForDatastore() string {
	splitApp := strings.Split(os.Getenv("GAE_APPLICATION"), "~")
	return splitApp[1]
}

func addContext(ctx context.Context, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(ctx)
		fn(w, req)
	}
}
