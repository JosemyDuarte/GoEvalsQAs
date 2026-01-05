package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dute/go_evals/pkg/eval"
)

func main() {
	inputPath := flag.String("input", "data/golden_set.jsonl", "Path to golden set")
	outputPath := flag.String("output", "data/eval_dataset_generated.jsonl", "Path to save generated results")
	command := flag.String("cmd", "ollama", "CLI command to run (e.g., gemini, ollama)")
	argsStr := flag.String("args", "run,llama3.3", "Comma-separated arguments for the CLI command")
	flag.Parse()

	args := strings.Split(*argsStr, ",")
	provider := &eval.CLIProvider{
		Command: *command,
		Args:    args,
	}

	inFile, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inFile.Close()

	outFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	scanner := bufio.NewScanner(inFile)

	fmt.Printf("Starting generation using command: %s %v\n", provider.Command, provider.Args)
	count := 0

	for scanner.Scan() {
		var input eval.EvalCase
		if err := json.Unmarshal(scanner.Bytes(), &input); err != nil {
			continue
		}

		// Retry logic for CLI flakes
		var answer string
		var genErr error
		for i := 0; i < 3; i++ {
			answer, genErr = provider.Generate(input.Prompt)
			if genErr == nil && answer != "" {
				break
			}
			time.Sleep(2 * time.Second)
		}

		if genErr != nil {
			fmt.Printf("\n[Fail] ID %s: %v", input.ID, genErr)
			answer = "ERROR_GENERATING"
		}

		output := eval.GeneratedCase{
			EvalCase:     input,
			SystemAnswer: answer,
		}

		jsonLine, _ := json.Marshal(output)
		writer.Write(jsonLine)
		writer.WriteString("\n")
		writer.Flush()

		count++
		fmt.Printf("\rProcessed: %d", count)
	}
	fmt.Println("\nDone.")
}
