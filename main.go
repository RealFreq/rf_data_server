package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	matchPattern = "^[0-9]{4}-[0-9]{2}-[0-9]{2}, [0-9]{2}:[0-9]{2}:[0-9]{2}, *"
	timeLayout   = "2006-01-02 15:04:05 (MST)"
	apiUrl       = "http://sheltered-anchorage-8761.herokuapp.com"
)

type rfRecord struct {
	latitude         string
	longitude        string
	recorded_at      time.Time
	hz_start         string
	hz_end           string
	hz_step          string
	averaged_samples string
	samples          []string
}

func main() {
	// TODO Add support for reading from file
	// TODO Add flags for server options
	scanInput()
}

func scanInput() {
	lines := make(chan string)
	raw_record_queue := make(chan string)
	record_queue := make(chan []string)
	tx_queue := make(chan rfRecord)

	scanner := bufio.NewScanner(os.Stdin)

	go parseLines(lines, raw_record_queue)
	go scanCsv(raw_record_queue, record_queue)
	go parseRecord(record_queue, tx_queue)
	go uploadRecord(tx_queue)

	for {
		scanner.Scan()
		lines <- scanner.Text()
		if err := scanner.Err(); err != nil {
			// TODO handle error gracefully
			log.Fatal(err)
		}
	}
}

func parseLines(lines <-chan string, raw_record_queue chan<- string) {
	for line := range lines {
		matched, err := regexp.MatchString(matchPattern, line)
		if matched {
			raw_record_queue <- line
		}
		if err != nil {
			// TODO handle error gracefully
			log.Fatal(err)
		}
	}
}

func scanCsv(raw_record_queue <-chan string, record_queue chan<- []string) {
	for raw_record := range raw_record_queue {
		reader := csv.NewReader(strings.NewReader(raw_record))
		record, err := reader.Read()
		if err != nil {
			// TODO handle error gracefully
			log.Fatal(err)
		}
		record_queue <- record
	}
}

func parseRecord(record_queue <-chan []string, tx_queue chan<- rfRecord) {
	for record := range record_queue {
		tx_queue <- buildRfRecord(record)
	}
}

func uploadRecord(tx_queue <-chan rfRecord) {
	for record := range tx_queue {
		// TODO handle error

		params := url.Values{}
		params.Set("rf_record[latitude]", record.latitude)
		params.Set("rf_record[longitude]", record.longitude)
		params.Set("rf_record[recorded_at]", record.recorded_at.String())
		params.Set("rf_record[hz_start]", record.hz_start)
		params.Set("rf_record[hz_end]", record.hz_end)
		params.Set("rf_record[hz_step]", record.hz_step)
		params.Set("rf_record[num_samples_per_reading]", record.averaged_samples)
		samples, _ := json.Marshal(record.samples)
		params.Set("rf_record[samples]", string(samples))

		resp, _ := http.PostForm(apiUrl+"/rf_records.json", params)

		defer resp.Body.Close()

		log.Printf("response Status:", resp.Status)
		log.Printf("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("response Body:", string(body))
	}
}

func buildRfRecord(record []string) rfRecord {
	return rfRecord{
		latitude:         "32.7513718",          // Need to get from GPS feed
		longitude:        "-117.14624170000002", // Need to get from GPS feed
		recorded_at:      getTimestampFromRecord(record),
		hz_start:         record[2],
		hz_end:           record[3],
		hz_step:          record[4],
		averaged_samples: record[5],
		samples:          record[6 : len(record)-1],
	}
}

func getTimestampFromRecord(record []string) time.Time {
	record_date := record[0]
	record_time := record[1]
	record_tz, _ := time.Now().Local().Zone()

	composed_timestamp := fmt.Sprintf("%s%s (%s)", record_date, record_time, record_tz)

	recorded_at, err := time.Parse(timeLayout, composed_timestamp)
	if err != nil {
		// TODO handle error gracefully
		log.Fatal(err)
	}

	return recorded_at
}
