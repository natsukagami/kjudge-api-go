// Package languages provides compiler-specific tasks
package languages

import (
	"fmt"

	"github.com/natsukagami/kjudge-api-go/task"
	"github.com/natsukagami/kjudge-api-go/task/queue"
)

// Interface represents a compiler's set of functions
type Interface interface {
	Compile(name, folder string) error
	CompileGrader(name, folder, problemFolder string) error
	CompileComparator(problemFolder string) error
	Ext() string
}

// Error is a compiler error
type Error struct {
	Err      string
	Exitcode int
}

func (e Error) Error() string {
	return fmt.Sprintf("Compile error (exitcode %d): %s", e.Exitcode, e.Err)
}

func doCompile(t *task.Task) error {
	res := queue.Enqueue(t)
	if res.ExitCode != 0 {
		return Error{res.Stderr, res.ExitCode}
	}
	return nil
}

// NoLanguageError is an error thrown when a submission is found with no
// supported language.
type NoLanguageError struct {
}

func (e NoLanguageError) Error() string {
	return "No language found"
}

// All provides an array of supported languages
var All = []Interface{Cpp{}}
