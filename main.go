package main

func main() {
	raw_data := make(chan string)
	rf_data := make(chan RfRecord)

	Parse(raw_data, rf_data)

	// Upload data to server

	// Infinite loop that reads in the data
	ReadData(raw_data)
}
