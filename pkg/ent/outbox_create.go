// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/ent/outbox"
)

// OutboxCreate is the builder for creating a Outbox entity.
type OutboxCreate struct {
	config
	mutation *OutboxMutation
	hooks    []Hook
}

// SetAggregateType sets the "aggregate_type" field.
func (oc *OutboxCreate) SetAggregateType(s string) *OutboxCreate {
	oc.mutation.SetAggregateType(s)
	return oc
}

// SetAggregateID sets the "aggregate_id" field.
func (oc *OutboxCreate) SetAggregateID(s string) *OutboxCreate {
	oc.mutation.SetAggregateID(s)
	return oc
}

// SetEvent sets the "event" field.
func (oc *OutboxCreate) SetEvent(s string) *OutboxCreate {
	oc.mutation.SetEvent(s)
	return oc
}

// SetPayload sets the "payload" field.
func (oc *OutboxCreate) SetPayload(b []byte) *OutboxCreate {
	oc.mutation.SetPayload(b)
	return oc
}

// SetRetryAt sets the "retry_at" field.
func (oc *OutboxCreate) SetRetryAt(t time.Time) *OutboxCreate {
	oc.mutation.SetRetryAt(t)
	return oc
}

// SetRetryCount sets the "retry_count" field.
func (oc *OutboxCreate) SetRetryCount(i int) *OutboxCreate {
	oc.mutation.SetRetryCount(i)
	return oc
}

// SetNillableRetryCount sets the "retry_count" field if the given value is not nil.
func (oc *OutboxCreate) SetNillableRetryCount(i *int) *OutboxCreate {
	if i != nil {
		oc.SetRetryCount(*i)
	}
	return oc
}

// SetID sets the "id" field.
func (oc *OutboxCreate) SetID(i int64) *OutboxCreate {
	oc.mutation.SetID(i)
	return oc
}

// Mutation returns the OutboxMutation object of the builder.
func (oc *OutboxCreate) Mutation() *OutboxMutation {
	return oc.mutation
}

// Save creates the Outbox in the database.
func (oc *OutboxCreate) Save(ctx context.Context) (*Outbox, error) {
	return withHooks(ctx, oc.sqlSave, oc.mutation, oc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (oc *OutboxCreate) SaveX(ctx context.Context) *Outbox {
	v, err := oc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (oc *OutboxCreate) Exec(ctx context.Context) error {
	_, err := oc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (oc *OutboxCreate) ExecX(ctx context.Context) {
	if err := oc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (oc *OutboxCreate) check() error {
	if _, ok := oc.mutation.AggregateType(); !ok {
		return &ValidationError{Name: "aggregate_type", err: errors.New(`ent: missing required field "Outbox.aggregate_type"`)}
	}
	if _, ok := oc.mutation.AggregateID(); !ok {
		return &ValidationError{Name: "aggregate_id", err: errors.New(`ent: missing required field "Outbox.aggregate_id"`)}
	}
	if _, ok := oc.mutation.Event(); !ok {
		return &ValidationError{Name: "event", err: errors.New(`ent: missing required field "Outbox.event"`)}
	}
	if _, ok := oc.mutation.Payload(); !ok {
		return &ValidationError{Name: "payload", err: errors.New(`ent: missing required field "Outbox.payload"`)}
	}
	if _, ok := oc.mutation.RetryAt(); !ok {
		return &ValidationError{Name: "retry_at", err: errors.New(`ent: missing required field "Outbox.retry_at"`)}
	}
	return nil
}

func (oc *OutboxCreate) sqlSave(ctx context.Context) (*Outbox, error) {
	if err := oc.check(); err != nil {
		return nil, err
	}
	_node, _spec := oc.createSpec()
	if err := sqlgraph.CreateNode(ctx, oc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	oc.mutation.id = &_node.ID
	oc.mutation.done = true
	return _node, nil
}

func (oc *OutboxCreate) createSpec() (*Outbox, *sqlgraph.CreateSpec) {
	var (
		_node = &Outbox{config: oc.config}
		_spec = sqlgraph.NewCreateSpec(outbox.Table, sqlgraph.NewFieldSpec(outbox.FieldID, field.TypeInt64))
	)
	if id, ok := oc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := oc.mutation.AggregateType(); ok {
		_spec.SetField(outbox.FieldAggregateType, field.TypeString, value)
		_node.AggregateType = value
	}
	if value, ok := oc.mutation.AggregateID(); ok {
		_spec.SetField(outbox.FieldAggregateID, field.TypeString, value)
		_node.AggregateID = value
	}
	if value, ok := oc.mutation.Event(); ok {
		_spec.SetField(outbox.FieldEvent, field.TypeString, value)
		_node.Event = value
	}
	if value, ok := oc.mutation.Payload(); ok {
		_spec.SetField(outbox.FieldPayload, field.TypeBytes, value)
		_node.Payload = value
	}
	if value, ok := oc.mutation.RetryAt(); ok {
		_spec.SetField(outbox.FieldRetryAt, field.TypeTime, value)
		_node.RetryAt = &value
	}
	if value, ok := oc.mutation.RetryCount(); ok {
		_spec.SetField(outbox.FieldRetryCount, field.TypeInt, value)
		_node.RetryCount = value
	}
	return _node, _spec
}

// OutboxCreateBulk is the builder for creating many Outbox entities in bulk.
type OutboxCreateBulk struct {
	config
	err      error
	builders []*OutboxCreate
}

// Save creates the Outbox entities in the database.
func (ocb *OutboxCreateBulk) Save(ctx context.Context) ([]*Outbox, error) {
	if ocb.err != nil {
		return nil, ocb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ocb.builders))
	nodes := make([]*Outbox, len(ocb.builders))
	mutators := make([]Mutator, len(ocb.builders))
	for i := range ocb.builders {
		func(i int, root context.Context) {
			builder := ocb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*OutboxMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ocb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ocb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ocb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ocb *OutboxCreateBulk) SaveX(ctx context.Context) []*Outbox {
	v, err := ocb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ocb *OutboxCreateBulk) Exec(ctx context.Context) error {
	_, err := ocb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ocb *OutboxCreateBulk) ExecX(ctx context.Context) {
	if err := ocb.Exec(ctx); err != nil {
		panic(err)
	}
}
