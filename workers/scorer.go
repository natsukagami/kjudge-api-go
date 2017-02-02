package workers

import (
	"math"

	kjudge "github.com/natsukagami/kjudge-api-go"
)

const scorersCount = 1

// ScoreError is an error that happens when scoring failed
type scoreError struct {
	message string
}

func (e scoreError) Error() string {
	return "Scoring failed: " + e.message
}

// ScoringFailed wraps a Submission that failed scoring
type scoringFailed struct {
	*kjudge.Submission
	err error
}

// Sub returns the submission
func (s scoringFailed) Sub() *kjudge.Submission {
	return s.Submission
}

func (s scoringFailed) Error() string {
	return s.err.Error()
}

func scorer(in <-chan testingSuccess, success chan<- *kjudge.Submission, fail chan<- failure) {
	for sub := range in {
		var err error
		switch sub.Problem.ScoringMode {
		case kjudge.ScoreSum:
			sub.Result, err = sumScoring(sub.Results, sub.Problem.Tests)
		case kjudge.ScoreGroupMin:
			sub.Result, err = groupRatioScoring(sub.Results, sub.Problem.Tests, sub.Problem.SubtaskScoring, func(x float64, y float64) float64 {
				return math.Min(x, y)
			})
		case kjudge.ScoreGroupMul:
			sub.Result, err = groupRatioScoring(sub.Results, sub.Problem.Tests, sub.Problem.SubtaskScoring, func(x float64, y float64) float64 {
				return x * y
			})
		}
		if err != nil {
			fail <- scoringFailed{sub.Submission, err}
			continue
		}
		success <- sub.Submission
	}
}

func sumScoring(results []kjudge.TestResult, tests []kjudge.Test) (kjudge.Result, error) {
	var n = len(results)
	if n != len(tests) {
		return kjudge.Result{}, scoreError{"The number of tests is not consistent?"}
	}
	res := kjudge.Result{
		Score: 0,
		Subtasks: []kjudge.SubtaskResult{
			kjudge.SubtaskResult{Score: 0, Tests: results},
		},
	}
	for i := 0; i < n; i++ {
		results[i].Score *= tests[i].Score
		res.Score += results[i].Score
		res.Subtasks[0].Score += results[i].Score
	}
	return res, nil
}

func groupRatioScoring(results []kjudge.TestResult, tests []kjudge.Test, scoring []kjudge.SubtaskScore, comb func(float64, float64) float64) (kjudge.Result, error) {
	var n = len(results)
	if n != len(tests) {
		return kjudge.Result{}, scoreError{"The number of tests is not consistent?"}
	}
	res := kjudge.Result{Score: 0, Subtasks: make([]kjudge.SubtaskResult, len(scoring))}
	cur := 0
	for id, score := range scoring {
		ratio := 1.0
		res.Subtasks[id] = kjudge.SubtaskResult{Score: 0, Tests: results[cur : cur+score.Num]}
		if len(res.Subtasks[id].Tests) != score.Num {
			return kjudge.Result{}, scoreError{"The number of tests is not consistent?"}
		}
		for i := 0; i < score.Num; i++ {
			cur++
			ratio = comb(ratio, results[cur].Score)
		}
		res.Subtasks[id].Score = ratio * score.Score
		res.Score += ratio * score.Score
	}
	return res, nil
}
