package mysql

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"sync"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/logging"
)

type gTIDStreamer struct {
	*config.Config
	streamer *replication.BinlogStreamer
	syncer   *replication.BinlogSyncer
	gtid     mysql.GTIDSet
	mu       sync.RWMutex
}

var (
	_ Streamer = (*gTIDStreamer)(nil)
)

func NewGTIDStreamer(ctx context.Context, conf *config.Config) (Streamer, error) {
	return newGTIDStreamer(ctx, conf)
}

func newGTIDStreamer(ctx context.Context, conf *config.Config) (*gTIDStreamer, error) {
	conn, err := conf.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	gtidSet, err := loadGTID(conn, conf)
	if err != nil {
		return nil, err
	}

	serverId, err := findServerId(conn)
	if err != nil {
		return nil, err
	}
	syncer, err := conf.NewBinlogSyncer(serverId,
		config.WithLogger(logging.NewBinlogLogger()),
		config.WithDisableRetrySync(), // リトライしない
	)
	if err != nil {
		return nil, err
	}
	streamer, err := syncer.StartSyncGTID(gtidSet)
	if err != nil {
		return nil, err
	}
	st := &gTIDStreamer{
		Config:   conf,
		streamer: streamer,
		syncer:   syncer,
		gtid:     gtidSet,
	}
	return st, nil
}

func (st *gTIDStreamer) GetEvent(ctx context.Context) (*replication.BinlogEvent, error) {
	ev, err := st.streamer.GetEvent(ctx)
	if err != nil {
		return nil, err
	}
	switch e := ev.Event.(type) {
	case *replication.RowsEvent:
		e.Dump(os.Stdout)
	case *replication.GTIDEvent:
		e.Dump(os.Stdout)
	case *replication.XIDEvent:
		if e.GSet != nil {
			st.updateGTID(e.GSet)
			slog.Info(e.GSet.String())
		}
	}
	return ev, nil
}

func (st *gTIDStreamer) Close() {
	st.syncer.Close()
}

func (st *gTIDStreamer) Save() error {
	return st.Config.Saver.Save(config.GTID, config.BinlogSaveFormat{
		GTID: st.gtid.Encode(),
	})
}

func (st *gTIDStreamer) updateGTID(gtid mysql.GTIDSet) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.gtid = gtid
}

func loadGTID(conn *sql.DB, conf *config.Config) (gtidSet mysql.GTIDSet, err error) {
	var gtid []byte

	if format, loadErr := conf.Saver.Load(config.GTID); loadErr == nil {
		gtid = format.GTID
	}
	if len(gtid) == 0 {
		v, err := loadGlobalBinlogGTID(conn)
		if err != nil {
			return nil, err
		}
		gtid = []byte(v)
	}

	gtidSet, err = mysql.DecodeMysqlGTIDSet(gtid)
	if err != nil {
		slog.Warn(err.Error())
		gtidSet, err = mysql.ParseMysqlGTIDSet(string(gtid))
		if err != nil {
			return nil, err
		}
	}
	return
}
