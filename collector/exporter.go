package collector

import (
	"context"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// Metric descriptors.
var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "exporter", "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

type Metrics struct {
	TotalScrapes    prometheus.Counter
	ScrapeErrors    *prometheus.CounterVec
	Error           prometheus.Gauge
	LinuxExporterUp prometheus.Gauge
}
type Exporter struct {
	ctx      context.Context
	logger   log.Logger
	scrapers []Scraper
	metrics  Metrics
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.LinuxExporterUp.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(e.ctx, ch)
	e.metrics.ScrapeErrors.Collect(ch)
	ch <- e.metrics.LinuxExporterUp

}

//通过类型转换方式来检查是否实现接口
var _ prometheus.Collector = &Exporter{}

func (e *Exporter) scrape(ctx context.Context, ch chan<- prometheus.Metric) {
	e.metrics.TotalScrapes.Inc()

	scrapeTime := time.Now()

	e.metrics.LinuxExporterUp.Set(1)
	e.metrics.Error.Set(0)
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "connection")

	var wg sync.WaitGroup
	defer wg.Wait()
	for _, scraper := range e.scrapers {
		wg.Add(1)
		go func(scraper Scraper) {
			defer wg.Done()
			label := "collect." + scraper.Name()
			scrapeTime := time.Now()
			if err := scraper.Scrape(ctx, ch, log.With(e.logger, "scrapers", scraper.Name())); err != nil {
				logrus.Errorf("Error from scraper:%v,err:%v", scraper.Name(), err)
				e.metrics.ScrapeErrors.WithLabelValues(label).Inc()
				e.metrics.Error.Set(1)
			}
			ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), label)

		}(scraper)
	}
}

//创建Exporter
func New(ctx context.Context, metric Metrics, scrapers []Scraper, logger log.Logger) *Exporter {
	return &Exporter{
		ctx:      ctx,
		metrics:  metric,
		scrapers: scrapers,
		logger:   logger,
	}
}

//创建指标
func NewMetrics() Metrics {
	subsystem := "exporter"
	return Metrics{
		TotalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "scrapes_total",
			Help:      "Total number of times linux was scraped for metrics.",
		}),

		ScrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occurred scraping a Linux.",
		}, []string{"collector"}),
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
