package process

import (
	"bufio"
	"encoding/csv"
	"flowLogParser/pkg/types"
	"flowLogParser/pkg/utils"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// ProcessChunk processes a chunk of lines, updating the tagCounts and portProtocolCounts maps.
func ProcessChunk(lines []string, table types.LookupTable, tagCounts map[string]int, portProtocolCounts map[string]int, mutex *sync.Mutex, protocolMap map[int]string) {
	for _, line := range lines {
		flowLog, err := utils.ParseFlowLog(line, protocolMap)
		if err != nil {
			log.Printf("Failed to parse flow log: %v", err)
			continue
		}

		key := types.LookupKey{DstPort: flowLog.DstPort, Protocol: strings.ToLower(flowLog.Protocol)}

		mutex.Lock()
		if tags, exists := table[key]; exists {
			for _, tag := range tags {
				tagCounts[tag]++
			}
		} else {
			tagCounts["Untagged"]++
		}
		portProtocolKey := fmt.Sprintf("%s,%s", flowLog.DstPort, strings.ToLower(flowLog.Protocol))
		portProtocolCounts[portProtocolKey]++
		mutex.Unlock()
	}
}

// CountTagMatches reads the file with flowLogs and counts the number of matches for each tag and port/protocol combination.
func CountTagMatches(filePath string, table types.LookupTable, protocolMap map[int]string) (map[string]int, map[string]int, error) {
	flowLogFile, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer flowLogFile.Close()

	tagCounts := make(map[string]int)
	portProtocolCounts := make(map[string]int)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	reader := bufio.NewReader(flowLogFile)
	chunkSize := 1024 * 64
	numWorkers := runtime.NumCPU()

	linesChan := make(chan []string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for lines := range linesChan {
				ProcessChunk(lines, table, tagCounts, portProtocolCounts, &mutex, protocolMap)
			}
		}()
	}

	buffer := make([]byte, chunkSize)
	lines := make([]string, 0)
	for {
		n, err := reader.Read(buffer)
		if err != nil && n == 0 {
			break
		}
		chunk := string(buffer[:n])
		lines = append(lines, strings.Split(chunk, "\n")...)

		if len(lines) >= numWorkers {
			linesChan <- lines
			lines = nil
		}
	}

	if len(lines) > 0 {
		linesChan <- lines
	}
	close(linesChan)
	wg.Wait()

	return tagCounts, portProtocolCounts, nil
}

// LoadProtocolMapping loads the protocol number-to-name mapping from a CSV file.
func LoadProtocolMapping(filename string) (map[int]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open protocol mapping file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read protocol mapping file: %w", err)
	}

	protocolMap := make(map[int]string)
	for _, row := range rows[1:] {
		decimal, _ := strconv.Atoi(row[0])
		keyword := row[1]
		protocolMap[decimal] = keyword
	}
	return protocolMap, nil
}
