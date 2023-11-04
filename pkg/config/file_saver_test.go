package config

import (
	"os"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func TestValidateFileSaver(t *testing.T) {
	saver := FileSaver{
		Name: "",
	}
	err := validation.Validate(&saver)
	require.Equal(t, "Name: cannot be blank.", err.Error())
}
