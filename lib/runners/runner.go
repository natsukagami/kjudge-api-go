// Package runners provides an implementation of sandbox wrappers.
package runners

import "fmt"

// Result represents the outcome of an execution.
type Result struct {
	Time   int64
	Mem    int64
	Status string
}

// Error represents a sandbox Error
type Error struct {
	Err      string
	Exitcode int
}

func (e Error) Error() string {
	return fmt.Sprintf("Sandbox error (exitcode %d): %s", e.Exitcode, e.Err)
}

// Interface presents a sandbox implementation
type Interface interface {
	Prepare() error
	Run(cmd, cwd string, time, mem int64) (*Result, error)
	Cleanup() error
}

// SandboxError returns the sandbox's signal meaning
func SandboxError(signal string) string {
	switch signal {
	case "RE":
		return "Runtime Error"
	case "SG":
		return "Killed by Signal (Memory Limit Exceeded?)"
	case "TO":
		return "Time Limit Exceeded"
	case "MLE":
		return "Memory Limit Exceeded"
	default:
		return "Unknown Error " + signal
	}
}
