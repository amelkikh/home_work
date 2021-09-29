package main

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	path, err := os.Getwd()
	require.NoError(t, err)
	dataPath := strings.Join([]string{path, "testdata", "env"}, string(os.PathSeparator))
	expectedEnvs := Environment{
		"BAR": {
			Value:      "bar",
			NeedRemove: false,
		},
		"EMPTY": {
			Value:      "",
			NeedRemove: false,
		},
		"FOO": {
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		},
		"HELLO": {
			Value:      "\"hello\"",
			NeedRemove: false,
		},
		"UNSET": {
			Value:      "",
			NeedRemove: true,
		},
		"ZERO_BYTE": {
			Value:      "ABC\nDEF",
			NeedRemove: false,
		},
	}
	envs, err := ReadDir(dataPath)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(envs, expectedEnvs))
}
