package kjudge

import (
	"os"
	"path"
)

// ScoringMode represents how the problem is scored.
type ScoringMode byte

//
const (
	ScoreSum ScoringMode = iota
	ScoreGroupMin
	ScoreGroupMul
)

// SubtaskScore represents a Subtask's score configuration
type SubtaskScore struct {
	Num   int
	Score float64
}

// Problem represents a problem model.
type Problem struct {
	Name           string
	DisplayName    string
	Folder         string
	Comparator     bool
	Grader         bool
	Tests          []Test
	ScoringMode    ScoringMode
	SubtaskScoring []SubtaskScore
	Time           int64
	Mem            int64
}

// ComparatorFound returns if the comparator is already compiled
func (p Problem) ComparatorFound() bool {
	_, err := os.Stat(path.Join(p.Folder, "compare"))
	return err == nil
}
