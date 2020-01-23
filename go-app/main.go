// Main API handler for this Golang App Engine module

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"po/go-app/po"

	"cloud.google.com/go/datastore"
)

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	dsClient, err := datastore.NewClient(ctx, "cdac-purchaseorder")
	if err != nil {
		log.Printf("Error making datastore client: %v", err)
		return
	}
	handler := po.NewPurchaseOrderHandler(dsClient)
	mux := http.NewServeMux()
	mux.HandleFunc("/goapi/v1/po/", func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(ctx)
		poHandler(w, req, handler)
	})
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func poHandler(w http.ResponseWriter, r *http.Request, handler po.PurchaseOrderHandler) {
	ctx := r.Context()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
		}
		email := r.FormValue("email")

		pos := handler.GetPurchaseOrders(ctx, email)

		resp := map[string]interface{}{
			"status": 200,
			"data":   pos,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
	// Hitting us with an unsupported method
	w.WriteHeader(405)
	return
}
