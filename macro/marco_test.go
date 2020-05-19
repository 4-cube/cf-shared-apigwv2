package macro_test

import (
	"github.com/4-cube/cf-shared-apigwv2/macro"
	"github.com/4-cube/cf-shared-apigwv2/x"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestMacro_ProcessFragment(t *testing.T) {
	paths := x.TestFilePaths("*-input.json")
	for _, path := range paths {
		base := filepath.Base(path)
		t.Run(base, func(t *testing.T) {
			fragment, err := ioutil.ReadFile(path)
			assert.NoError(t, err)

			expected, err := ioutil.ReadFile(strings.Replace(path, "-input.json", "-output.json", 1))

			m := macro.NewMacro(fragment, logrus.New())
			output, err := m.ProcessFragment()
			assert.NoError(t, err)
			require.NotNil(t, output)
			assert.JSONEq(t, string(expected), string(output))
		})
	}
}
