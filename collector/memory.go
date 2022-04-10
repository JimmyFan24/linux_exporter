package collector

import (
	"context"
	"fmt"
	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	//exporter namespace
	namespace = "linux"
	//subsystem,such as cpu,diskspace etc
	mem = "memory"
)

const (
	memQueryShell   = `free -h`
	name_memoryInfo = "mem_info"
)

//1.内存数据抓取结构体
type MemScraper struct{}

//2.定义指标
var (
	memUsageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, mem, "memory_usage"),
		"already usage of memory.",
		[]string{"usage"}, nil,
	)
	memTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, mem, "memory_total"),
		"total size of memory.",
		[]string{"usage"}, nil,
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
	//memtotalcmd :=  exec.Command("bash","-c","free -h | awk 'NR==2{print $2}'")
	//memtotalcmd :=  exec.Command("bash","-c","ipconfig")
	back := &backend.Local{}
	shell, err := ps.New(back)
	if err != nil {
		panic(err)
	}
	defer shell.Exit()

	// ... 和它交互
	stdout, _, err := shell.Execute("ipconfig")
	if err != nil {
		panic(err)
	}
	ch <- prometheus.MustNewConstMetric(
		newDesc("memory", "meminfo", "memory info fo windows host"),
		prometheus.UntypedValue,
		88.888)
	fmt.Println(stdout)
	return nil
	/*output,err := memtotalcmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(output)

	return nil*/
}

//实现Scraper接口
var _ Scraper = &MemScraper{}
