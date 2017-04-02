package languages

import (
	"path"
	"path/filepath"

	"github.com/natsukagami/kjudge-api-go/lib/fs"
	"github.com/natsukagami/kjudge-api-go/task"
)

// cpp implements Interface for the GNU C++ Compiler
type cpp struct {
}

// Ext returns the source code's extension
func (c cpp) Ext() string {
	return ".cpp"
}

func (c cpp) exec() string {
	return "g++"
}

func (c cpp) defaultArgs() []string {
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
func (c cpp) Compile(name, folder string) error {
	tsk := task.New(c.exec(), append(c.defaultArgs(), name, name+c.Ext()), folder)
	return doCompile(tsk)
}

// CompileGrader compiles a *.cpp, along with the problem's grader file into an executable
func (c cpp) CompileGrader(name, folder, problemFolder string) error {
	fs.Copy(path.Join(problemFolder, "grader.cpp"), path.Join(folder, "grader.cpp"))
	fs.Copy(path.Join(problemFolder, "grader.h"), path.Join(folder, "grader.h"))
	tsk := task.New(c.exec(), append(c.defaultArgs(), name, name+c.Ext(), "grader.cpp"), folder)
	return doCompile(tsk)
}

// CompileComparator compile a problem's compare.cpp comparator into an executable.
func (c cpp) CompileComparator(problemFolder string) error {
	tsk := task.New(c.exec(), append(c.defaultArgs(), "compare", "compare"+c.Ext()), problemFolder)
	return doCompile(tsk)
}

// RunCommand returns the required command to run the executable.
func (c cpp) Executable(name string, folder string) string {
	return filepath.Join(folder, name)
}
