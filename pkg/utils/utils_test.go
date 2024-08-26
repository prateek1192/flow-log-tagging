package utils

import (
	"flowLogParser/pkg/types"
	"os"
	"testing"
)

func TestParseFlowLog(t *testing.T) {
	// Define the test input line and protocol map
	line := "2 123456789012 eni-0a1b2c3d 10.0.1.201 198.51.100.2 443 49153 6 25 20000 1620140761 1620140821 ACCEPT OK"
	protocolMap := map[int]string{
		6: "TCP",
	}

	// Parse the flow log
	flowLog, err := ParseFlowLog(line, protocolMap)
	if err != nil {
		t.Fatalf("Failed to parse flow log: %v", err)
	}

	// Assert each field in the FlowLog struct
	if flowLog.Version != 2 {
		t.Errorf("Expected version 2, got %d", flowLog.Version)
	}
	if flowLog.AccountId != "123456789012" {
		t.Errorf("Expected accountId 123456789012, got %s", flowLog.AccountId)
	}
	if flowLog.InterfaceId != "eni-0a1b2c3d" {
		t.Errorf("Expected interfaceId eni-0a1b2c3d, got %s", flowLog.InterfaceId)
	}
	if flowLog.SrcAddr != "10.0.1.201" {
		t.Errorf("Expected srcAddr 10.0.1.201, got %s", flowLog.SrcAddr)
	}
	if flowLog.DstAddr != "198.51.100.2" {
		t.Errorf("Expected dstAddr 198.51.100.2, got %s", flowLog.DstAddr)
	}
	if flowLog.SrcPort != "443" {
		t.Errorf("Expected srcPort 443, got %s", flowLog.SrcPort)
	}
	if flowLog.DstPort != "49153" {
		t.Errorf("Expected dstPort 49153, got %s", flowLog.DstPort)
	}
	if flowLog.Protocol != "TCP" {
		t.Errorf("Expected protocol TCP, got %s", flowLog.Protocol)
	}
	if flowLog.Packets != 25 {
		t.Errorf("Expected packets 25, got %d", flowLog.Packets)
	}
	if flowLog.Bytes != 20000 {
		t.Errorf("Expected bytes 20000, got %d", flowLog.Bytes)
	}
	if flowLog.Start != 1620140761 {
		t.Errorf("Expected start time 1620140761, got %d", flowLog.Start)
	}
	if flowLog.End != 1620140821 {
		t.Errorf("Expected end time 1620140821, got %d", flowLog.End)
	}
	if flowLog.Action != "ACCEPT" {
		t.Errorf("Expected action ACCEPT, got %s", flowLog.Action)
	}
	if flowLog.LogStatus != "OK" {
		t.Errorf("Expected logStatus OK, got %s", flowLog.LogStatus)
	}
}

func TestLoadLookupTable(t *testing.T) {
	csvContent := "dstport,protocol,tag\n443,6,sv_P1\n80,6,sv_P2"
	tempFile, err := os.CreateTemp("", "lookup-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	lookupTable, err := CreateLookupTable(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load lookup table: %v", err)
	}

	if len(lookupTable) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(lookupTable))
	}
}

func TestLoadLookupTableMultipleSvs(t *testing.T) {
	csvContent := "dstport,protocol,tag\n443,6,sv_P1\n443,6,sv_P2\n80,6,sv_P3"
	tempFile, err := os.CreateTemp("", "lookup-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	lookupTable, err := CreateLookupTable(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load lookup table: %v", err)
	}

	if len(lookupTable) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(lookupTable))
	}

	key := types.LookupKey{DstPort: "443", Protocol: "6"}

	tags, exists := lookupTable[key]
	if !exists {
		t.Fatalf("Expected key %v to exist in lookupTable", key)
	}

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags for key %v, got %d", key, len(tags))
	}

	expectedTags := []string{"sv_P1", "sv_P2"}
	for _, tag := range expectedTags {
		if !contains(tags, tag) {
			t.Errorf("Expected tag %s to be present for key %v", tag, key)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
