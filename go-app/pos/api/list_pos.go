package poapis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"po/go-app/pos"
	"strconv"
)

type ListPOsResponse struct {
	Status          int                 `json:"status"`
	Data            []pos.PurchaseOrder `json:"data"`
	Draw            int                 `json:"draw"`
	RecordsTotal    int                 `json:"recordsTotal"`
	RecordsFiltered int                 `json:"recordsFiltered"`
}

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
	length, err := strconv.Atoi(r.FormValue("length"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, fmt.Errorf("length: %w", err))
		return
	}
	response, err := s.poService.ListPurchaseOrders(ctx, email, 0, length)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		return
	}

	resp := ListPOsResponse{
		Status:          200,
		Data:            response.POs,
		RecordsTotal:    response.Total,
		RecordsFiltered: response.Total,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
