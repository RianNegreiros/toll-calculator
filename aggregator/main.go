package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/RianNegreiros/toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	grpcListenAddr := flag.String("listen-addr", ":3001", "http listen address")
	httpListenAddr := flag.String("listen-addr", ":3000", "grpc listen address")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(*grpcListenAddr)
	makeHTTPTransport(*httpListenAddr, svc)
}

func makeGRPCTransport(listenAddr string) error {
	fmt.Println("grpc transport listening on", listenAddr)
	ln, err := net.Listen("TCP", listenAddr)
	if err != nil {
		return err
	}
	server := grpc.NewServer([]grpc.ServerOption{}...)
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("http transport listening on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu_id"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu id query param"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
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
