package main

import (
	"bufio"
	"log"
	"os"
)

//
// Read data from STDIN and send each line to the output channel
//
// @param output chan<- Output channel to use for sending each line read
//											in from STDIN
//
func ReadData(output chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		// TODO Should this be as an else-block to the error check?
		output <- scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println("Error reading from STDIN: %s", err)
		}
	}
}
