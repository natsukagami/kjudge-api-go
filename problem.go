package kjudge

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/natsukagami/kjudge-api-go/lib/languages"
	"github.com/pkg/errors"
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
	Name           string         `json:"name"`
	DisplayName    string         `json:"displayName"`
	Folder         string         `json:"folder"`
	Comparator     bool           `json:"comparator"`
	Grader         bool           `json:"grader"`
	Tests          []*Test        `json:"tests"`
	ScoringMode    ScoringMode    `json:"scoringMode"`
	SubtaskScoring []SubtaskScore `json:"subtaskScore"`
	Time           int64          `json:"timeLimit"`
	Mem            int64          `json:"memoryLimit"`
}

// ComparatorFound returns if the comparator is already compiled
func (p *Problem) ComparatorFound() bool {
	_, err := os.Stat(path.Join(p.Folder, "compare"))
	return err == nil
}

// CompileComparator ensures the comparator is available in the problem's folder.
func (p *Problem) CompileComparator() error {
	if p.ComparatorFound() {
		return nil
	}
	cpp, err := languages.New(".cpp")
	if err != nil {
		return errors.Wrap(err, "Comparator compile error")
	}
	return errors.Wrap(
		cpp.CompileComparator(p.Folder),
		"Comparator compile error")
}

// Validate runs a validation test on the problem.
// It makes sure that:
// - There are enough tests configured each subtask.
// - Test input/output files must exists.
// - If comparator is enabled, it should be compiled and ready.
// - If grader is enabled, grading libraries should be there. That includes `grader.cpp` for the
// grader implementation, `grader.h` for the header information and `grader.zip` for a sample grader
// implementation that contestants can download.
func (p *Problem) Validate() error {
	// Subtask validation
	if p.ScoringMode == ScoreGroupMin || p.ScoringMode == ScoreGroupMul {
		var ptr int
		for id, subtask := range p.SubtaskScoring {
			if ptr+subtask.Num > len(p.Tests) {
				return errors.Wrap(
					fmt.Errorf("Subtask %d: Not enough test", id+1),
					"Problem validation Error")
			}
			ptr += subtask.Num
		}
	}
	// Test files validation
	for id, test := range p.Tests {
		if _, err := os.Stat(test.Input); err != nil {
			return errors.Wrap(
				errors.Wrapf(err, "Test %d, input file", id),
				"Problem validation Error")
		}
		if _, err := os.Stat(test.Output); err != nil {
			return errors.Wrap(
				errors.Wrapf(err, "Test %d, output file", id),
				"Problem validation Error")
		}
	}
	// Comparator
	if p.Comparator {
		if err := p.CompileComparator(); err != nil {
			return errors.Wrapf(err, "Problem validation Error")
		}
	}
	// Grader
	if p.Grader {
		for _, file := range [3]string{"grader.cpp", "grader.h", "grader.zip"} {
			if _, err := os.Stat(filepath.Join(p.Folder, file)); err != nil {
				return errors.Wrapf(err, "Problem validation Error")
			}
		}
	}
	return nil
}
