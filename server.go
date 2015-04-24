package main

func main() {
	parser_queue := make(chan string)
	logger_queue := make(chan RfRecord)

	go Parse(parser_queue, logger_queue)
	go Logger(logger_queue)

	host, port := ServerConfig()
	RfDataServer(host, uint32(port), parser_queue)
}
