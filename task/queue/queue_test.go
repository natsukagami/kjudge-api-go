package queue

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/natsukagami/kjudge-api-go/task"
)

func TestQueue(t *testing.T) {
	main()
	var n = 10
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		go func(id int) {
			tsk := task.NewTask(fmt.Sprintf("echo"), []string{fmt.Sprintf("%d", id)}, "")
			num := Enqueue(&tsk)
			if x, e := strconv.ParseInt(strings.Replace(num.Stdout, "\n", "", -1), 10, 32); e != nil || x != int64(id) {
				t.Error("Invalid pattern: " + e.Error())
			} else {
				t.Log(id, " is correct")
			}
			wg.Done()
		}(i)
	}
	wg.Add(n)
	wg.Wait()
}

func BenchmarkQueue(b *testing.B) {
	main()
	var n = b.N
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			tsk := task.NewTask(fmt.Sprintf("echo"), []string{fmt.Sprintf("%d", id)}, "")
			num := Enqueue(&tsk)
			if x, e := strconv.ParseInt(strings.Replace(num.Stdout, "\n", "", -1), 10, 32); e != nil || x != int64(id) {
				b.Error("Invalid pattern: " + e.Error())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
