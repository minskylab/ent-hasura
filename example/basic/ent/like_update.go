// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/minskylab/ent-hasura/example/basic/ent/like"
	"github.com/minskylab/ent-hasura/example/basic/ent/predicate"
	"github.com/minskylab/ent-hasura/example/basic/ent/user"
)

// LikeUpdate is the builder for updating Like entities.
type LikeUpdate struct {
	config
	hooks    []Hook
	mutation *LikeMutation
}

// Where appends a list predicates to the LikeUpdate builder.
func (lu *LikeUpdate) Where(ps ...predicate.Like) *LikeUpdate {
	lu.mutation.Where(ps...)
	return lu
}

// SetCreatedAt sets the "created_at" field.
func (lu *LikeUpdate) SetCreatedAt(t time.Time) *LikeUpdate {
	lu.mutation.SetCreatedAt(t)
	return lu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lu *LikeUpdate) SetNillableCreatedAt(t *time.Time) *LikeUpdate {
	if t != nil {
		lu.SetCreatedAt(*t)
	}
	return lu
}

// SetCreatorID sets the "creator" edge to the User entity by ID.
func (lu *LikeUpdate) SetCreatorID(id int) *LikeUpdate {
	lu.mutation.SetCreatorID(id)
	return lu
}

// SetCreator sets the "creator" edge to the User entity.
func (lu *LikeUpdate) SetCreator(u *User) *LikeUpdate {
	return lu.SetCreatorID(u.ID)
}

// Mutation returns the LikeMutation object of the builder.
func (lu *LikeUpdate) Mutation() *LikeMutation {
	return lu.mutation
}

// ClearCreator clears the "creator" edge to the User entity.
func (lu *LikeUpdate) ClearCreator() *LikeUpdate {
	lu.mutation.ClearCreator()
	return lu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lu *LikeUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(lu.hooks) == 0 {
		if err = lu.check(); err != nil {
			return 0, err
		}
		affected, err = lu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LikeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = lu.check(); err != nil {
				return 0, err
			}
			lu.mutation = mutation
			affected, err = lu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(lu.hooks) - 1; i >= 0; i-- {
			if lu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = lu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, lu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (lu *LikeUpdate) SaveX(ctx context.Context) int {
	affected, err := lu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lu *LikeUpdate) Exec(ctx context.Context) error {
	_, err := lu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lu *LikeUpdate) ExecX(ctx context.Context) {
	if err := lu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lu *LikeUpdate) check() error {
	if _, ok := lu.mutation.CreatorID(); lu.mutation.CreatorCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"creator\"")
	}
	return nil
}

func (lu *LikeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   like.Table,
			Columns: like.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: like.FieldID,
			},
		},
	}
	if ps := lu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: like.FieldCreatedAt,
		})
	}
	if lu.mutation.CreatorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   like.CreatorTable,
			Columns: []string{like.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := lu.mutation.CreatorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   like.CreatorTable,
			Columns: []string{like.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, lu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{like.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// LikeUpdateOne is the builder for updating a single Like entity.
type LikeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *LikeMutation
}

// SetCreatedAt sets the "created_at" field.
func (luo *LikeUpdateOne) SetCreatedAt(t time.Time) *LikeUpdateOne {
	luo.mutation.SetCreatedAt(t)
	return luo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (luo *LikeUpdateOne) SetNillableCreatedAt(t *time.Time) *LikeUpdateOne {
	if t != nil {
		luo.SetCreatedAt(*t)
	}
	return luo
}

// SetCreatorID sets the "creator" edge to the User entity by ID.
func (luo *LikeUpdateOne) SetCreatorID(id int) *LikeUpdateOne {
	luo.mutation.SetCreatorID(id)
	return luo
}

// SetCreator sets the "creator" edge to the User entity.
func (luo *LikeUpdateOne) SetCreator(u *User) *LikeUpdateOne {
	return luo.SetCreatorID(u.ID)
}

// Mutation returns the LikeMutation object of the builder.
func (luo *LikeUpdateOne) Mutation() *LikeMutation {
	return luo.mutation
}

// ClearCreator clears the "creator" edge to the User entity.
func (luo *LikeUpdateOne) ClearCreator() *LikeUpdateOne {
	luo.mutation.ClearCreator()
	return luo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (luo *LikeUpdateOne) Select(field string, fields ...string) *LikeUpdateOne {
	luo.fields = append([]string{field}, fields...)
	return luo
}

// Save executes the query and returns the updated Like entity.
func (luo *LikeUpdateOne) Save(ctx context.Context) (*Like, error) {
	var (
		err  error
		node *Like
	)
	if len(luo.hooks) == 0 {
		if err = luo.check(); err != nil {
			return nil, err
		}
		node, err = luo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*LikeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = luo.check(); err != nil {
				return nil, err
			}
			luo.mutation = mutation
			node, err = luo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(luo.hooks) - 1; i >= 0; i-- {
			if luo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = luo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, luo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (luo *LikeUpdateOne) SaveX(ctx context.Context) *Like {
	node, err := luo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (luo *LikeUpdateOne) Exec(ctx context.Context) error {
	_, err := luo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (luo *LikeUpdateOne) ExecX(ctx context.Context) {
	if err := luo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (luo *LikeUpdateOne) check() error {
	if _, ok := luo.mutation.CreatorID(); luo.mutation.CreatorCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"creator\"")
	}
	return nil
}

func (luo *LikeUpdateOne) sqlSave(ctx context.Context) (_node *Like, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   like.Table,
			Columns: like.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: like.FieldID,
			},
		},
	}
	id, ok := luo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Like.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := luo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, like.FieldID)
		for _, f := range fields {
			if !like.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != like.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := luo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := luo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: like.FieldCreatedAt,
		})
	}
	if luo.mutation.CreatorCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   like.CreatorTable,
			Columns: []string{like.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := luo.mutation.CreatorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   like.CreatorTable,
			Columns: []string{like.CreatorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Like{config: luo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, luo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{like.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
