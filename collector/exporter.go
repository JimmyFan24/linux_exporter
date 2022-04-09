package collector

import (
	"context"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)


type Metrics struct {
	TotalScrapes prometheus.Counter
	ScrapeErrors *prometheus.CounterVec
	Error prometheus.Gauge
	LinuxExporterUp prometheus.Gauge
}
type Exporter struct {
	ctx      context.Context
	logger   log.Logger
	scrapers []Scraper
	metrics  Metrics
}

func (e Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.LinuxExporterUp.Desc()
}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(e.ctx, ch)
	e.metrics.ScrapeErrors.Collect(ch)
	ch <- e.metrics.LinuxExporterUp

}

func (e Exporter)scrape(ctx context.Context,ch chan<- prometheus.Metric){

}

//创建Exporter
func New(ctx context.Context,metric Metrics,scrapes []Scrapes,logger log.Logger)*Exporter{

}
//创建指标
func NewMetrics()Metrics{
	subsystem := "exporter"
	return Metrics{
		TotalScrapes:prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name: "scrapes_total",
			Help: "Total number of times linux was scraped for metrics.",
		}),

		ScrapeErrors:prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem:subsystem,
			Name:"scrape_errors_total",
			Help:"Total number of times an error occurred scraping a Linux.",
		},[]string{"collector"}),
		Error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from Linux resulted in an error (1 for error, 0 for success).",
		}),
		LinuxExporterUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether the LinuxExporter server is up.",
		}),
	}

}
