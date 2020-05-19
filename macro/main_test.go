package macro

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestMacro_ProcessFragment(t *testing.T) {
	paths, err := filepath.Glob("test-templates/*-input.json")
	assert.NoError(t, err)

	for _, path := range paths {
		base := filepath.Base(path)
		t.Run(base, func(t *testing.T) {
			fragment, err := ioutil.ReadFile(path)
			assert.NoError(t, err)

			expected, err := ioutil.ReadFile(strings.Replace(path, "-input.json", "-output.json", 1))

			macro := NewMacro(fragment, logrus.New())
			output, err := macro.ProcessFragment()
			assert.NoError(t, err)
			require.NotNil(t, output)
			assert.JSONEq(t, string(expected), string(output))
		})
	}
}
