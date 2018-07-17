// Main API handler for this Golang App Engine module

package main

import (
	"context"
	"encoding/json"
	"net/http"

	"strings"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/goapi/v1/po/", poHandler)
}

func shouldAttachEmail(email string) bool {
	usersThatCanSeeAllPOs := []string{
		"gdholtslander",
		"gholtslander",
		"smyhre",
		"dwiebe",
		"test@example.com",
		"jheindle",
		"rhoult",
		"rsmith",
	}

	var userIsAdmin bool
	for _, name := range usersThatCanSeeAllPOs {
		if name == email {
			userIsAdmin = true
			break
		}
	}

	return email != "" && !userIsAdmin
}

func poHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.Method == http.MethodGet {
		r.ParseForm()
		email := r.FormValue("email")
		log.Infof(ctx, "Setting up a new query for POs")
		q := datastore.NewQuery("PurchaseOrder")
		if shouldAttachEmail(email) {
			log.Infof(ctx, "Using poID: %s", email)
			tokens := strings.Split(email, "@")
			q = q.Filter("purchaser =", tokens[0])
		}
		q = q.BatchSize(5000)
		log.Infof(ctx, "Executing query")
		pos := getAllPurchaseOrders(ctx, q)

		resp := map[string]interface{}{
			"status": 200,
			"data":   pos,
		}
		json.NewEncoder(w).Encode(resp)
	}
	return
}

func getAllPurchaseOrders(ctx context.Context, q *datastore.Query) []PurchaseOrder {
	var pos []PurchaseOrder
	log.Infof(ctx, "About to get POs")
	for t := q.Run(ctx); ; {
		var po PurchaseOrder
		_, err := t.Next(&po)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error somehow. Skip it maybe?
			log.Infof(ctx, "Error received: %s", err.Error())
			break
		}
		po.CalculateIsAddressed()
		pos = append(pos, po)
	}
	log.Infof(ctx, "Done getting POs, there are %d of them", len(pos))
	return pos
}
