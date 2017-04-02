// Package task also provides an universal queue of tasks.
package task

// Runners defines the number of concurrent task runners.
const (
	Runners = 2000
)

// Item represents a queue item
type item struct {
	*Task
	ch chan *Result
}

var (
	input    = make(chan item)
	block    = make(chan struct{})
	unblock  = make(chan struct{})
	priorize = make(chan item)
)

// Enqueue adds the specified task into the queue
func Enqueue(t *Task) *Result {
	ch := make(chan *Result)
	go func() {
		input <- item{t, ch}
	}()
	return <-ch
}

// PriorizedEnqueue adds the task into the priorized channel.
// It runs only when no other tasks are running
func PriorizedEnqueue(t *Task) *Result {
	ch := make(chan *Result)
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
		for n := 0; n < Runners; n++ {
			block <- struct{}{}
		}
		runTask(&x)
		for n := 0; n < Runners; n++ {
			unblock <- struct{}{}
		}
	}
}

func init() {
	// Initialize workers, because channels are synchronous
	for i := 0; i < Runners; i++ {
		go taskRunner(input)
	}
	go priorizedRunner(priorize)
}
