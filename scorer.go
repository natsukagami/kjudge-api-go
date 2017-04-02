package kjudge

import (
	"math"

	"github.com/pkg/errors"
)

// Scoring is implemented quite obviously.

// AssignScore gives a submission its score by calculating through test results.
func (s *Submission) AssignScore(res []*TestResult) (err error) {
	switch s.Problem.ScoringMode {
	case ScoreSum:
		s.Result, err = sumScoring(res, s.Problem.Tests)
	case ScoreGroupMin:
		s.Result, err = groupRatioScoring(
			res,
			s.Problem.Tests,
			s.Problem.SubtaskScoring,
			math.Min)
	case ScoreGroupMul:
		s.Result, err = groupRatioScoring(
			res,
			s.Problem.Tests,
			s.Problem.SubtaskScoring,
			func(x, y float64) float64 { return x * y })
	}
	return errors.Wrap(err, "Scoring Error")
}

func sumScoring(results []*TestResult, tests []*Test) (*Result, error) {
	var n = len(results)
	if n != len(tests) {
		return nil, errors.New("The number of tests is not consistent?")
	}
	res := &Result{
		Score: 0,
		Subtasks: []*SubtaskResult{
			&SubtaskResult{Score: 0, Tests: results},
		},
	}
	for i := 0; i < n; i++ {
		results[i].Score *= tests[i].Score
		res.Score += results[i].Score
		res.Subtasks[0].Score += results[i].Score
	}
	return res, nil
}

func groupRatioScoring(results []*TestResult,
	tests []*Test,
	scoring []SubtaskScore,
	comb func(float64, float64) float64,
) (*Result, error) {
	var n = len(results)
	if n != len(tests) {
		return nil, errors.New("The number of tests is not consistent?")
	}
	res := &Result{Score: 0, Subtasks: make([]*SubtaskResult, len(scoring))}
	cur := 0
	for id, score := range scoring {
		ratio := 1.0
		res.Subtasks[id] = &SubtaskResult{Score: 0, Tests: results[cur : cur+score.Num]}
		for i := 0; i < score.Num; i++ {
			ratio = comb(ratio, results[cur].Score)
			cur++
		}
		res.Subtasks[id].Score = ratio * score.Score
		res.Score += ratio * score.Score
	}
	return res, nil
}
