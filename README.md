# Unbabel CLI

A simple command-line application that parses a stream of events and produces a file with the moving average of the translation delivery time for the last N minutes.

## Building and Running the Application

### Prerequisites
- Go 1.15 or higher

### Build
1. Clone the repository:
   ```sh
   git clone git@github.com:ajejoseph22/babel-cl.git
   ```
2. Build the application:
   ```sh
   cd babel-cl
   go build -o babel main.go
   ```
3. Run the application:
   ```sh
   ./babel --input_file ./files/input.json --output_file <file path> --window_size <window size>
   ```
   The output file will be generated in the specified path
4. Run the tests:
   ```sh
   go test ./...
   ```