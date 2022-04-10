package collector

import (
	"context"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type CpuScraper struct{}

func (c CpuScraper) Name() string {
	return "cpu_info"
}

func (c CpuScraper) Help() string {
	return "collect cpu infomation of linux host"
}

func (c CpuScraper) Version() float64 {
	return 7.6
}

func (c CpuScraper) Scrape(ctx context.Context, ch chan<- prometheus.Metric, logger log.Logger) error {
	ch <- prometheus.MustNewConstMetric(
		newDesc("cpu", "cpuinfo", "cpu info fo windows host"),
		prometheus.UntypedValue,
		0.80)
	return nil
}

var _ Scraper = &CpuScraper{}
