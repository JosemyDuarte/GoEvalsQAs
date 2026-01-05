package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/dute/go_evals/pkg/eval"
)

func main() {
	inputPath := flag.String("input", "data/eval_dataset_generated.jsonl", "Path to generated dataset")
	outputPath := flag.String("output", "data/manual_review.csv", "Path to save CSV template")
	flag.Parse()

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

	writer := csv.NewWriter(outFile)
	writer.Comma = ';' // Consistent with align_judge
	defer writer.Flush()

	writer.Write([]string{"ID", "Reference_Answer", "System_Answer", "Human_Label_Correct (1/0)", "Human_Notes"})

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		var item eval.GeneratedCase
		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			continue
		}

		writer.Write([]string{
			item.ID,
			item.Reference,
			item.SystemAnswer,
			"",
			"",
		})
	}
	fmt.Printf("Done! CSV template saved to %s\n", *outputPath)
}
