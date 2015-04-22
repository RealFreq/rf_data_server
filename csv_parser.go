package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RawRfRecord struct {
	recorded_at      time.Time
	hz_start         string
	hz_end           string
	hz_step          string
	averaged_samples string
	samples          []string
}

const (
	matchPattern = "^[0-9]{4}-[0-9]{2}-[0-9]{2}, [0-9]{2}:[0-9]{2}:[0-9]{2}, *"
	timeLayout   = "2006-01-02 15:04:05 (MST)"
)

func Parse(data <-chan string, rf_data chan<- RfRecord) {
	for line := range data {
		matched, err := regexp.MatchString(matchPattern, line)
		if err != nil {
			log.Printf("Error matching data: %s\n", err)
			continue
		}

		if false == matched {
			continue
		}

		// TODO Validate the line before trying to parse it
		records, err := buildRfRecords(line)

		go func() {
			for _, record := range records {
				rf_data <- record
			}
		}()
	}
}

func buildRfRecords(line string) ([]RfRecord, error) {
	parsed_csv := parseCsv(line)

	if len(parsed_csv) == 0 {
		return nil, errors.New("Empty raw record")
	}

	raw_record := RawRfRecord{
		recorded_at:      getTimestampFromRecord(parsed_csv),
		hz_start:         parsed_csv[2],
		hz_end:           parsed_csv[3],
		hz_step:          parsed_csv[4],
		averaged_samples: parsed_csv[5],
		samples:          parsed_csv[6 : len(parsed_csv)-1],
	}

	result, err := processRawRfRecord(raw_record)
	return result, err
}

func processRawRfRecord(rawRecord RawRfRecord) ([]RfRecord, error) {
	// TODO This should be handled by a validator in Parse
	step, err := stringToFloat64(rawRecord.hz_step)
	if err != nil {
		msg := fmt.Sprintf("Error parsing step for: %+v - %s\n", rawRecord)
		log.Printf(msg)
		return nil, errors.New(msg)
	}

	start, err := stringToFloat64(rawRecord.hz_start)
	if err != nil {
		msg := fmt.Sprintf("Error parsing step for: %+v - %s\n", rawRecord)
		log.Printf(msg)
		return nil, errors.New(msg)
	}

	records := make([]RfRecord, len(rawRecord.samples))
	for i, raw_power := range rawRecord.samples {
		frequency := start + (float64(i) * step)
		power, err := stringToFloat64(raw_power)

		if err != nil {
			log.Printf("Error parsing power value(%s): %s\n", raw_power, err)
			continue
		}

		records[i] = RfRecord{
			recorded_at: rawRecord.recorded_at,
			frequency:   frequency,
			power:       power,
		}
	}

	return records, nil
}

func stringToFloat64(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
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
