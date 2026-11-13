package porepository

import (
	"po/go-app/pos"
	"strings"

	"cloud.google.com/go/datastore"
)

func buildPOQuery(email string) *datastore.Query {
	q := datastore.NewQuery("PurchaseOrder")
	if email != "" {
		tokens := strings.Split(email, "@")
		// Facilities users need to see all of their POs AS WELL AS
		// any POs marked for as facilities POs
		if pos.FacilitiesUser(email) {
			facilitiesFilter := datastore.OrFilter{
				Filters: []datastore.EntityFilter{
					datastore.PropertyFilter{
						FieldName: "for_facilities",
						Operator:  "=",
						Value:     true,
					},
					datastore.PropertyFilter{
						FieldName: "purchaser",
						Operator:  "=",
						Value:     tokens[0],
					},
				},
			}
			q = q.FilterEntity(facilitiesFilter)
		} else {
			q = q.FilterField("purchaser", "=", tokens[0])
		}
	}
	return q
}
