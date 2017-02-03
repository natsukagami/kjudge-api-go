// Package queue provides an universal queue of tasks
package queue

import "github.com/natsukagami/kjudge-api-go/task"

const runners = 20

// Item represents a queue item
type item struct {
	*task.Task
	ch chan task.Result
}

var (
	input    = make(chan item)
	block    = make(chan bool)
	unblock  = make(chan bool)
	priorize = make(chan item)
)

// Enqueue adds the specified task into the queue
func Enqueue(t *task.Task) task.Result {
	ch := make(chan task.Result)
	go func() {
		input <- item{t, ch}
	}()
	return <-ch
}

// PriorizedEnqueue adds the task into the priorized channel.
// It runs only when no other tasks are running
func PriorizedEnqueue(t *task.Task) task.Result {
	ch := make(chan task.Result)
	go func() {
		priorize <- item{t, ch}
	}()
	return <-ch
}

func runTask(x *item) {
	res := x.Run()
	x.ch <- res
	close(x.ch)
}

func taskRunner(i <-chan item) {
	for {
		select {
		case x := <-i:
			runTask(&x)
		case <-block:
			<-unblock
		}
	}
}

func priorizedRunner(i <-chan item) {
	for x := range i {
		for n := 0; n < runners; n++ {
			block <- true
		}
		runTask(&x)
		for n := 0; n < runners; n++ {
			unblock <- true
		}
	}
}

func init() {
	// Initialize workers, because channels are synchronous
	for i := 0; i < runners; i++ {
		go taskRunner(input)
	}
	go priorizedRunner(priorize)
}
