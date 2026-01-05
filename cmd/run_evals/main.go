package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dute/go_evals/pkg/eval"
)

func main() {
	inputPath := flag.String("input", "data/eval_dataset_generated.jsonl", "Path to generated dataset")
	workerCount := flag.Int("workers", 5, "Number of parallel judge workers")
	command := flag.String("cmd", "gemini", "Judge CLI command")
	argsStr := flag.String("args", "", "Comma-separated arguments for the judge CLI")
	flag.Parse()

	args := []string{}
	if *argsStr != "" {
		args = strings.Split(*argsStr, ",")
	}

	provider := &eval.CLIProvider{
		Command: *command,
		Args:    args,
	}

	start := time.Now()

	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var jobsList []eval.GeneratedCase
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var job eval.GeneratedCase
		if err := json.Unmarshal(scanner.Bytes(), &job); err != nil {
			continue
		}
		jobsList = append(jobsList, job)
	}
	fmt.Printf("loaded %d cases. Starting %d workers...\n", len(jobsList), *workerCount)

	jobs := make(chan eval.GeneratedCase, len(jobsList))
	results := make(chan eval.EvalResult, len(jobsList))
	var wg sync.WaitGroup

	for w := 1; w <= *workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				prompt := fmt.Sprintf(`
You are a strict factual judge.
Reference Answer: "%s"
Student Answer: "%s"

Does the Student Answer contain the core fact from the Reference Answer? 
Reply ONLY with "1" for Yes or "0" for No.
`, job.Reference, job.SystemAnswer)

				var verdict string
				var err error
				for i := 0; i < 3; i++ {
					verdict, err = provider.Generate(prompt)
					if err == nil {
						break
					}
					time.Sleep(1 * time.Second)
				}

				passed := strings.Contains(verdict, "1")
				results <- eval.EvalResult{
					ID:     job.ID,
					Passed: passed,
				}
				fmt.Printf(".")
			}
		}()
	}

	for _, j := range jobsList {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	passCount := 0
	totalCount := 0
	failLog, _ := os.Create("data/failures.log")
	defer failLog.Close()

	for res := range results {
		totalCount++
		if res.Passed {
			passCount++
		} else {
			failLog.WriteString(fmt.Sprintf("Failed ID: %s\n", res.ID))
		}
	}

	accuracy := (float64(passCount) / float64(totalCount)) * 100
	fmt.Println("\n\n--- EVALUATION COMPLETE ---")
	fmt.Printf("Time Taken: %v\n", time.Since(start))
	fmt.Printf("Total Samples: %d\n", totalCount)
	fmt.Printf("Passing: %d\n", passCount)
	fmt.Printf("Failing: %d\n", totalCount-passCount)
	fmt.Printf("FINAL SCORE: %.2f%%\n", accuracy)
}
