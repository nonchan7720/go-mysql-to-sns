package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileSaver(t *testing.T) {
	const filename = "binlog_test"
	saver := FileSaver{
		Name: filename,
	}
	defer func() {
		_ = os.Remove(filename)
	}()
	err := saver.Save("binlog", 100)
	require.Nil(t, err)
	file, pos, err := saver.Load()
	require.Nil(t, err)
	require.Equal(t, file, "binlog")
	require.Equal(t, pos, 100)
}
