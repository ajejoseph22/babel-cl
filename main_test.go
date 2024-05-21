package main

import (
	"testing"
)

func TestReadEvents(t *testing.T) {
	events := readEvents("files/input.json")
	if len(events) == 0 {
		t.Errorf("No events were read from the file")
	}

	expectedEvents := []Event{
		{
			Timestamp: "2018-12-26 18:11:08.509654",
			Duration:  20,
		},
		{
			Timestamp: "2018-12-26 18:15:19.903159",
			Duration:  31,
		},
		{
			Timestamp: "2018-12-26 18:23:19.903159",
			Duration:  54,
		},
	}

	if len(events) != len(expectedEvents) {
		t.Errorf("Expected %d events, got %d", len(expectedEvents), len(events))
	}

	for i, event := range events {
		if event.Timestamp != expectedEvents[i].Timestamp || event.Duration != expectedEvents[i].Duration {
			t.Errorf("At index %d: expected %v, got %v", i, expectedEvents[i], event)
		}
	}
}

func TestCalculateMovingAverages(t *testing.T) {
	events := []Event{
		{
			Timestamp: "2018-12-26 18:11:08.509654",
			Duration:  20,
		},
		{
			Timestamp: "2018-12-26 18:15:19.903159",
			Duration:  31,
		},
		{
			Timestamp: "2018-12-26 18:23:19.903159",
			Duration:  54,
		},
	}
	outputs := calculateMovingAverages(events, 10)

	expectedOutputs := []Output{
		{"2018-12-26 18:11:00", 0},
		{"2018-12-26 18:12:00", 20},
		{"2018-12-26 18:13:00", 20},
		{"2018-12-26 18:14:00", 20},
		{"2018-12-26 18:15:00", 20},
		{"2018-12-26 18:16:00", 25.5},
		{"2018-12-26 18:17:00", 25.5},
		{"2018-12-26 18:18:00", 25.5},
		{"2018-12-26 18:19:00", 25.5},
		{"2018-12-26 18:20:00", 25.5},
		{"2018-12-26 18:21:00", 25.5},
		{"2018-12-26 18:22:00", 31},
		{"2018-12-26 18:23:00", 31},
		{"2018-12-26 18:24:00", 42.5},
	}

	if len(outputs) != len(expectedOutputs) {
		t.Errorf("Expected %d outputs, got %d", len(expectedOutputs), len(outputs))
	}

	for i, output := range outputs {
		if output.Date != expectedOutputs[i].Date || output.AverageDeliveryTime != expectedOutputs[i].AverageDeliveryTime {
			t.Errorf("At index %d: expected %v, got %v", i, expectedOutputs[i], output)
		}
	}

}
