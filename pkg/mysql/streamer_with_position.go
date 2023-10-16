package mysql

import (
	"context"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
)

type positionStreamer struct {
	*config.Config
	streamer *replication.BinlogStreamer
	syncer   *replication.BinlogSyncer
}

var (
	_ Streamer = (*positionStreamer)(nil)
)

func NewPositionStreamer(ctx context.Context, conf *config.Config) (Streamer, error) {
	return newPositionStreamer(ctx, conf)
}

func newPositionStreamer(ctx context.Context, conf *config.Config) (*gTIDStreamer, error) {
	var (
		file string
		pos  int
	)
	conn, err := conf.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if format, loadErr := conf.Saver.Load(config.Position); loadErr == nil {
		file = format.File
		pos = format.Position
	}
	if file == "" && pos == 0 {
		file, pos, err = loadBinlogPosition(conn)
		if err != nil {
			return nil, err
		}
	}

	serverId, err := findServerId(conn)
	if err != nil {
		return nil, err
	}
	syncer, err := conf.NewBinlogSyncer(serverId)
	if err != nil {
		return nil, err
	}

	streamer, err := syncer.StartSync(mysql.Position{Name: file, Pos: uint32(pos)})
	if err != nil {
		return nil, err
	}
	st := &gTIDStreamer{
		Config:   conf,
		streamer: streamer,
		syncer:   syncer,
	}
	return st, nil
}

func (st *positionStreamer) GetEvent(ctx context.Context) (*replication.BinlogEvent, error) {
	return st.streamer.GetEvent(ctx)
}

func (st *positionStreamer) Close() {
	st.syncer.Close()
}

func (st *positionStreamer) Save() error {
	pos := st.syncer.GetNextPosition()
	return st.Config.Saver.Save(config.Position, config.BinlogSaveFormat{
		File:     pos.Name,
		Position: int(pos.Pos),
	})
}
