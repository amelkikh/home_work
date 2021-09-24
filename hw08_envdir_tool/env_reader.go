package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrDirNotExists = errors.New("no such directory")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, ErrDirNotExists
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.Contains(f.Name(), "=") {
			continue
		}
		fPath := strings.Join([]string{dir, f.Name()}, string(os.PathSeparator))

		if env, err := parseFile(fPath); err == nil {
			envs[f.Name()] = env
		}
	}

	return envs, nil
}

func parseFile(path string) (EnvValue, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return EnvValue{}, err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return EnvValue{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err != nil {
		return EnvValue{}, err
	}

	var line []byte
	for scanner.Scan() {
		line = bytes.ReplaceAll(scanner.Bytes(), []byte{0x00}, []byte("\n"))
		break
	}
	data := strings.TrimRightFunc(string(line), func(r rune) bool {
		return r == ' ' || r == '\t'
	})

	return EnvValue{
		Value:      data,
		NeedRemove: stat.Size() == 0,
	}, nil
}
