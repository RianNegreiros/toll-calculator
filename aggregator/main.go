package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/RianNegreiros/toll-calculator/types"
)

func main() {
	listenAddr := flag.String("listen-addr", ":3000", "server listen address")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = newInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
	fmt.Println("HTTP Trasnport listening on", listenAddr)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
