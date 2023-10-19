package mysql

import (
	"context"
	"log/slog"
	"math"
	"sort"
	"sync"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/config"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/ent"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/ent/outbox"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/interfaces"
	"github.com/nonchan7720/go-mysql-to-sns/pkg/mysql/client"
	"golang.org/x/sync/errgroup"
)

type transaction func(ctx context.Context, db *ent.Client, execFunc func(ctx context.Context, tx *ent.Tx) error) (err error)

type OutboxPolling struct {
	*config.OutboxPolling
	client   *client.Client
	stop     sync.Once
	stopping chan struct{}

	timeNow          func() time.Time
	publisher        interfaces.Publisher
	runInTransaction transaction
}

func NewOutboxPolling(
	polling *config.OutboxPolling,
	publisher interfaces.Publisher,
	runInTransaction transaction,
) (*OutboxPolling, error) {
	c, err := client.NewDB(&polling.Config)
	if err != nil {
		return nil, err
	}
	return &OutboxPolling{
		OutboxPolling: polling,
		client:        c,
		timeNow:       time.Now,
		publisher:     publisher,
		stop:          sync.Once{},
		stopping:      make(chan struct{}),
	}, nil
}

func (p *OutboxPolling) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		<-p.stopping
		cancel()
		return nil
	})

	group.Go(func() error {
		return p.processing(ctx)
	})

	return group.Wait()
}

func (p *OutboxPolling) processing(ctx context.Context) error {
	pollingTimer := time.NewTicker(p.OutboxConfig.PollingInterval)
	for {
		events, err := p.findEvents(ctx)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
		}
		if len(events) > 0 {
			p.processingEvent(ctx, events)
		}

		select {
		case <-pollingTimer.C:
			continue
		case <-ctx.Done():
			pollingTimer.Stop()
			return ctx.Err()
		}
	}
}

func (p *OutboxPolling) processingEvent(ctx context.Context, events []*outboxPollingEvent) {
	mp := eventToMapProducerEvent(p.Publisher.FindProducer, events)
	for producer := range mp {
		event := mp[producer]
		mpGroupEvent := eventToGroupingAggregateId(event)
		for mGroupId := range mpGroupEvent {
			groupEvents := mpGroupEvent[mGroupId]
			if mGroupId != "" {
				eventsSort(groupEvents)
			}
			for idx := range groupEvents {
				event := groupEvents[idx]
				isOk := p.publishEvent(ctx, producer, event)
				if mGroupId != "" && !isOk {
					// メッセージグループがありかつ失敗の場合はこれ以上送らない
					break
				}
			}
		}
	}
}

func (p *OutboxPolling) findEvents(ctx context.Context) ([]*outboxPollingEvent, error) {
	q := p.client.Outbox.
		Query().
		Where(outbox.Or(outbox.RetryCountIsNil(), outbox.RetryCountLT(p.OutboxConfig.MaxRetryCount))).
		Order(outbox.ByID(entsql.OrderAsc())).Limit(1000)
	events, err := q.All(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]*outboxPollingEvent, len(events))
	for idx := range events {
		event := events[idx]
		results[idx] = newOutboxEvent(event)
	}
	return results, nil
}

func eventToMapProducerEvent(producerFindFunc func(string) (string, error), events []*outboxPollingEvent) map[string][]*outboxPollingEvent {
	mp := map[string][]*outboxPollingEvent{}
	for _, event := range events {
		producer, err := producerFindFunc(event.AggregateType)
		if err != nil {
			slog.
				With(slog.String("AggregateType", event.AggregateType)).
				Error("Producer information matching the AggregateType could not be obtained.")
			continue
		}
		mp[producer] = append(mp[producer], event)
	}
	return mp
}

func eventToGroupingAggregateId(events []*outboxPollingEvent) map[string][]*outboxPollingEvent {
	mp := map[string][]*outboxPollingEvent{}
	for _, event := range events {
		mp[event.AggregateID] = append(mp[event.AggregateID], event)
	}
	return mp
}

func eventsSort(events []*outboxPollingEvent) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
}

func (p *OutboxPolling) publishEvent(ctx context.Context, producer string, event *outboxPollingEvent) (isOk bool) {
	idField := slog.Int64("ID", event.ID)
	if event.CanNotRetry(p.timeNow()) {
		slog.With(idField).Warn("waiting for the retry")
		return false
	}

	defer func() {
		fn := p.txUpdateOrDelete(isOk)
		err := p.runInTransaction(ctx, p.client.Client, func(ctx context.Context, tx *ent.Tx) error {
			return fn(ctx, tx, event)
		})
		if err != nil {
			isOk = false
			slog.With(idField).Error(err.Error())
		}
	}()

	err := p.publisher.PublishOutbox(ctx, producer, interfaces.Outbox{
		AggregateId: event.AggregateID,
		EventType:   event.Event,
		Payload:     event.Payload,
	})
	return err == nil
}

func (p *OutboxPolling) txDeleteRecord(ctx context.Context, tx *ent.Tx, event *outboxPollingEvent) error {
	if err := tx.Outbox.DeleteOneID(event.ID).Exec(ctx); err != nil {
		return err
	}
	slog.With(slog.Int64("ID", event.ID)).Info("Published and deleted.")
	return nil
}

func (p *OutboxPolling) txUpdateRetryField(ctx context.Context, tx *ent.Tx, event *outboxPollingEvent) error {
	event.IncrementToRetryCount()
	if event.CheckMaxRetryCount(p.OutboxConfig.MaxRetryCount) {
		slog.With(
			slog.Int64("ID", event.ID),
			slog.Int("RetryCount", event.RetryCount),
		).Error("RetryCount has reached its maximum value.")
	}
	duration := getNextRetryDuration(p.OutboxConfig.RetryBackOff, event.RetryCount)
	nextRetryTime := p.timeNow().Add(duration)
	event.RetryAt = &nextRetryTime
	q := tx.Outbox.UpdateOneID(event.ID)
	q = q.SetRetryCount(event.RetryCount)
	q = q.SetRetryAt(*event.RetryAt)
	return q.Exec(ctx)
}

func (p *OutboxPolling) txUpdateOrDelete(isOk bool) func(ctx context.Context, tx *ent.Tx, event *outboxPollingEvent) error {
	if !isOk {
		return p.txUpdateRetryField
	} else {
		return p.txDeleteRecord
	}
}

// getNextRetryDuration return the next retry duration in seconds based on the attempt,
// the formula: `backoff * 2 ^ attempt`
func getNextRetryDuration(backoff time.Duration, attempt int) time.Duration {
	return time.Duration(backoff.Seconds()*math.Pow(2, float64(attempt))) * time.Second
}
