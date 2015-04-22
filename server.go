package main

const (
	HOST = "0.0.0.0"
	PORT = 10000
)

func main() {
	parser_queue := make(chan string)
	logger_queue := make(chan RfRecord)

	go Parse(parser_queue, logger_queue)
	go Logger(logger_queue)

	RfDataServer(HOST, PORT, parser_queue)
}
