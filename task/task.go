// Package task provides the Task struct and interface.
package task

import (
	"bytes"
	"os/exec"
	"syscall"

	debug "github.com/tj/go-debug"
)

var taskDebug = debug.Debug("kjudge:task")

// Result is the output of execution of a Task.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Task represents an unit-level operation of the package.
// Its job is acually a direct fork to Linux command line programs.
type Task struct {
	Command string
	Args    []string
	Cwd     string
	id      int
}

// Run executes the command synchronorously.
// The result is a Task.Result object.
func (t *Task) Run() (r Result) {
	taskDebug("Task %d: Running `%s`", t.id, t.Command)
	osCmd := exec.Command(t.Command, t.Args...)
	var stdout, stderr bytes.Buffer
	osCmd.Stderr = &stderr
	osCmd.Stdout = &stdout
	err := osCmd.Run()
	taskDebug("Task %d: Done", t.id)
	r.Stdout = stdout.String()
	r.Stderr = stderr.String()
	if err == nil {
		r.ExitCode = 0
	} else {
		switch err.(type) {
		case *exec.ExitError:
			r.ExitCode = (err.(*exec.ExitError)).Sys().(syscall.WaitStatus).ExitStatus()
		default:
			r.ExitCode = -1
		}
	}
	return
}

// ID returns the task's id
func (t Task) ID() int {
	return t.id
}

var index = 0

// NewTask creates a new Task. Use this function to
// initialize the Task with an ID.
func NewTask(command string, args []string, cwd string) Task {
	index++
	return Task{
		Command: command,
		Args:    args,
		Cwd:     cwd,
		id:      index}
}
