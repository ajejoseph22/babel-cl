package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Event struct {
	Timestamp      string `json:"timestamp"`
	TranslationID  string `json:"translation_id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	ClientName     string `json:"client_name"`
	EventName      string `json:"event_name"`
	Duration       int    `json:"duration"`
	NrWords        int    `json:"nr_words"`
}

type Output struct {
	Date                string  `json:"date"`
	AverageDeliveryTime float64 `json:"average_delivery_time"`
}

func main() {
	inputFile := flag.String("input_file", "events.json", "Input JSON file with events")
	outputFile := flag.String("output_file", "output.json", "Output JSON file for results")
	windowSize := flag.Int("window_size", 10, "Window size for the moving average")
	flag.Parse()

	events := readEvents(*inputFile)
	outputs := calculateMovingAverages(events, *windowSize)
	writeOutput(outputs, *outputFile)
}

func readEvents(filename string) []Event {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event Event
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	return events
}

func calculateMovingAverages(events []Event, windowSize int) []Output {
	if len(events) == 0 {
		return []Output{}
	}

	// Parsing timestamps and converting to TimeEvent
	type TimeEvent struct {
		Time     time.Time
		Duration int
	}
	var timeEvents []TimeEvent
	for _, event := range events {
		t, err := time.Parse("2006-01-02 15:04:05.000000", event.Timestamp)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err)
			os.Exit(1)
		}
		timeEvents = append(timeEvents, TimeEvent{Time: t, Duration: event.Duration})
	}

	// Calculating moving averages
	var outputs []Output
	startTime := timeEvents[0].Time.Truncate(time.Minute)
	endTime := timeEvents[len(timeEvents)-1].Time.Add(time.Minute).Truncate(time.Minute)

	// Iterate from startTime to endTime by minute
	for currentTime := startTime; !currentTime.After(endTime); currentTime = currentTime.Add(time.Minute) {
		var sum int
		var count int
		for _, te := range timeEvents {
			lowerBoundWindow := currentTime.Add(-time.Duration(windowSize) * time.Minute)
			if te.Time.After(lowerBoundWindow) && te.Time.Before(currentTime) {
				sum += te.Duration
				count++
			}
		}
		average := 0.0
		if count > 0 {
			average = float64(sum) / float64(count)
		}
		outputs = append(outputs, Output{
			Date:                currentTime.Format("2006-01-02 15:04:05"),
			AverageDeliveryTime: average,
		})
	}

	return outputs
}

func writeOutput(outputs []Output, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, output := range outputs {
		outJson, err := json.Marshal(output)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		writer.WriteString(string(outJson) + "\n")
	}
	writer.Flush()
}
