.PHONY: prep generate export align run test tidy clean

prep:
	go run cmd/prep_dataset/main.go

generate:
	go run cmd/generate/main.go

export:
	go run cmd/export_csv/main.go

align:
	go run cmd/align_judge/main.go

run:
	go run cmd/run_evals/main.go

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -f data/golden_set.jsonl data/eval_dataset_generated.jsonl data/manual_review.csv data/failures.log
