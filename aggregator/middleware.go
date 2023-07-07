package main

import (
	"time"

	"github.com/RianNegreiros/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{next: next}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"func": "AggregateDistance",
		}).Info("Aggregated distance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}
