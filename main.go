package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LogEntry represents a single line from the NDJSON log file
type LogEntry struct {
	Source           string      `json:"source"`
	SubscriptionName string      `json:"subscription-name"`
	Timestamp        int64       `json:"timestamp"`
	Time             string      `json:"time"`
	Prefix           string      `json:"prefix,omitempty"`
	Updates          []UpdateMsg `json:"updates,omitempty"`
}

// UpdateMsg represents an update with a path and values
type UpdateMsg struct {
	Path   string                 `json:"Path"`
	Values map[string]interface{} `json:"values"`
}

func main() {
	inputFile := flag.String("file", "", "Input file containing gNMI subscribe messages in NDJSON format")
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -file flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Open the input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read file line by line (NDJSON format)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Parse the JSON line
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Show first 200 characters of the problematic line for debugging
			preview := line
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			fmt.Fprintf(os.Stderr, "Error parsing JSON on line %d: %v\n", lineNum, err)
			fmt.Fprintf(os.Stderr, "Line preview: %s\n", preview)
			continue
		}

		// Process the entry
		processEntry(entry)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
}

func processEntry(entry LogEntry) {
	// Format timestamp (nanoseconds since Unix epoch)
	timestamp := time.Unix(0, entry.Timestamp).Format(time.RFC3339Nano)

	// Process each update
	for _, update := range entry.Updates {
		// Build the full path from prefix and update path
		fullPath := ""
		if entry.Prefix != "" {
			fullPath = entry.Prefix
		}

		if update.Path != "" {
			if fullPath != "" {
				fullPath = fullPath + "/" + update.Path
			} else {
				fullPath = update.Path
			}
		}

		// Process each value in the update
		for _, value := range update.Values {
			// Use fullPath directly - it contains the complete path from update.Path
			// which includes selectors like [name=Ethernet1/1]
			// The keys in values often have incomplete paths missing these selectors
			valuePath := fullPath
			if valuePath == "" {
				// Fallback: if no path from prefix/update.Path, this shouldn't happen
				// but handle it gracefully
				continue
			}

			// Format the value
			valueStr := formatValue(value)

			// Output with timestamp prefix
			fmt.Printf("[%s] %s = %s\n", timestamp, valuePath, valueStr)
		}
	}
}

// formatValue converts a value to a string representation
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		// Use strconv.FormatFloat with 'f' format to avoid scientific notation
		// prec -1 means use the smallest number of digits necessary
		str := strconv.FormatFloat(v, 'f', -1, 64)
		return str
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "null"
	case map[string]interface{}:
		// Check if this is a leaflist (has "element" key with array)
		if elements, ok := v["element"]; ok {
			if elemArray, ok := elements.([]interface{}); ok {
				return formatLeaflist(elemArray)
			}
		}
		// For other complex types, marshal to JSON
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	default:
		// For complex types, marshal to JSON
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(jsonBytes)
	}
}

// formatLeaflist formats a gNMI leaflist value as a comma-separated list
func formatLeaflist(elements []interface{}) string {
	var values []string
	for _, elem := range elements {
		// Each element should have a "Value" key with the typed value
		if elemMap, ok := elem.(map[string]interface{}); ok {
			if valueMap, ok := elemMap["Value"].(map[string]interface{}); ok {
				// Extract the actual value from the typed value
				var elemValue interface{}
				for _, v := range valueMap {
					elemValue = v
					break
				}
				values = append(values, formatValue(elemValue))
			}
		}
	}
	return "[" + strings.Join(values, ", ") + "]"
}
