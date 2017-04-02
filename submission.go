package kjudge

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/natsukagami/kjudge-api-go/lib/languages"
	"github.com/pkg/errors"
)

// Submission represents a submission model.
type Submission struct {
	Problem    *Problem           `json:"problem"`
	Content    []byte             `json:"content"`
	Ext        string             `json:"language"`
	Language   languages.Language `json:"-"`
	Result     *Result            `json:"result"`
	JudgeError error              `json:"judgeError"`
}

// NewSub creates a new submission parsed from a file given in path.
func NewSub(file string, prob *Problem) (s *Submission, err error) {
	if stat, err := os.Stat(file); err != nil {
		return nil, errors.Wrap(err, "Submission create error")
	} else if stat.IsDir() {
		return nil, errors.New("Submission create error: given path is a directory")
	}
	s = &Submission{
		Problem: prob,
		Result:  nil,
		Ext:     filepath.Ext(file),
	}
	if s.Content, err = ioutil.ReadFile(file); err != nil {
		return nil, errors.Wrap(err, "Submission create error")
	}
	if s.Language, err = languages.New(s.Ext); err != nil {
		return nil, errors.Wrap(err, "Submission create error")
	}
	return
}

// Judge performs judging on the current submission.
func (s *Submission) Judge() {
	if err := s.Problem.Validate(); err != nil {
		s.JudgeError = err
		return
	}
	folder, err := ioutil.TempDir("", "submission-")
	if err != nil {
		s.JudgeError = err
		return
	}
	defer folderClean(folder)
	if err := ioutil.WriteFile(
		filepath.Join(folder, s.Problem.Name+s.Ext),
		s.Content,
		0755,
	); err != nil {
		s.JudgeError = err
		return
	}
	if err := s.Compile(folder); err != nil {
		s.JudgeError = err
		return
	}
	testResult, err := s.RunTests(folder)
	if err != nil {
		s.JudgeError = err
		return
	}
	if err := s.AssignScore(testResult); err != nil {
		s.JudgeError = err
		return
	}
}
