package client

import (
	"context"

	"github.com/RianNegreiros/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
