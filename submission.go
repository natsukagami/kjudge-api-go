package kjudge

import (
	"os"
	"path"

	"github.com/natsukagami/kjudge-api-go/tasks/languages"
)

// Submission represents a submission model.
type Submission struct {
	Folder     string
	Problem    Problem
	Result     Result
	JudgeError string
}

// Language returns the language of the submission, or nil
// if none is supported.
func (s Submission) Language() languages.Interface {
	for _, l := range languages.All {
		if _, e := os.Stat(path.Join(s.Folder, s.Problem.Name+l.Ext())); e == nil {
			return l
		}
	}
	return nil
}
