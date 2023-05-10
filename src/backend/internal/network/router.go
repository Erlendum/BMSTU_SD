package network

import (
	"backend/internal/network/handlers"
	"net/http"
)

func NewRouter(h *handlers.Handlers) *http.ServeMux {
	mux := &http.ServeMux{}
	mux.HandleFunc("/create_instrument", http.HandlerFunc(h.InstrumentHandler.Create))
	mux.HandleFunc("/instruments", http.HandlerFunc(h.InstrumentHandler.GetList))
	mux.HandleFunc("/delete_instrument", http.HandlerFunc(h.InstrumentHandler.Delete))
	mux.HandleFunc("/update_instrument", http.HandlerFunc(h.InstrumentHandler.Update))
	mux.HandleFunc("/create_user", http.HandlerFunc(h.UserHandler.Create))
	mux.HandleFunc("/get_user", http.HandlerFunc(h.UserHandler.Get))
	mux.HandleFunc("/comparison_list", http.HandlerFunc(h.UserHandler.GetComparisonList))
	mux.HandleFunc("/add_instrument_to_comparison_list", http.HandlerFunc(h.ComparisonListHandler.AddInstrument))
	mux.HandleFunc("/delete_instrument_from_comparison_list", http.HandlerFunc(h.ComparisonListHandler.DeleteInstrument))
	return mux
}
