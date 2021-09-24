package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	osEnv := map[string]string{
		"ADDED": "from original env",
		"UNSET": "SHOULD_REMOVE",
		"HELLO": "SHOULD_REPLACE",
		"FOO":   "SHOULD_REPLACE",
		"EMPTY": "SHOULD_BE_EMPTY",
	}
	for k, v := range osEnv {
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}

	path, err := os.Getwd()
	require.NoError(t, err)

	dataPath := strings.Join([]string{path, "testdata", "env"}, string(os.PathSeparator))
	envs, err := ReadDir(dataPath)
	require.NoError(t, err)

	envCmd := exec.Command("sh", "-c", "echo \"HELLO is ($HELLO)\nBAR is (${BAR})\nFOO is (${FOO})\nUNSET is (${UNSET})\nADDED is (${ADDED})\nEMPTY is (${EMPTY})\nZERO_BYTE is (${ZERO_BYTE})\nTEST_KEY= is (${TEST_KEY=})\"")

	err = applyEnvs(envs)
	require.NoError(t, err)

	envCmd.Env = os.Environ()
	var b2 bytes.Buffer
	envCmd.Stdout = &b2

	err = envCmd.Start()
	require.NoError(t, err)
	err = envCmd.Wait()
	require.NoError(t, err)

	expectedOut := []byte(`HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
ZERO_BYTE is (ABC
DEF)
TEST_KEY= is ()
`)
	t.Log(b2.Bytes())
	t.Log(expectedOut)

	t.Log(b2.String())
	t.Log(string(expectedOut))
	require.True(t, bytes.Equal(b2.Bytes(), expectedOut))
}
