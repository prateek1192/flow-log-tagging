package main

import (
	"sync"
	"testing"
)

func TestProcessChunk(t *testing.T) {

	lookupTable := LookupTable{
		LookupKey{DstPort: "443", Protocol: "tcp"}: {"sv_P1"},
	}

	protocolMap := map[int]string{
		6:  "tcp",
		17: "udp",
	}

	tagCounts := make(map[string]int)
	portProtocolCounts := make(map[string]int)
	var mutex sync.Mutex

	lines := []string{
		"2 123456789012 eni-1a2b3c4d 10.0.1.102 172.217.7.228 49152 443 6 8 4000 1620140661 1620140721 ACCEPT OK",
	}
	ProcessChunk(lines, lookupTable, tagCounts, portProtocolCounts, &mutex, protocolMap)
	expectedTagCounts := map[string]int{
		"sv_P1": 1,
	}

	expectedPortProtocolCounts := map[string]int{
		"443,tcp": 1,
	}

	for tag, expectedCount := range expectedTagCounts {
		if count, ok := tagCounts[tag]; !ok || count != expectedCount {
			t.Errorf("Expected tag count for %s to be %d, got %d", tag, expectedCount, count)
		}
	}

	for portProtocol, expectedCount := range expectedPortProtocolCounts {
		if count, ok := portProtocolCounts[portProtocol]; !ok || count != expectedCount {
			t.Errorf("Expected port/protocol count for %s to be %d, got %d", portProtocol, expectedCount, count)
		}
	}

	// Verify that there are no unexpected tags
	for tag := range tagCounts {
		if _, ok := expectedTagCounts[tag]; !ok {
			t.Errorf("Unexpected tag %s found in counts", tag)
		}
	}

	for portProtocol := range portProtocolCounts {
		if _, ok := expectedPortProtocolCounts[portProtocol]; !ok {
			t.Errorf("Unexpected port/protocol combination %s found in counts", portProtocol)
		}
	}
}

func TestProcessChunkMultiple(t *testing.T) {

	lookupTable := LookupTable{
		LookupKey{DstPort: "443", Protocol: "tcp"}: {"sv_P1", "sv_P2"},
	}

	protocolMap := map[int]string{
		6:  "tcp",
		17: "udp",
	}

	tagCounts := make(map[string]int)
	portProtocolCounts := make(map[string]int)
	var mutex sync.Mutex

	lines := []string{
		"2 123456789012 eni-1a2b3c4d 10.0.1.102 172.217.7.228 49152 443 6 8 4000 1620140661 1620140721 ACCEPT OK",
		"2 123456789012 eni-1a2b3c4d 10.0.1.102 172.217.7.228 49152 423 17 8 4000 1620140661 1620140721 ACCEPT OK",
	}
	ProcessChunk(lines, lookupTable, tagCounts, portProtocolCounts, &mutex, protocolMap)
	expectedTagCounts := map[string]int{
		"sv_P1":    1,
		"sv_P2":    1,
		"Untagged": 1,
	}

	expectedPortProtocolCounts := map[string]int{
		"443,tcp": 1,
		"423,udp": 1,
	}

	for tag, expectedCount := range expectedTagCounts {
		if count, ok := tagCounts[tag]; !ok || count != expectedCount {
			t.Errorf("Expected tag count for %s to be %d, got %d", tag, expectedCount, count)
		}
	}

	for portProtocol, expectedCount := range expectedPortProtocolCounts {
		if count, ok := portProtocolCounts[portProtocol]; !ok || count != expectedCount {
			t.Errorf("Expected port/protocol count for %s to be %d, got %d", portProtocol, expectedCount, count)
		}
	}

	// Verify that there are no unexpected tags
	for tag := range tagCounts {
		if _, ok := expectedTagCounts[tag]; !ok {
			t.Errorf("Unexpected tag %s found in counts", tag)
		}
	}

	for portProtocol := range portProtocolCounts {
		if _, ok := expectedPortProtocolCounts[portProtocol]; !ok {
			t.Errorf("Unexpected port/protocol combination %s found in counts", portProtocol)
		}
	}
}
