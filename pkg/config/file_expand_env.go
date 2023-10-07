package config

import (
	"io"
	"os"
)

type ExpandEnv struct {
	io.ReadCloser
}

func (f *ExpandEnv) Read(p []byte) (n int, err error) {
	n, err = f.ReadCloser.Read(p)
	if err == nil {
		expandedData := os.ExpandEnv(string(p[:n]))
		copy(p, []byte(expandedData))
		n = len(expandedData)
	}
	return
}

func NewExpandEnv(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return NewExpandEnvWithReadeCloser(f), nil
}

func NewExpandEnvWithReadeCloser(f io.ReadCloser) io.ReadCloser {
	return &ExpandEnv{ReadCloser: f}
}
