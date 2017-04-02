package kjudge

// Test represents a test model.
type Test struct {
	Input  string  `json:"input"`
	Output string  `json:"output"`
	Score  float64 `json:"score"`
}
