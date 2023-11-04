package config

import (
	"bytes"
	"io"
	"os"
)

func NewExpandEnv(name string) (io.Reader, error) {
	buf, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewExpandEnvWithReader(bytes.NewBufferString(os.ExpandEnv(string(buf)))), nil
}

func NewExpandEnvWithReader(f io.Reader) io.Reader {
	return f
}
