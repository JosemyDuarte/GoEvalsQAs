package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dute/go_evals/pkg/eval"
)

func main() {
	labelsPath := flag.String("labels", "data/manual_review_done.csv", "Path to human labels CSV")
	command := flag.String("cmd", "gemini", "Judge CLI command")
	argsStr := flag.String("args", "", "Comma-separated arguments for the judge CLI")
	flag.Parse()

	args := strings.Split(*argsStr, ",")
	provider := &eval.CLIProvider{
		Command: *command,
		Args:    args,
	}

	file, err := os.Open(*labelsPath)
	if err != nil {
		fmt.Printf("Error opening labels file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	matches := 0
	total := 0
	failuresCaught := 0
	actualFailures := 0

	fmt.Printf("Running Judge Alignment for %d records\n", len(records)-1)

	for i, row := range records {
		if i == 0 {
			continue
		}

		if len(row) < 4 {
			continue
		}

		id := row[0]
		ref := row[1]
		sysAns := row[2]
		humanLabel := strings.TrimSpace(row[3])

		if humanLabel == "" {
			continue
		}

		prompt := fmt.Sprintf(`
You are a strict factual judge.
Reference Answer: "%s"
Student Answer: "%s"

Does the Student Answer contain the core fact from the Reference Answer? 
Ignore minor phrasing differences. 
If the reference is a date or name, it MUST be present.
Reply ONLY with "1" for Yes or "0" for No. Don't add explanations.
`, ref, sysAns)

		judgeOut, err := provider.Generate(prompt)
		if err != nil {
			fmt.Printf("Error judging %s: %v\n", id, err)
			continue
		}

		judgeVerdict := "0"
		if strings.Contains(judgeOut, "1") {
			judgeVerdict = "1"
		}

		isMatch := (judgeVerdict == humanLabel)
		if isMatch {
			matches++
		} else {
			fmt.Printf("[MISMATCH] ID: %s | Human: %s | Judge: %s\n", id, humanLabel, judgeVerdict)
		}

		if humanLabel == "0" {
			actualFailures++
			if judgeVerdict == "0" {
				failuresCaught++
			}
		}

		total++
		fmt.Printf("\rProcessed: %d | Agreement: %.0f%%", total, (float64(matches)/float64(total))*100)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n\n--- Alignment Report ---")
	fmt.Printf("Total Labeled: %d\n", total)
	fmt.Printf("Agreement Rate: %.2f%%\n", (float64(matches)/float64(total))*100)

	if actualFailures > 0 {
		fmt.Printf("Defect Recall: %.2f%% (%d/%d bad answers caught)\n",
			(float64(failuresCaught)/float64(actualFailures))*100, failuresCaught, actualFailures)
	}
}
