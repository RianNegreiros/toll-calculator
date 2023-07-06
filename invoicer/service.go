package main

import "github.com/RianNegreiros/toll-calculator/types"

type Invoicer interface {
	AggregateDistance(types.Distance) error
}
