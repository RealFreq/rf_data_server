package main

import (
	"fmt"
	"testing"
	"time"
)

const (
	test_csv_data = "2015-03-09, 18:06:18, 24000000, 26797996, 87437.38, 256, -22.74, -22.29, -20.07"
)

func TestParseCsv(t *testing.T) {
	expected := []string{"2015-03-09", "18:06:18", "24000000", "26797996", "87437.38", "256", "-22.74", "-22.29", "-20.07"}

	result := parseCsv(test_csv_data)

	if len(result) != len(expected) {
		t.Errorf("\nexpected %s\nreceived %s\n", expected, result)
		return
	}

	for i, val := range result {
		if expected[i] != val {
			t.Errorf("\nexpected %s\nreceived %s\n", expected, result)
			return
		}
	}
}

func TestGetTimestampFromRecord(t *testing.T) {
	data := []string{"2015-03-09", "18:06:18"}
	tz, _ := time.Now().Local().Zone()
	expected, _ := time.Parse(timeLayout,
		fmt.Sprintf("%s %s (%s)", data[0], data[1], tz))
	result := getTimestampFromRecord(data)

	if expected != result {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.Format(timeLayout), result.Format(timeLayout))
	}
}

func TestBuildRfRecord(t *testing.T) {
	data := parseCsv(test_csv_data)
	lat_long := getGpsCoordinates()

	expected := RfRecord{
		latitude:         lat_long[0],
		longitude:        lat_long[1],
		recorded_at:      getTimestampFromRecord(data),
		hz_start:         data[2],
		hz_end:           data[3],
		hz_step:          data[4],
		averaged_samples: data[5],
		samples:          data[6 : len(data)-1],
	}

	result := buildRfRecord(test_csv_data)

	if result.latitude != expected.latitude {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.latitude, result.latitude)
		return
	}

	if result.longitude != expected.longitude {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.longitude, result.longitude)
		return
	}
	if result.recorded_at != expected.recorded_at {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.recorded_at, result.recorded_at)
		return
	}
	if result.hz_start != expected.hz_start {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.hz_start, result.hz_start)
		return
	}
	if result.hz_end != expected.hz_end {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.hz_end, result.hz_end)
		return
	}
	if result.hz_step != expected.hz_step {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.hz_step, result.hz_step)
		return
	}
	if result.averaged_samples != expected.averaged_samples {
		t.Errorf("\nexpected %s\nreceived %s\n", expected.averaged_samples, result.averaged_samples)
		return
	}
	for i, val := range result.samples {
		if val != expected.samples[i] {
			t.Errorf("\nexpected %s\nreceived %s\n", expected.samples[i], val)
			return
		}
	}
}
