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
	type TimeEvent struct {
		Time     time.Time
		Duration int
	}

	if len(events) == 0 {
		return []Output{}
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

	var outputs []Output
	currentMinute := timeEvents[0].Time.Truncate(time.Minute)
	endTime := timeEvents[len(timeEvents)-1].Time.Truncate(time.Minute).Add(time.Minute)

	window := []TimeEvent{}
	sum := 0
	translationCount := 0

	for !currentMinute.After(endTime) {
		// Remove events that are outside the window
		for len(window) > 0 && window[0].Time.Before(currentMinute.Add(-time.Duration(windowSize)*time.Minute)) {
			sum -= window[0].Duration
			translationCount--
			window = window[1:]
		}

		// Add events that are in the current minute
		for len(timeEvents) > 0 && !timeEvents[0].Time.After(currentMinute) {
			sum += timeEvents[0].Duration
			translationCount++
			window = append(window, timeEvents[0])
			timeEvents = timeEvents[1:]
		}

		// Compute the average
		average := 0.0
		if translationCount > 0 {
			average = float64(sum) / float64(translationCount)
		}

		// Append the result
		outputs = append(outputs, Output{
			Date:                currentMinute.Format("2006-01-02 15:04:05"),
			AverageDeliveryTime: average,
		})

		// Move on to the next minute
		currentMinute = currentMinute.Add(time.Minute)
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
