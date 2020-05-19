package x

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	ProjectRoot = filepath.Join(filepath.Dir(b), "../")
)

const TestTemplates = "test-templates"

func LoadTestFile(f string) []byte {
	fp := filepath.Join(ProjectRoot, TestTemplates, f)
	jf, _ := ioutil.ReadFile(fp)
	return jf
}

func TestFilePaths(pattern string) []string {
	paths, _ := filepath.Glob(filepath.Join(ProjectRoot, TestTemplates, pattern))
	return paths
}
