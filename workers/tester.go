package workers

import (
	"path"
	"sync"

	kjudge "github.com/natsukagami/kjudge-api-go"
	"github.com/natsukagami/kjudge-api-go/misc"
	"github.com/natsukagami/kjudge-api-go/tasks/comparators"
	"github.com/natsukagami/kjudge-api-go/tasks/fs"
	"github.com/natsukagami/kjudge-api-go/tasks/runners"
)

// TestingFailed wraps a Submission with testing errors
type testingFailed struct {
	*kjudge.Submission
	err error
}

// Sub returns the submission pointer
func (c testingFailed) Sub() *kjudge.Submission {
	return c.Submission
}

func (c testingFailed) Error() string {
	return "Compilation failed: " + c.err.Error()
}

// TestingSuccess wraps a Submission with successful testing, adding
// an array of raw test results.
type testingSuccess struct {
	*kjudge.Submission
	Results []kjudge.TestResult
}

// TestingItem wraps a Submission with a pending test.
type testingItem struct {
	*kjudge.Submission
	TestID int
	Result kjudge.TestResult
	Err    error
	Out    chan<- testingItem
}

const (
	testersCount     = 2
	testRunnersCount = 7
)

var testRuns = make(chan testingItem)

func tester(in <-chan *kjudge.Submission, success chan<- testingSuccess, fail chan<- failure) {
	for sub := range in {
		out := make(chan testingItem)
		res := make([]kjudge.TestResult, len(sub.Problem.Tests))
		var err error
		wg := sync.WaitGroup{}
		wg.Add(len(sub.Problem.Tests))
		for i := range sub.Problem.Tests {
			go func(id int) {
				testRuns <- testingItem{Submission: sub, TestID: id, Out: out}
				r := <-out
				if r.Err != nil {
					err = r.Err
				} else {
					res[r.TestID] = r.Result
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		if err != nil {
			fail <- testingFailed{sub, err}
		} else {
			success <- testingSuccess{sub, res}
		}
	}
}

func testRunner() {
	for test := range testRuns {
		test.Out <- runTest(test)
	}
}

func runTest(test testingItem) testingItem {
	sub := test.Submission
	t := &sub.Problem.Tests[test.TestID]
	// Runs the test
	// Step 1. Create a temp folder
	folder := path.Join("/tmp", misc.RandString(20))
	if err := fs.Mkdir(folder); err != nil {
		test.Err = err
		return test
	}
	// Step 2. Copy the input files
	if err := fs.Copy(t.Input, path.Join(folder, "input.txt")); err != nil {
		test.Err = err
		return test
	}
	if err := fs.Copy(path.Clean(sub.Folder)+"/.", folder); err != nil {
		test.Err = err
		return test
	}
	if err := fs.Chmod(folder, "777"); err != nil {
		test.Err = err
		return test
	}
	// Step 3: Runs the code in a sandbox
	isolate := runners.Make()
	defer isolate.Cleanup()
	if err := isolate.Prepare(); err != nil {
		test.Err = err
		return test
	}
	r, err := isolate.Run(
		path.Join(folder, sub.Problem.Name),
		folder,
		sub.Problem.Time,
		sub.Problem.Mem,
	)
	if err != nil {
		test.Err = err
		return test
	}
	if r.Status != "OK" {
		test.Result = kjudge.TestResult{
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
		v, err = comparators.Diff(path.Join(folder, "out.txt"), t.Output)
	}
	if err != nil {
		test.Err = err
		return test
	}
	test.Result = kjudge.TestResult{
		Verdict: v.Comment,
		Score:   v.Score,
		Time:    r.Time,
		Mem:     r.Mem,
	}
	return test
}
