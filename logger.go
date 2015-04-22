package main

import (
	"fmt"
	"github.com/marpaia/graphite-golang"
	"log"
)

const (
	GRAPHITE_HOST = "54.191.107.29"
	GRAPHITE_PORT = 2003
)

var logSrv *graphite.Graphite

func init() {
	var err error
	logSrv, err = graphite.NewGraphite(GRAPHITE_HOST, GRAPHITE_PORT)
	if err != nil {
		log.Fatalf("Could not connect to Graphite server: %s\n", err)
	}
}

func Logger(logger_queue <-chan RfRecord) {
	var name string
	var value string

	for record := range logger_queue {
		name = fmt.Sprintf("frequency.%f.power", record.frequency)
		value = fmt.Sprintf("%f", record.power)

		log.Printf("Sending %s %s to graphite\n", name, value)
		logSrv.SendMetric(graphite.NewMetric(name, value, record.recorded_at.Unix()))
	}
}
