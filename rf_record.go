package main

import (
	"time"
)

type RfRecord struct {
	recorded_at time.Time
	frequency   float64
	power       float64
}
