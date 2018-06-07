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

func poHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Setting up datastore client")
	if r.Method == http.MethodGet {
		log.Infof(ctx, "Setting up a new query for POs")
		q := datastore.NewQuery("PurchaseOrder")
		r.ParseForm()
		email := r.FormValue("email")
		if email != "" {
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
			break
		}
		po.CalculateIsAddressed()
		pos = append(pos, po)
	}
	log.Infof(ctx, "Done getting POs, there are %d of them", len(pos))
	return pos
}

//func createPo(ctx context.Context) (*datastore.Entity, error) {
//	po := new(PurchaseOrder)
//	po.PoID = "Whatever"
//	po.PrettyPoID = 1
//	po.Purchaser = "Graham"
//	po.Supplier = "Crystal"
//	po.Product = "Rings"
//	po.Price = 10.40
//
//	key := datastore.NewIncompleteKey(ctx, "PurchaseOrder", nil)
//	_, err := datastore.Put(ctx, key, po)
//	if err != nil {
//		log.Errorf(ctx, err.Error())
//	}
//
//	return nil, nil
//}
