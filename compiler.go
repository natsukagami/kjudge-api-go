package kjudge

import "github.com/pkg/errors"

// Compile runs the compiling of the submission.
func (sub *Submission) Compile(folder string) error {
	if sub.Problem.Grader {
		return errors.Wrap(
			sub.Language.CompileGrader(sub.Problem.Name, folder, sub.Problem.Folder),
			"Compilation Error")
	}
	return errors.Wrap(
		sub.Language.Compile(sub.Problem.Name, folder),
		"Compilation Error")
}

// func compiler(in <-chan *Submission, success chan<- *Submission, fail chan<- failure) {
// 	for sub := range in {
// 		lang := sub.Language()
// 		if lang == nil {
// 			fail <- compileFailed{sub, languages.NoLanguageError{}}
// 			continue
// 		}
// 		var err error
// 		if sub.Problem.Grader {
// 			err = lang.CompileGrader(sub.Problem.Name, sub.Folder, sub.Problem.Folder)
// 		} else {
// 			err = lang.Compile(sub.Problem.Name, sub.Folder)
// 		}
// 		if err != nil {
// 			fail <- compileFailed{sub, err}
// 		} else {
// 			success <- sub
// 		}
// 	}
// }
