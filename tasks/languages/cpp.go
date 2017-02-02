package languages

import (
	"path"

	"github.com/natsukagami/kjudge-api-go/task"
	"github.com/natsukagami/kjudge-api-go/tasks/fs"
)

// Cpp implements Interface for the GNU C++ Compiler
type Cpp struct {
}

// Ext returns the source code's extension
func (c Cpp) Ext() string {
	return ".cpp"
}

func (c Cpp) exec() string {
	return "g++"
}

func (c Cpp) defaultArgs() []string {
	return []string{
		"-O2",
		"-s",
		"-std=c++11",
		"-static",
		"-lm",
		"-o",
	}
}

// Compile compiles a *.cpp file into an executable
func (c Cpp) Compile(name, folder string) error {
	tsk := task.NewTask(c.exec(), append(c.defaultArgs(), name, name+c.Ext()), folder)
	return doCompile(&tsk)
}

// CompileGrader compiles a *.cpp, along with the problem's grader file into an executable
func (c Cpp) CompileGrader(name, folder, problemFolder string) error {
	fs.Copy(path.Join(problemFolder, "grader.cpp"), path.Join(folder, "grader.cpp"))
	fs.Copy(path.Join(problemFolder, "grader.h"), path.Join(folder, "grader.h"))
	tsk := task.NewTask(c.exec(), append(c.defaultArgs(), name, name+c.Ext(), "grader.cpp"), folder)
	return doCompile(&tsk)
}

// CompileComparator compile a problem's compare.cpp comparator into an executable.
func (c Cpp) CompileComparator(problemFolder string) error {
	tsk := task.NewTask(c.exec(), append(c.defaultArgs(), "compare", "compare"+c.Ext()), problemFolder)
	return doCompile(&tsk)
}
