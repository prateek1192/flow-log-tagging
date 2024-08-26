package utils

import (
	"encoding/csv"
	"flowLogParser/pkg/types"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"strconv"
	"strings"
)

// ParseFlowLog parses a line into a types.FlowLog struct
func ParseFlowLog(line string, protocolMap map[int]string) (types.FlowLog, error) {
	fields := strings.Fields(line)
	if len(fields) < 14 {
		return types.FlowLog{}, fmt.Errorf("invalid line format: %s", line)
	}

	version, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return types.FlowLog{}, fmt.Errorf("error parsing version: %w", err)
	}
	protocolNumber, err := strconv.ParseInt(fields[7], 10, 32)
	if err != nil {
		return types.FlowLog{}, fmt.Errorf("failed to parse protocol number: %w", err)
	}
	protocolName, exists := protocolMap[int(protocolNumber)]
	if !exists {
		return types.FlowLog{}, fmt.Errorf("unknown protocol number: %d", protocolNumber)
	}
	packets, err := strconv.ParseInt(fields[8], 10, 64)
	if err != nil {
		return types.FlowLog{}, fmt.Errorf("error parsing packets: %w", err)
	}
	bytes, err := strconv.ParseInt(fields[9], 10, 64)
	if err != nil {
		return types.FlowLog{}, fmt.Errorf("error parsing bytes: %w", err)
	}
	start, err := strconv.ParseInt(fields[10], 10, 64)
	if err != nil {
		return types.FlowLog{}, fmt.Errorf("error parsing start time: %w", err)
	}
	end, err := strconv.ParseInt(fields[11], 10, 64)
	if err != nil {
		return types.FlowLog{}, fmt.Errorf("error parsing end time: %w", err)
	}

	return types.FlowLog{
		Version:     int32(version),
		AccountId:   fields[1],
		InterfaceId: fields[2],
		SrcAddr:     fields[3],
		DstAddr:     fields[4],
		SrcPort:     fields[5],
		DstPort:     fields[6],
		Protocol:    protocolName,
		Packets:     packets,
		Bytes:       bytes,
		Start:       start,
		End:         end,
		Action:      fields[12],
		LogStatus:   fields[13],
	}, nil
}

// CreateLookupTable opens the lookup table file and reads it into memory
func CreateLookupTable(filePath string) (types.LookupTable, error) {

	lookupFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer lookupFile.Close()

	lookupEntries := make([]*types.LookupEntry, 0)

	err = gocsv.UnmarshalFile(lookupFile, &lookupEntries)
	if err != nil {
		return nil, err
	}

	lookupTable := make(types.LookupTable)

	for i := 0; i < len(lookupEntries); i++ {
		key := types.LookupKey{DstPort: lookupEntries[i].DstPort, Protocol: strings.ToLower(lookupEntries[i].Protocol)}
		lookupTable[key] = append(lookupTable[key], lookupEntries[i].Tag)

	}
	return lookupTable, nil
}

func WriteTagCounts(tagCounts map[string]int, portProtocolCounts map[string]int, outputFilePath string) error {
	countFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating file for writing counts: %v", err)
	}
	defer countFile.Close()

	writer := csv.NewWriter(countFile)
	defer writer.Flush()

	err = writer.Write([]string{"Tag", "Count"})
	if err != nil {
		return fmt.Errorf("error writing to CSV file: %w", err)
	}

	for tag, count := range tagCounts {
		err = writer.Write([]string{tag, fmt.Sprintf("%d", count)})
		if err != nil {
			return fmt.Errorf("error writing tag count to CSV: %w", err)
		}
	}

	err = writer.Write([]string{"Port", "Protocol", "Count"})
	if err != nil {
		return fmt.Errorf("error writing to CSV file: %w", err)
	}

	for portProtocol, count := range portProtocolCounts {

		err = writer.Write([]string{portProtocol, fmt.Sprintf("%d", count)})
		if err != nil {
			return fmt.Errorf("error writing port/protocol count to CSV: %w", err)
		}
	}

	return nil
}
