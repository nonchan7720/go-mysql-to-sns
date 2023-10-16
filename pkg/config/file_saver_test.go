package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileSaverWithPosition(t *testing.T) {
	const filename = "binlog_position_test"
	saver := FileSaver{
		Name: filename,
	}
	defer func() {
		_ = os.Remove(filename)
	}()
	err := saver.Save(Position, BinlogSaveFormat{"binlog", 100, nil})
	require.Nil(t, err)
	format, err := saver.Load(Position)
	require.Nil(t, err)
	require.Equal(t, format.File, "binlog")
	require.Equal(t, format.Position, 100)
	require.Equal(t, format.GTID, []byte(nil))
}

func TestFileSaverWithGTID(t *testing.T) {
	const filename = "binlog_gtid_test"
	saver := FileSaver{
		Name: filename,
	}
	defer func() {
		_ = os.Remove(filename)
	}()
	err := saver.Save(GTID, BinlogSaveFormat{GTID: []byte("de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2")})
	require.Nil(t, err)
	format, err := saver.Load(GTID)
	require.Nil(t, err)
	require.Equal(t, format.File, "")
	require.Equal(t, format.Position, 0)
	require.Equal(t, format.GTID, []byte("de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"))
}
