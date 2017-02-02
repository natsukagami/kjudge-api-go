package runners

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/natsukagami/kjudge-api-go/task"
	"github.com/natsukagami/kjudge-api-go/task/queue"
)

// This file implements a Runner struct using the isolate command.

var boxID = make(chan int)

// Isolate is an implemention of isolate sandbox wrapper.
type Isolate struct {
	id int
}

// Dir returns the directory of the sandbox
func (i *Isolate) Dir() string {
	return fmt.Sprintf("/var/lib/isolate/%d/box", i.id)
}

// Prepare prepares the sandbox for usage
func (i *Isolate) Prepare() error {
	tsk := task.NewTask("isolate", []string{
		"--init",
		"--cg",
		"-b", string(i.id),
	}, "")
	res := queue.Enqueue(&tsk)
	if res.ExitCode != 0 {
		go func() { // Must be done async or the code will block
			boxID <- i.id
		}()
		return Error{res.Stderr, res.ExitCode}
	}
	return nil
}

// Run runs the command under the sandbox's restrictions
func (i *Isolate) Run(cmd, cwd string, time, mem int64) (r Result, e error) {
	tsk := task.NewTask("isolate", []string{
		"-b", strconv.Itoa(i.id),
		"--cg",
		"--run",
		"--dir=" + path.Join(i.Dir(), "env") + "=" + cwd + ":rw",
		"-t", fmt.Sprintf("%.3f", float64(time)/1000),
		"-w", fmt.Sprintf("%.3f", float64(time)/1000+0.5),
		"-m", fmt.Sprintf("%d", mem),
		"-i", path.Join(i.Dir(), "env", "input.txt"),
		"-o", path.Join(i.Dir(), "env", "output.txt"),
		"-M", path.Join(cwd, "meta.txt"),
	}, "")
	queue.Enqueue(&tsk)
	// Opens the meta file
	dat, err := ioutil.ReadFile(path.Join(cwd, "meta.txt"))
	if err != nil {
		e = err
		return
	}
	meta := strings.Split(string(dat), "\n")
	var timeWall int64 = -1
	r.Status = "OK"
	for _, str := range meta {
		line := strings.Split(str, ":")
		switch line[0] {
		case "time-wall":
			t, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				e = err
				return
			}
			timeWall = int64(t * 1000)
		case "time":
			t, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				e = err
				return
			}
			r.Time = int64(t * 1000)
		case "cg-mem":
			t, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				e = err
				return
			}
			r.Mem = int64(t)
		case "status":
			r.Status = line[1]
		}
		if timeWall >= 0 {
			r.Time = timeWall
		}
	}
	return
}

// Cleanup cleans the sandbox and make it ready for another use
func (i *Isolate) Cleanup() error {
	tsk := task.NewTask("isolate", []string{"--cleanup", "-b", strconv.Itoa(i.id)}, "")
	res := queue.Enqueue(&tsk)
	go func() {
		boxID <- i.id
	}()
	if res.ExitCode != 0 {
		return Error{res.Stderr, res.ExitCode}
	}
	return nil
}

// Make creates a new isolate sandbox
func Make() *Isolate {
	box := Isolate{<-boxID}
	return &box
}
