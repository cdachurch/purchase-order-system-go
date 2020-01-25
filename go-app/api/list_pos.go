package poapis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (s *poAPIServer) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "something went wrong: %v", err)
		return
	}
	email := r.FormValue("email")
	draw, err := strconv.Atoi(r.FormValue("draw"))
	start, err := strconv.Atoi(r.FormValue("start"))
	length, err := strconv.Atoi(r.FormValue("length"))
	response, err := s.poGetter.ListPurchaseOrders(ctx, email, start, length)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "something bad happened: %v", err)
		return
	}

	resp := map[string]interface{}{
		"status":          200,
		"data":            response.POs,
		"draw":            draw,
		"recordsTotal":    response.Total,
		"recordsFiltered": response.Total,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	return
}
