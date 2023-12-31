package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"pasza.org/sr-challenge/evaluator"
	"pasza.org/sr-challenge/parser"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func validateCommandLine() (inputPath, outputPath string) {
	argv := os.Args
	argc := len(argv)

	switch argc {
	case 2:
		inputPath = argv[1]
	case 3:
		inputPath = argv[1]
		outputPath = argv[2]
		if fileExists(outputPath) {
			log.Fatalf("Output file already exists")
			os.Exit(1)
		}
	default:
		fmt.Printf("Usage: %s <input_file> [output_file]\n", argv[0])
		os.Exit(1)
	}
	if !fileExists(inputPath) {
		log.Fatal("Input path does not exist")
		os.Exit(1)
	}
	return
}

// Writes output to file or to stdout
func writeOutput(outputPath string, result [][]evaluator.CalculatedValue) {
	var f *os.File
	var err error
	if outputPath != "" {
		f, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open output file")
			os.Exit(1)
		}
		defer f.Close()
	} else {
		f = os.Stdout
	}
	writer := bufio.NewWriter(f)
	// do the writing here
	for _, row := range result {
		first := true
		for _, cell := range row {
			if first {
				first = false
			} else {
				writer.WriteString(" |")
			}
			writer.WriteString(cell.String())
		}
		writer.WriteString("\n")
	}
	// done writing
	defer writer.Flush()
}

func main() {
	inputPath, outputPath := validateCommandLine()
	input, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	csv, ok, err := parser.ParseCSV(string(input))
	if err != nil {
		log.Fatalf("failed to parse: %v\n", err)
		os.Exit(1)
	}

	if !ok {
		log.Fatal("expected CSV not matched\n")
		os.Exit(1)
	}
	// evaluate
	result := evaluator.Evaluate(csv)
	// format output
	writeOutput(outputPath, result)
}
