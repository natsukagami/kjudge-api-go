// Package fs provide a library of common filesystem tasks
package fs

import (
	"fmt"

	"github.com/natsukagami/kjudge-api-go/task"
)

// Error is a fs-specific error, occured when FS tasks failed
type Error struct {
	Stderr   string
	ExitCode int
}

func (f Error) Error() string {
	return fmt.Sprintf("FS Error (exitcode %d): %s", f.ExitCode, f.Stderr)
}

func runFsTask(T *task.Task) error {
	res := task.Enqueue(T)
	if res.ExitCode != 0 {
		return Error{res.Stderr, res.ExitCode}
	}
	return nil
}

// Copy creates a task that performs recursive copy.
func Copy(source, dest string) error {
	tsk := task.New("cp", []string{"-a", source, dest}, "")
	return runFsTask(tsk)
}

// Move creates a task that performs recursive move.
func Move(source, dest string) error {
	tsk := task.New("mv", []string{"-a", source, dest}, "")
	return runFsTask(tsk)
}

// Mkdir creates new folders according to the path specified
func Mkdir(source string) error {
	tsk := task.New("mkdir", []string{"-p", source, "-m", "777"}, "")
	return runFsTask(tsk)
}

// Chmod changes a file/folder's permission
func Chmod(source, mode string) error {
	tsk := task.New("chmod", []string{"-R", mode, source}, "")
	return runFsTask(tsk)
}

// Remove removes the specified folder/file
func Remove(source string) error {
	tsk := task.New("rm", []string{"-rf", source}, "")
	return runFsTask(tsk)
}
