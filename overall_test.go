package kjudge

import (
	"fmt"
	"testing"
)

func TestOverall(t *testing.T) {
	problem := Problem{
		Name:        "rolls1",
		DisplayName: "Rolls 1",
		Folder:      "/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/",
		Comparator:  false,
		Grader:      false,
		Tests:       make([]*Test, 0),
		ScoringMode: ScoreGroupMul,
		SubtaskScoring: []SubtaskScore{
			SubtaskScore{Num: 3, Score: 5},
			SubtaskScore{Num: 3, Score: 5},
		},
		Time: 1000,
		Mem:  262144,
	}
	for i := 1; i <= 6; i++ {
		problem.Tests = append(problem.Tests, &Test{
			Input:  fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.in", i),
			Output: fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.out", i),
			Score:  1.0,
		})
	}
	sub, err := NewSub("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/rolls1.cpp", &problem)
	if err != nil {
		t.Fatal(err)
	}
	sub.Judge()
	if sub.JudgeError != nil {
		t.Fatal(sub.JudgeError)
	}
	t.Log(sub)
	if sub.Result.Score != 10.0 {
		t.Error("Score unexpected")
	}
}

func BenchmarkOverall(b *testing.B) {
	problem := Problem{
		Name:        "rolls1",
		DisplayName: "Rolls 1",
		Folder:      "/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/",
		Comparator:  false,
		Grader:      false,
		Tests:       make([]*Test, 0),
		ScoringMode: ScoreGroupMul,
		SubtaskScoring: []SubtaskScore{
			SubtaskScore{Num: 3, Score: 5},
			SubtaskScore{Num: 0, Score: 5},
		},
		Time: 1000,
		Mem:  262144,
	}
	for i := 1; i <= 3; i++ {
		problem.Tests = append(problem.Tests, &Test{
			Input:  fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.in", i),
			Output: fmt.Sprintf("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/%02d.out", i),
			Score:  1.0,
		})
	}
	done := make(chan struct{})
	for i := 0; i < b.N; i++ {
		go func() {
			sub, err := NewSub("/home/natsukagami/MEGASync/doituyen_2016/20161014/ROLLS1/rolls1.cpp", &problem)
			if err != nil {
				b.Fatal(err)
			}
			sub.Judge()
			done <- struct{}{}
		}()
	}
	for i := 0; i < b.N; i++ {
		<-done
	}
	// if sub.Result.Score != 10.0 {
	// 	b.Error("Score unexpected")
	// }
}
