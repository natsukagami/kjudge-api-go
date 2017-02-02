package workers

import kjudge "github.com/natsukagami/kjudge-api-go"

// Input is the place you push all the submissions
var Input = make(chan *kjudge.Submission)

// Output is the place you get judged submissions out
var Output = make(chan *kjudge.Submission)

func init() {
	compiled := make(chan *kjudge.Submission)
	evaluated := make(chan testingSuccess)
	fail := make(chan failure)
	for i := 0; i < compilersCount; i++ {
		go compiler(Input, compiled, fail)
	}
	for i := 0; i < testersCount; i++ {
		go tester(compiled, evaluated, fail)
	}
	for i := 0; i < testRunnersCount; i++ {
		go testRunner()
	}
	for i := 0; i < scorersCount; i++ {
		go scorer(evaluated, Output, fail)
	}
	for i := 0; i < failHandlersCount; i++ {
		go failHandler(fail, Output)
	}
}
