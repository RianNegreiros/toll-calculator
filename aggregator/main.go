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

	store := NewMemoryStore()
	var (
		svc = newInvoiceAggregator(store)
	)

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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
