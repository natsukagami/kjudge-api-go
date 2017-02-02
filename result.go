package kjudge

// TestResult represents the outcome of a test
type TestResult struct {
	Verdict string
	Score   float64
	Time    int64
	Mem     int64
}

// SubtaskResult represents the outcome of a subtask
type SubtaskResult struct {
	Score float64
	Tests []TestResult
}

// Result represents the outcome of a submission
type Result struct {
	Score    float64
	Subtasks []SubtaskResult
}
