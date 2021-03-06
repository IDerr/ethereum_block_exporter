package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/iderr/ethereum_block_exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "undefined"

func main() {
	flag.Usage = func() {
		const (
			usage = "Usage: ethereum_exporter [option] [arg]\n\n" +
				"Prometheus exporter for Ethereum client metrics\n\n" +
				"Options and arguments:\n"
		)

		fmt.Fprint(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()

		os.Exit(2)
	}

	url := flag.String("url", "http://localhost:8545", "Ethereum JSON-RPC URL")
	addr := flag.String("addr", ":9368", "listen address")
	ver := flag.Bool("v", false, "print version number and exit")

	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
	}

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	rpc, err := rpc.Dial(*url)
	if err != nil {
		log.Fatal(err)
	}

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		collector.NewEthBlockNumber(rpc),
	)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
		ErrorHandling: promhttp.ContinueOnError,
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Ethereum Block Exporter</title></head>
            <body>
            <h1>Ethereum Block Exporter</h1>
            <p><a href='/metrics'>Metrics</a></p>
            </body>
            </html>`))
	})
	http.Handle("/metrics", handler)
	log.Println(`Listening on ` + *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
