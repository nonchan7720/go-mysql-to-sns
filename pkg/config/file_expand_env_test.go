package config

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type textScanner struct {
	io.Reader
}

func (t *textScanner) Close() error {
	return nil
}

func TestExpandEnvReader(t *testing.T) {
	data := `key: ${TEST}`
	t.Setenv("TEST", "Value")
	var value = struct {
		Key string `yaml:"key"`
	}{}
	scanner := &textScanner{Reader: strings.NewReader(data)}
	f := NewExpandEnvWithReadeCloser(scanner)
	err := yaml.NewDecoder(f).Decode(&value)
	require.NoError(t, err)
	require.Equal(t, "Value", value.Key)
}
