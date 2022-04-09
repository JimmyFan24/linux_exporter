package main

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/alecthomas/kingpin.v2"
	"linux_exporter/collector"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//1.启动参数

var (
	listenAddress = pflag.String("listen-addr", "127.0.0.1:9999", "Address to listen on for web interface and telemetry.")
	metricPath    = pflag.String("web-telemetry-path", "/metrics", "Path under which to expose metrics.")
)

//2.定义数据抓取的scrape的集合
var scrapers = map[collector.Scraper]bool{
	collector.MemScraper{}: true,
}

func main() {
	scraperFlag := map[collector.Scraper]*bool{}
	for scraper, enableByDefault := range scrapers {
		defaultOn := "false"
		if enableByDefault {
			defaultOn = "true"
		}
		strDefaultOn, _ := strconv.ParseBool(defaultOn)
		f := pflag.Bool("collect."+scraper.Name(), strDefaultOn, scraper.Help())
		scraperFlag[scraper] = f

	}

	//parse flags
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("liunx_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)
	// landingPage contains the HTML served at '/'.

	var landingPage = []byte(`<html>
<head><title>MySQLd exporter</title></head>
<body>
<h1>MySQLd exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)
	logrus.Infof("Starting linux_exporter,version is %v", version.Info())
	logrus.Infof("Build context,%v", version.BuildContext())

	// Register only scrapers enabled by flag.
	enabledScrapers := []collector.Scraper{}
	for scraper, enable := range scraperFlag {
		if *enable {
			logrus.Infof("Scraper enabled:%v", scraper.Name())
			enabledScrapers = append(enabledScrapers, scraper)
		}
	}

	handlerFunc := newHandler(collector.NewMetrics(), enabledScrapers, logger)
	http.Handle(*metricPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	//exporter := collector.New(ctx)
	//http.Handle(*metricPath,promhttp.Handler())
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(landingPage)
	})

	srv := &http.Server{
		Addr: *listenAddress,
	}
	if err := web.ListenAndServe(srv, "", logger); err != nil {
		logrus.Errorf("Error starting HTTP server:%v", err)
		os.Exit(1)
	}

}
func newHandler(metric collector.Metrics, scrapers []collector.Scraper, logger log.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logrus.Info("newhandler func running")
		filteredScrapers := scrapers
		params := request.URL.Query()["collect[]"]
		logrus.Debugf("collect[] params:%v", strings.Join(params, ","))
		ctx := request.Context()
		registry := prometheus.NewRegistry()
		registry.MustRegister(collector.New(ctx, metric, filteredScrapers, logger))
		gatherers := prometheus.Gatherers{
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(writer, request)
	}
}
