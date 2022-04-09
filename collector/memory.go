package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"context"
	"github.com/go-kit/log"
	"os/exec"
)

const (
	//exporter namespace
	namespace = "linux"
	//subsystem,such as cpu,diskspace etc
	mem = "memory"

)

const(
	memQueryShell = `free -h`
	name_memoryInfo ="mem_info"
)
//1.内存数据抓取结构体
type MemScraper struct {}
//2.定义指标
var(
	memUsageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,mem,"memory_usage"),
		"already usage of memory.",
		[]string{"usage"},nil,
		)
	memTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,mem,"memory_total"),
		"total size of memory.",
		[]string{"usage"},nil,
	)
)
func (m MemScraper) Name() string {
	return name_memoryInfo
}

func (m MemScraper) Help() string {
	return "collect memory infomation of linux host"
}

func (m MemScraper) Version() float64 {
	return 7.6
}

func (m MemScraper) Scrape(ctx context.Context, ch chan<- prometheus.Metric, logger log.Logger) error {
	memtotalcmd :=  exec.Command("bash","-c","free -h | awk 'NR==2{print $2}'")

	output,err := memtotalcmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(output)

	return nil
}

//实现Scraper接口
var _ Scraper = &MemScraper{}