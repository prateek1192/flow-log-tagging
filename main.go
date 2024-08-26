package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	lookupFilePath := flag.String("lookupfile", "lookup_table.csv", "Path to the lookup table file")
	flowLogFilePath := flag.String("flowlogfile", "flowLog.txt", "Path to the flow log file")
	outputFilePath := flag.String("outputfile", "tag_counts.csv", "Path to the output file")

	// Parse command-line flags
	flag.Parse()

	protocolMap, err := LoadProtocolMapping("protocol-numbers-1.csv")
	if err != nil {
		log.Fatalf("Failed to load protocol mapping: %v", err)
	}

	lookupTable, err := CreateLookupTable(*lookupFilePath)
	if err != nil {
		log.Fatalf("Error loading lookup table: %v\n", err)
	}

	tagCounts, portProtocolCounts, err := CountTagMatches(*flowLogFilePath, lookupTable, protocolMap)
	if err != nil {
		log.Fatalf("Error processing flow logs: %v\n", err)
	}

	err = WriteTagCounts(tagCounts, portProtocolCounts, *outputFilePath)
	if err != nil {
		log.Fatalf("Error writing tag counts: %v\n", err)
	}
	for k, v := range tagCounts {
		fmt.Printf("Key %v: value %v\n", k, v)
	}
	for k, v := range portProtocolCounts {
		fmt.Printf("Key %v: value %v\n", k, v)
	}

	log.Printf("Tag counts written to %s\n", outputFilePath)
}
