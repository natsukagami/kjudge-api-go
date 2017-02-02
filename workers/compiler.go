// Package workers provides a series of channel workers
package workers

import (
	kjudge "github.com/natsukagami/kjudge-api-go"
	"github.com/natsukagami/kjudge-api-go/tasks/languages"
)

// CompilersCount represents the number of concurrent compilers.
const compilersCount = 1

// CompileFailed is a wrapper for compile-failed submissions.
type compileFailed struct {
	*kjudge.Submission
	err error
}

// Sub returns the submission pointer
func (c compileFailed) Sub() *kjudge.Submission {
	return c.Submission
}

func (c compileFailed) Error() string {
	return "Compilation failed: " + c.err.Error()
}

func compiler(in <-chan *kjudge.Submission, success chan<- *kjudge.Submission, fail chan<- failure) {
	for sub := range in {
		lang := sub.Language()
		if lang == nil {
			fail <- compileFailed{sub, languages.NoLanguageError{}}
			continue
		}
		var err error
		if sub.Problem.Grader {
			err = lang.CompileGrader(sub.Problem.Name, sub.Folder, sub.Problem.Folder)
		} else {
			err = lang.Compile(sub.Problem.Name, sub.Folder)
		}
		if err != nil {
			fail <- compileFailed{sub, err}
		} else {
			success <- sub
		}
	}
}
