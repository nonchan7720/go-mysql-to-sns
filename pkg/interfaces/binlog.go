package interfaces

type Binlog interface {
	LoadBinlogFilePos() (file string, position int, err error)
	LoadBinlogMapEventFilePos() (file string, position int, err error)
	Save(file string, position int) error
}
