package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	matchPattern = "^[0-9]{4}-[0-9]{2}-[0-9]{2}, [0-9]{2}:[0-9]{2}:[0-9]{2}, *"
	timeLayout   = "2006-01-02 15:04:05 (MST)"
)

func Parse(csv_data <-chan string, rf_data chan<- RfRecord) {
	var record RfRecord

	for line := range csv_data {
		matched, err := regexp.MatchString(matchPattern, line)
		if matched {
			record = buildRfRecord(line)
		}
		if err != nil {
			log.Println("Error matching data: %s\n", err)
		} else {
			rf_data <- record
		}
	}
}

// TODO add err return for nil
func buildRfRecord(csv_line string) RfRecord {
	parsed_csv := parseCsv(csv_line)

	if len(parsed_csv) == 0 {
		return RfRecord{}
	}

	lat_long := getGpsCoordinates()

	return RfRecord{
		latitude:         lat_long[0],
		longitude:        lat_long[1], // TODO Need to get from GPS feed
		recorded_at:      getTimestampFromRecord(parsed_csv),
		hz_start:         parsed_csv[2],
		hz_end:           parsed_csv[3],
		hz_step:          parsed_csv[4],
		averaged_samples: parsed_csv[5],
		samples:          parsed_csv[6 : len(parsed_csv)-1],
	}
}

// TODO add err return for empty slice
func parseCsv(data string) []string {
	reader := csv.NewReader(strings.NewReader(data))
	record, err := reader.Read()
	if err != nil {
		log.Println("Error parsing CSV: %s\n", err)
		return make([]string, 0)
	}

	for i, val := range record {
		record[i] = strings.TrimSpace(val)
	}

	return record
}

func getTimestampFromRecord(record []string) time.Time {
	record_date := record[0]
	record_time := record[1]
	record_tz, _ := time.Now().Local().Zone()

	composed_timestamp := fmt.Sprintf("%s %s (%s)", record_date, record_time, record_tz)

	recorded_at, err := time.Parse(timeLayout, composed_timestamp)
	if err != nil {
		// TODO handle error gracefully
		log.Fatal(err)
	}

	return recorded_at
}

func getGpsCoordinates() []string {
	coords := make([]string, 2)
	coords[0] = "32.7513718"
	coords[1] = "-117.14624170000002"

	return coords
}
