package eval

// EvalCase represents a single evaluation scenario.
type EvalCase struct {
	ID        string `json:"id"`
	Prompt    string `json:"prompt"`
	Reference string `json:"reference"`
}

// GeneratedCase represents a completed attempt by a model.
type GeneratedCase struct {
	EvalCase
	SystemAnswer string `json:"system_answer"`
}

// EvalResult represents the verdict of a judge on a generated case.
type EvalResult struct {
	ID     string
	Passed bool
	Reason string
}
