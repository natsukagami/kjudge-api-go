package workers

import kjudge "github.com/natsukagami/kjudge-api-go"

// Failure represents a failed submission
type failure interface {
	Sub() *kjudge.Submission
	Error() string
}

const failHandlersCount = 7

func failHandler(in <-chan failure, out chan<- *kjudge.Submission) {
	for sub := range in {
		sub.Sub().JudgeError = sub.Error()
		out <- sub.Sub()
	}
}
