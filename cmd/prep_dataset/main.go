package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dute/go_evals/pkg/eval"
)

// HotpotRaw represents the input structure from HotpotQA.
type HotpotRaw struct {
	ID       string          `json:"_id"`
	Question string          `json:"question"`
	Answer   string          `json:"answer"`
	Level    string          `json:"level"`
	Context  [][]interface{} `json:"context"`
}

func main() {
	inputPath := flag.String("input", "data/hotpot_dev.json", "Path to raw HotpotQA dataset")
	outputPath := flag.String("output", "data/golden_set.jsonl", "Path to save the golden set")
	sampleLimit := flag.Int("limit", 50, "Number of hard cases to sample")
	flag.Parse()

	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("Reading raw dataset...")

	var rawData []HotpotRaw
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rawData); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	var hardCases []HotpotRaw
	for _, item := range rawData {
		if item.Level == "hard" {
			hardCases = append(hardCases, item)
		}
	}

	if len(hardCases) == 0 {
		fmt.Println("No 'hard' cases found in the dataset.")
		return
	}

	// Shuffle to get random samples
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(hardCases), func(i, j int) { hardCases[i], hardCases[j] = hardCases[j], hardCases[i] })

	limit := *sampleLimit
	if len(hardCases) < limit {
		limit = len(hardCases)
	}
	subset := hardCases[:limit]
	fmt.Printf("Selected %d 'hard' samples for Golden Set.\n", len(subset))

	outFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()
	writer := json.NewEncoder(outFile)

	for _, raw := range subset {
		flatContext := eval.FlattenContext(raw.Context)
		fullPrompt := fmt.Sprintf("Context:\n%s\n\nQuestion: %s", flatContext, raw.Question)

		evalCase := eval.EvalCase{
			ID:        raw.ID,
			Prompt:    fullPrompt,
			Reference: raw.Answer,
		}

		writer.Encode(evalCase)
	}

	fmt.Printf("Done! Saved to %s\n", *outputPath)
}
