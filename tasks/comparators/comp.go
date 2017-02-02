// Package comparators provide a list of tasks useful as text comparators.
package comparators

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/natsukagami/kjudge-api-go/task"
	"github.com/natsukagami/kjudge-api-go/task/queue"
)

// Result is a struct that holds a comparator's output
type Result struct {
	Score   float64
	Comment string
}

// Error represents a diff checker error
type Error struct {
	Stderr   string
	ExitCode int
}

func (e Error) Error() string {
	return fmt.Sprintf("Comparator exited %d: %s", e.ExitCode, e.Stderr)
}

// Diff uses the built-in GNU diff command to compare output files.
// It automatically ignores all whitespaces.
// Note that output and answer arguments are filepaths rather than
// text contents.
func Diff(output, answer string) (r Result, e error) {
	tsk := task.NewTask("diff", []string{"-wq", output, answer}, "")
	res := queue.Enqueue(&tsk)
	if res.ExitCode == 0 {
		r = Result{1.0, "Output is Correct"}
	} else if res.ExitCode == 1 {
		r = Result{0.0, "Output isn't Correct"}
	} else {
		e = Error{res.Stderr, res.ExitCode}
	}
	return
}

// Comparator uses the (compiled?) comparator program to compare output files.
// The comparator should take 3 arguments [input] [output] [answer] as filepaths
// and writes only ONE number to stdout that is a float between 0 and 1 as the
// score, and write anything on stderr as the output.
func Comparator(problemFolder, input, output, answer string) (r Result, e error) {
	tsk := task.NewTask("./compare", []string{input, output, answer}, problemFolder)
	res := queue.Enqueue(&tsk)
	if res.ExitCode == 0 {
		if x, err := strconv.ParseFloat(strings.Replace(res.Stdout, "\n", "", -1), 64); e == nil {
			r.Score = math.Max(0, math.Min(1, x))
		} else {
			e = err
			return
		}
		r.Comment = res.Stderr
	} else {
		e = Error{res.Stderr, res.ExitCode}
	}
	return
}
