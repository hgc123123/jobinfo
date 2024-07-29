package main

import (
	"flag"
	"net/http"

	"github.com/hgc123123/jobinfo/transfer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const officialPort = "9821"

func main() {
	portPtr := flag.String("port", officialPort, "the port to listen on")
	helpPtr := flag.Bool("help", false, "print usage")
	cgroupsRootPathPtr := flag.String("cgroups-root", "/sys/fs/cgroup", "path to the root of the cgroupsv1 hierarchy")
	methodPtr := flag.String("method", "slurm", "one of: file, slurm")
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}
	log.Printf("serving cgroups from hierarchy root %s", *cgroupsRootPathPtr)
	// Create and register our cgroups collector
	var cgroupsCollector prometheus.Collector
	if *methodPtr == "slurm" {
		cgroupsCollector = collectors.NewCgroupsSlurmCollector(*cgroupsRootPathPtr)
	} else {
		log.Fatalf("invalid method %s", *methodPtr)
	}
	prometheus.MustRegister(cgroupsCollector)
	// Serve Prometheus HTTP requests
	log.Printf("listening on port %s", *portPtr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+(*portPtr), nil))
}
