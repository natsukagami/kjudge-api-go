// Package queue provides an universal queue of tasks
package queue

import "github.com/natsukagami/kjudge-api-go/task"

const runners = 7

// Item represents a queue item
type item struct {
	*task.Task
	ch chan task.Result
}

var input = make(chan item)

// Enqueue adds the specified task into the queue
func Enqueue(t *task.Task) task.Result {
	ch := make(chan task.Result)
	go func() {
		input <- item{t, ch}
	}()
	return <-ch
}

func taskRunner(i <-chan item) {
	for x := range i {
		res := x.Run()
		x.ch <- res
		close(x.ch)
	}
}

func main() {
	// Initialize workers, because channels are synchronous
	for i := 0; i < runners; i++ {
		go taskRunner(input)
	}
}
