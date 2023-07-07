package main

import "github.com/RianNegreiros/toll-calculator/types"

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func newInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{store}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	return i.store.Insert(distance)
}
