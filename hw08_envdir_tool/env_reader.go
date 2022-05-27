package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirR, errD := os.ReadDir(dir)
	if errD != nil {
		return nil, errD
	}

	env := make(Environment, len(dirR))

	for _, file := range dirR {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "=") {
			continue
		}

		openFile, errF := os.Open(filepath.Join(dir, file.Name()))
		if errF != nil {
			return nil, errF
		}
		defer func() {
			_ = openFile.Close()
		}()

		env[file.Name()] = EnvValue{Value: "", NeedRemove: false}
		buf := bufio.NewReader(openFile)
		line, _, errL := buf.ReadLine()
		if errL != nil {
			if errL == io.EOF {
				env[file.Name()] = EnvValue{NeedRemove: true}

				continue
			}

			return nil, errL
		}

		str := bytes.TrimRight(line, " \t")
		str = bytes.ReplaceAll(str, []byte("\x00"), []byte("\n"))

		if len(str) == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}

			continue
		}

		env[file.Name()] = EnvValue{Value: string(str)}
	}

	return env, nil
}
