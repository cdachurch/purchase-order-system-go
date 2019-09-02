// Main API handler for this Golang App Engine module

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"strings"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

var (
	usersThatCanSeeAllPOs = []string{
		"gdholtslander",
		"gholtslander",
		"smyhre",
		"dwiebe",
		"test@example.com",
		"jheindle",
		"rhoult",
		"rsmith",
	}
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
	mux := http.NewServeMux()
	mux.HandleFunc("/goapi/v1/po/", func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(ctx)
		poHandler(w, req, dsClient)
	})
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
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

func poHandler(w http.ResponseWriter, r *http.Request, dsClient *datastore.Client) {
	ctx := r.Context()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
		}
		email := r.FormValue("email")
		log.Printf("Setting up a new query for POs")
		q := datastore.NewQuery("PurchaseOrder")
		if shouldAttachEmail(email) {
			log.Printf("Using poID: %s", email)
			tokens := strings.Split(email, "@")
			q = q.Filter("purchaser =", tokens[0])
		}
		q = q.Limit(5000)
		log.Printf("Executing query")
		pos := getAllPurchaseOrders(ctx, dsClient, q)

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

func getAllPurchaseOrders(ctx context.Context, dsClient *datastore.Client, q *datastore.Query) []PurchaseOrder {
	var pos []PurchaseOrder
	log.Printf("About to get POs")
	for t := dsClient.Run(ctx, q); ; {
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
		po.CalculateIsAddressed()
		po.FormatDates()
		pos = append(pos, po)
	}
	log.Printf("Done getting POs, there are %d of them", len(pos))
	return pos
}
