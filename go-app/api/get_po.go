package poapis

import (
	"encoding/json"
	"log"
	"net/http"
)

// GetPurchaseOrders handles http requests for /goapi/v1/po/
// In general, it wants to see an email param in the query parameters. If it doesn't see one
// and the user is an administrator, they'll see all the pos they have access to.
func (s *poAPIServer) GetPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodGet {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
		}
		email := r.FormValue("email")

		pos := s.poGetter.GetPurchaseOrders(ctx, email)

		resp := map[string]interface{}{
			"status": 200,
			"data":   pos,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}
	// Hitting us with an unsupported method
	w.WriteHeader(405)
	return
}
