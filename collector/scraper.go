package collector

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/log"
)

//interface of scraper

type Scraper interface {

	Name()string

	Help()string

	Version()float64

	Scrape(ctx context.Context,ch chan<-prometheus.Metric,logger log.Logger)error
}

