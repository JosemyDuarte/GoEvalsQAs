# Go Evals: A Product Evaluation Pipeline

This repository implements an iterative evaluation pipeline in Go, designed to improve LLM-based products through testing and human-in-the-loop calibration.

The methodology follows the principles outlined in **[Eugene Yan's "Evaluations for LLM-based Products"](https://eugeneyan.com/writing/product-evals/)** and uses [HotpotQA dataset](https://huggingface.co/datasets/hotpotqa/hotpot_qa) to show it in action.

## 🚀 Overview

Evaluating LLMs is hard because their outputs are non-deterministic. This pipeline solves this by:
1. **Defining a Golden Set**: High-quality, representative examples of expected behavior.
2. **Model-in-the-Loop Evaluation**: Using an "LLM-as-a-judge" to grade outputs at scale.
3. **Judge Alignment**: Calibrating the judge against human labels to ensure accuracy.
4. **Iterative Refinement**: Making it easy to swap models, prompts, and judging logic.

---

## 🏗 Architecture & Pipeline

The evaluation flow is divided into 5 distinct stages. You can run them manually or use the provided **`Makefile`**.

### Stage 1: Prep Dataset
**`make prep`** or `go run cmd/prep_dataset/main.go`
- Reads raw data (e.g., `data/hotpot_dev.json`).
- Filters for "hard" cases and samples them.
- **Flags**: `-input` (raw path), `-output` (golden path), `-limit` (sample size).

### Stage 2: Generate Answers
**`make generate`** or `go run cmd/generate/main.go`
- Takes the `golden_set.jsonl`.
- Runs the "Student" model via a CLI provider.
- **Flags**: `-input`, `-output`, `-cmd` (e.g., `ollama`), `-args` (e.g., `run,llama3.3`).

### Stage 3: Export for Manual Review
**`make export`** or `go run cmd/export_csv/main.go`
- Converts generated answers into a human-friendly CSV template.
- **Flags**: `-input`, `-output`.

> [!IMPORTANT]
> **Manual Intervention Required**: After Stage 3, open `data/manual_review.csv` in a spreadsheet tool. 
> 1. Grade the results: Enter `1` for correct, `0` for incorrect in the labeling column.
> 2. Save the result as **`data/manual_review_done.csv`** (or your preferred labeling path).
> 
> [!TIP]
> **Defect Recall Calibration**: To properly align the judge, you need failing cases ("0"). If your manual review has fewer than 5 failures, consider adding **synthetic cases** (manually create bad answers) to the CSV. This ensures the judge is "strict" enough to catch subtle errors.

### Stage 4: Align the Judge
**`make align`** or `go run cmd/align_judge/main.go`
- Compares the Judge's automated grades against your human labels.
- Use this to refine the judge's prompt until it matches human intuition.
- **Flags**: `-labels`, `-cmd`, `-args`.

### Stage 5: Run Evals at scale
**`make run`** or `go run cmd/run_evals/main.go`
- Uses the aligned judge to run the final evaluation on the dataset using a worker pool.
- **Flags**: `-input`, `-workers`, `-cmd`, `-args`.

### Stage 6: Comparative Evaluation (Optional)
This is where the pipeline truly shines. Once you have an aligned judge, you can run Stage 2 and Stage 5 for **different models** (e.g., Llama 3 vs. Gemma 2) and compare their final scores.

#### Example: Llama 3.3 (70B) vs. Gemma 3 (1B)
Running the exact same 45 "hard" cases through both models yields a clear winner:

| Model | Total Samples | Passing | Failing | **Final Score** |
| :--- | :---: | :---: | :---: | :---: |
| **Llama 3.3 (70B)** | 45 | 42 | 3 | **93.33%** |
| **Gemma 3 (1B)** | 45 | 23 | 22 | **51.11%** |

- A higher score from the judge indicates better factual alignment with your golden set.
- Use this to make data-driven decisions on which model to deploy for your specific product use-case.

<img width="2816" height="1536" alt="Gemini_Generated_Image_yka7b0yka7b0yka7" src="https://github.com/user-attachments/assets/c6d2702a-f56a-47e1-9b5e-17bd09fbb738" />

---

## 📁 Data Flow

| File | Description |
| :--- | :--- |
| `data/hotpot_dev.json` | Raw input dataset (HotpotQA). |
| `data/golden_set.jsonl` | Cleaned evaluation cases. |
| `data/eval_dataset_generated.jsonl` | Answers from the target model. |
| `data/manual_review.csv` | Template for human labeling. |
| `data/manual_review_done.csv` | Calibrated ground truth for the judge. |
| `data/failures.log` | Log of failing IDs (created during `make run`). |

---

## 🛠 Developer Guide

### Prerequisites
- Go 1.25+
- Ollama or Gemini-CLI

### Common Commands
- **Run Tests**: `make test`
- **Clean Artifacts**: `make clean`
- **Tidy Modules**: `make tidy`

### Setup
1. Clone the repo.
2. Download the main dataset:
   ```bash
   curl -o data/hotpot_dev.json http://curtis.ml.cmu.edu/datasets/hotpot/hotpot_dev_distractor_v1.json
   ```
3. Ensure you have your target model running (e.g., `ollama run llama3.3`).
4. Run the pipeline stages in order.

---

## 🧠 Philosophy: Iterative Eval

> "The goal of evals is not a single score, but a process to improve the product."

- **Don't wait for perfect data**: Start with 50-100 samples.
- **The Judge is a model too**: Treat the judge's prompt as code that needs unit tests (the `align_judge` stage).
- **Scale with Workers**: Use Go's concurrency to run hundreds of evals in seconds.
