package kjudge

import (
	"io/ioutil"
	"path"

	"github.com/natsukagami/kjudge-api-go/lib/comparators"
	"github.com/natsukagami/kjudge-api-go/lib/fs"
	"github.com/natsukagami/kjudge-api-go/lib/runners"
	"github.com/pkg/errors"
)

// Defines the number of concurrent test runners.
const (
	TestRunners = 2000
)

// TestingItem wraps a Submission with a pending test.
type testingItem struct {
	*Submission
	Folder string
	TestID int
	Result *TestResult
	Err    error
	Out    chan<- *testingItem
}

var testRuns = make(chan *testingItem)

// RunTests performs test-running on the submission.
func (s *Submission) RunTests(folder string) (t []*TestResult, err error) {
	out := make(chan *testingItem)
	t = make([]*TestResult, len(s.Problem.Tests))
	for i := range s.Problem.Tests {
		go func(id int) {
			testRuns <- &testingItem{Submission: s, TestID: id, Out: out, Folder: folder}
		}(i)
	}
	for i := 0; i < len(t); i++ {
		item := <-out
		if item.Err != nil {
			err = errors.Wrap(
				errors.Wrapf(item.Err, "Test %d", item.TestID),
				"Testing Error")
		}
		t[item.TestID] = item.Result
	}
	return
}

func testRunner() {
	for test := range testRuns {
		test.Out <- runTest(test)
	}
}

func init() {
	for i := 0; i < TestRunners; i++ {
		go testRunner()
	}
}

func isolateClean(box *runners.Isolate) {
	go box.Cleanup()
}

func runTest(test *testingItem) *testingItem {
	sub := test.Submission
	t := sub.Problem.Tests[test.TestID]
	// Runs the test
	// Step 1. Create a temp folder
	folder, err := ioutil.TempDir("", "tester-")
	if err != nil {
		test.Err = err
		return test
	}
	defer folderClean(folder)
	isolate := runners.New()
	defer isolateClean(isolate)
	// Step 2. Copy the input files
	{
		c := make(chan error)
		go func() { c <- fs.Copy(t.Input, path.Join(folder, "input.txt")) }()
		go func() { c <- fs.Copy(test.Folder+"/.", folder) }()
		go func() { c <- isolate.Prepare() }()
		for i := 0; i < 3; i++ {
			if err := <-c; err != nil {
				test.Err = err
				return test
			}
		}
		if err := fs.Chmod(folder, "777"); err != nil {
			test.Err = err
			return test
		}
	}
	// Step 3: Runs the code in a sandbox
	r, err := isolate.Run(
		sub.Language.Executable(sub.Problem.Name, "."),
		folder,
		sub.Problem.Time,
		sub.Problem.Mem,
	)
	if err != nil {
		test.Err = err
		return test
	}
	if r.Status != "OK" {
		test.Result = &TestResult{
			Verdict: runners.SandboxError(r.Status),
			Score:   0.0,
			Time:    r.Time,
			Mem:     r.Mem}
		return test
	}
	// Step 4: Compare the outputs
	var v comparators.Result
	if sub.Problem.Comparator {
		v, err = comparators.Comparator(sub.Problem.Folder, t.Input, path.Join(folder, "output.txt"), t.Output)
	} else {
		v, err = comparators.Diff(path.Join(folder, "output.txt"), t.Output)
	}
	if err != nil {
		test.Err = err
		return test
	}
	test.Result = &TestResult{
		Verdict: v.Comment,
		Score:   v.Score,
		Time:    r.Time,
		Mem:     r.Mem,
	}
	return test
}
