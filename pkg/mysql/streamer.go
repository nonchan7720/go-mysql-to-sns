package mysql

import (
	"context"

	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
)

type Streamer interface {
	GetEvent(ctx context.Context) (*replication.BinlogEvent, error)
	Save() error
	Close()
}

func NewStreamer(ctx context.Context, formatType config.BinlogSaveFormatType, conf *config.Config) (st Streamer, err error) {
	if formatType == config.GTID {
		st, err = NewGTIDStreamer(ctx, conf)
	} else {
		st, err = NewPositionStreamer(ctx, conf)
	}
	if err != nil {
		return
	}
	// 読み込んだ現地点を保存する
	if err = st.Save(); err != nil {
		return
	}
	return
}
