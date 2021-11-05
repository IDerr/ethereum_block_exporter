package collector

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthBlockNumber struct {
	rpc    *rpc.Client
	desc   *prometheus.Desc
	upDesc *prometheus.Desc
}

func NewEthBlockNumber(rpc *rpc.Client) *EthBlockNumber {
	return &EthBlockNumber{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_block_number",
			"number of the most recent block",
			nil,
			nil,
		),
		upDesc: prometheus.NewDesc(
			"eth_block_up",
			"has the rpc call succeded",
			nil,
			nil,
		),
	}
}

func (collector *EthBlockNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
	ch <- collector.upDesc
}

func (collector *EthBlockNumber) Collect(ch chan<- prometheus.Metric) {
	var result hexutil.Uint64
	if err := collector.rpc.Call(&result, "eth_blockNumber"); err != nil {
		ch <- prometheus.MustNewConstMetric(collector.upDesc, prometheus.GaugeValue, 0)
		fmt.Println(err)
		return
	}

	value := float64(result)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
	ch <- prometheus.MustNewConstMetric(collector.upDesc, prometheus.GaugeValue, 1)
}
