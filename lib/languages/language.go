// Package languages provides compiler-specific tasks
package languages

import (
	"errors"
	"fmt"

	"github.com/natsukagami/kjudge-api-go/task"
)

// Language represents a compiler's set of functions
type Language interface {
	Compile(name, folder string) error
	CompileGrader(name, folder, problemFolder string) error
	CompileComparator(problemFolder string) error
	Executable(name string, folder string) string
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
	res := task.Enqueue(t)
	if res.ExitCode != 0 {
		return Error{res.Stderr, res.ExitCode}
	}
	return nil
}

// New returns a language instance, based on the provided extension.
func New(ext string) (l Language, err error) {
	switch ext {
	case ".cpp", ".cc":
		l = &cpp{}
	default:
		err = errors.New("No language found for ext: " + ext)
	}
	return
}
