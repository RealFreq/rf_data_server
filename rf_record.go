package main

import (
	"time"
)

type RfRecord struct {
	latitude         string
	longitude        string
	recorded_at      time.Time
	hz_start         string
	hz_end           string
	hz_step          string
	averaged_samples string
	samples          []string
}
