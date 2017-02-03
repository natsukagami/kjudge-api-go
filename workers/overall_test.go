package workers

import (
	"fmt"
	"testing"

	kjudge "github.com/natsukagami/kjudge-api-go"
)

func TestOverall(t *testing.T) {
	problem := kjudge.Problem{
		Name:        "rolls1",
		DisplayName: "Rolls 1",
		Folder:      "/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/",
		Comparator:  false,
		Grader:      false,
		Tests:       make([]kjudge.Test, 0),
		ScoringMode: kjudge.ScoreGroupMul,
		SubtaskScoring: []kjudge.SubtaskScore{
			kjudge.SubtaskScore{Num: 30, Score: 5},
			kjudge.SubtaskScore{Num: 30, Score: 5},
		},
		Time: 1000,
		Mem:  262144,
	}
	for i := 1; i <= 60; i++ {
		problem.Tests = append(problem.Tests, kjudge.Test{
			Input:  fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.in", i),
			Output: fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.out", i),
			Score:  1.0,
		})
	}
	sub := kjudge.Submission{
		Folder:  "/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1",
		Problem: problem,
	}
	Input <- &sub
	<-Output
	if sub.Result.Score != 10.0 {
		t.Error("Score unexpected")
	}
}

func BenchmarkOverall(b *testing.B) {
	problem := kjudge.Problem{
		Name:        "rolls1",
		DisplayName: "Rolls 1",
		Folder:      "/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/",
		Comparator:  false,
		Grader:      false,
		Tests:       make([]kjudge.Test, 0),
		ScoringMode: kjudge.ScoreGroupMul,
		SubtaskScoring: []kjudge.SubtaskScore{
			kjudge.SubtaskScore{Num: 5, Score: 5},
			kjudge.SubtaskScore{Num: 5, Score: 5},
		},
		Time: 1000,
		Mem:  262144,
	}
	for i := 1; i <= 10; i++ {
		problem.Tests = append(problem.Tests, kjudge.Test{
			Input:  fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.in", i),
			Output: fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.out", i),
			Score:  1.0,
		})
	}
	sub := kjudge.Submission{
		Folder:  "/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1",
		Problem: problem,
	}
	for i := 0; i < b.N; i++ {
		go func() { Input <- &sub }()
	}
	for i := 0; i < b.N; i++ {
		b.Log(<-Output)
	}
	// if sub.Result.Score != 10.0 {
	// 	b.Error("Score unexpected")
	// }
}
