package workers

import (
	"path"

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
	return "Testing failed: " + c.err.Error()
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

var testRuns = make(chan testingItem)

func tester(in <-chan *kjudge.Submission, success chan<- testingSuccess, fail chan<- failure) {
	for sub := range in {
		out := make(chan testingItem)
		res := make([]kjudge.TestResult, len(sub.Problem.Tests))
		for i := range sub.Problem.Tests {
			go func(id int) {
				testRuns <- testingItem{Submission: sub, TestID: id, Out: out}
			}(i)
		}
		var err error
		for i := 0; i < len(sub.Problem.Tests); i++ {
			r := <-out
			if r.Err != nil {
				err = r.Err
				break
			} else {
				res[r.TestID] = r.Result
			}
		}
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

func folderClean(folder string) {
	go fs.Remove(folder)
}

func isolateClean(box *runners.Isolate) {
	go box.Cleanup()
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
	defer folderClean(folder)
	isolate := runners.Make()
	defer isolateClean(isolate)
	// Step 2. Copy the input files
	{
		c := make(chan error)
		go func() { c <- fs.Copy(t.Input, path.Join(folder, "input.txt")) }()
		go func() { c <- fs.Copy(path.Clean(sub.Folder)+"/.", folder) }()
		go func() { c <- fs.Chmod(folder, "777") }()
		go func() { c <- isolate.Prepare() }()
		for i := 0; i < 4; i++ {
			if err := <-c; err != nil {
				test.Err = err
				return test
			}
		}
	}
	// Step 3: Runs the code in a sandbox
	r, err := isolate.Run(
		sub.Language().RunCommand(sub.Problem.Name),
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
		v, err = comparators.Diff(path.Join(folder, "output.txt"), t.Output)
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
