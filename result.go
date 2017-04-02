package kjudge

// TestResult represents the outcome of a test
type TestResult struct {
	Verdict string  `json:"verdict"`
	Score   float64 `json:"score"`
	Time    int64   `json:"runningTime"`
	Mem     int64   `json:"memoryUsed"`
}

// SubtaskResult represents the outcome of a subtask
type SubtaskResult struct {
	Score float64       `json:"score"`
	Tests []*TestResult `json:"tests"`
}

// Result represents the outcome of a submission
type Result struct {
	Score    float64          `json:"score"`
	Subtasks []*SubtaskResult `json:"subtasks"`
}
