// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/unbindapp/unbind-api/ent/jwtkey"
	"github.com/unbindapp/unbind-api/ent/predicate"
)

// JWTKeyQuery is the builder for querying JWTKey entities.
type JWTKeyQuery struct {
	config
	ctx        *QueryContext
	order      []jwtkey.OrderOption
	inters     []Interceptor
	predicates []predicate.JWTKey
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the JWTKeyQuery builder.
func (jkq *JWTKeyQuery) Where(ps ...predicate.JWTKey) *JWTKeyQuery {
	jkq.predicates = append(jkq.predicates, ps...)
	return jkq
}

// Limit the number of records to be returned by this query.
func (jkq *JWTKeyQuery) Limit(limit int) *JWTKeyQuery {
	jkq.ctx.Limit = &limit
	return jkq
}

// Offset to start from.
func (jkq *JWTKeyQuery) Offset(offset int) *JWTKeyQuery {
	jkq.ctx.Offset = &offset
	return jkq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (jkq *JWTKeyQuery) Unique(unique bool) *JWTKeyQuery {
	jkq.ctx.Unique = &unique
	return jkq
}

// Order specifies how the records should be ordered.
func (jkq *JWTKeyQuery) Order(o ...jwtkey.OrderOption) *JWTKeyQuery {
	jkq.order = append(jkq.order, o...)
	return jkq
}

// First returns the first JWTKey entity from the query.
// Returns a *NotFoundError when no JWTKey was found.
func (jkq *JWTKeyQuery) First(ctx context.Context) (*JWTKey, error) {
	nodes, err := jkq.Limit(1).All(setContextOp(ctx, jkq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{jwtkey.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (jkq *JWTKeyQuery) FirstX(ctx context.Context) *JWTKey {
	node, err := jkq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first JWTKey ID from the query.
// Returns a *NotFoundError when no JWTKey ID was found.
func (jkq *JWTKeyQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = jkq.Limit(1).IDs(setContextOp(ctx, jkq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{jwtkey.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (jkq *JWTKeyQuery) FirstIDX(ctx context.Context) int {
	id, err := jkq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single JWTKey entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one JWTKey entity is found.
// Returns a *NotFoundError when no JWTKey entities are found.
func (jkq *JWTKeyQuery) Only(ctx context.Context) (*JWTKey, error) {
	nodes, err := jkq.Limit(2).All(setContextOp(ctx, jkq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{jwtkey.Label}
	default:
		return nil, &NotSingularError{jwtkey.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (jkq *JWTKeyQuery) OnlyX(ctx context.Context) *JWTKey {
	node, err := jkq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only JWTKey ID in the query.
// Returns a *NotSingularError when more than one JWTKey ID is found.
// Returns a *NotFoundError when no entities are found.
func (jkq *JWTKeyQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = jkq.Limit(2).IDs(setContextOp(ctx, jkq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{jwtkey.Label}
	default:
		err = &NotSingularError{jwtkey.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (jkq *JWTKeyQuery) OnlyIDX(ctx context.Context) int {
	id, err := jkq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of JWTKeys.
func (jkq *JWTKeyQuery) All(ctx context.Context) ([]*JWTKey, error) {
	ctx = setContextOp(ctx, jkq.ctx, ent.OpQueryAll)
	if err := jkq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*JWTKey, *JWTKeyQuery]()
	return withInterceptors[[]*JWTKey](ctx, jkq, qr, jkq.inters)
}

// AllX is like All, but panics if an error occurs.
func (jkq *JWTKeyQuery) AllX(ctx context.Context) []*JWTKey {
	nodes, err := jkq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of JWTKey IDs.
func (jkq *JWTKeyQuery) IDs(ctx context.Context) (ids []int, err error) {
	if jkq.ctx.Unique == nil && jkq.path != nil {
		jkq.Unique(true)
	}
	ctx = setContextOp(ctx, jkq.ctx, ent.OpQueryIDs)
	if err = jkq.Select(jwtkey.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (jkq *JWTKeyQuery) IDsX(ctx context.Context) []int {
	ids, err := jkq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (jkq *JWTKeyQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, jkq.ctx, ent.OpQueryCount)
	if err := jkq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, jkq, querierCount[*JWTKeyQuery](), jkq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (jkq *JWTKeyQuery) CountX(ctx context.Context) int {
	count, err := jkq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (jkq *JWTKeyQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, jkq.ctx, ent.OpQueryExist)
	switch _, err := jkq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (jkq *JWTKeyQuery) ExistX(ctx context.Context) bool {
	exist, err := jkq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the JWTKeyQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (jkq *JWTKeyQuery) Clone() *JWTKeyQuery {
	if jkq == nil {
		return nil
	}
	return &JWTKeyQuery{
		config:     jkq.config,
		ctx:        jkq.ctx.Clone(),
		order:      append([]jwtkey.OrderOption{}, jkq.order...),
		inters:     append([]Interceptor{}, jkq.inters...),
		predicates: append([]predicate.JWTKey{}, jkq.predicates...),
		// clone intermediate query.
		sql:       jkq.sql.Clone(),
		path:      jkq.path,
		modifiers: append([]func(*sql.Selector){}, jkq.modifiers...),
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Label string `json:"label,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.JWTKey.Query().
//		GroupBy(jwtkey.FieldLabel).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (jkq *JWTKeyQuery) GroupBy(field string, fields ...string) *JWTKeyGroupBy {
	jkq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &JWTKeyGroupBy{build: jkq}
	grbuild.flds = &jkq.ctx.Fields
	grbuild.label = jwtkey.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Label string `json:"label,omitempty"`
//	}
//
//	client.JWTKey.Query().
//		Select(jwtkey.FieldLabel).
//		Scan(ctx, &v)
func (jkq *JWTKeyQuery) Select(fields ...string) *JWTKeySelect {
	jkq.ctx.Fields = append(jkq.ctx.Fields, fields...)
	sbuild := &JWTKeySelect{JWTKeyQuery: jkq}
	sbuild.label = jwtkey.Label
	sbuild.flds, sbuild.scan = &jkq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a JWTKeySelect configured with the given aggregations.
func (jkq *JWTKeyQuery) Aggregate(fns ...AggregateFunc) *JWTKeySelect {
	return jkq.Select().Aggregate(fns...)
}

func (jkq *JWTKeyQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range jkq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, jkq); err != nil {
				return err
			}
		}
	}
	for _, f := range jkq.ctx.Fields {
		if !jwtkey.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if jkq.path != nil {
		prev, err := jkq.path(ctx)
		if err != nil {
			return err
		}
		jkq.sql = prev
	}
	return nil
}

func (jkq *JWTKeyQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*JWTKey, error) {
	var (
		nodes = []*JWTKey{}
		_spec = jkq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*JWTKey).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &JWTKey{config: jkq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(jkq.modifiers) > 0 {
		_spec.Modifiers = jkq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, jkq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (jkq *JWTKeyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := jkq.querySpec()
	if len(jkq.modifiers) > 0 {
		_spec.Modifiers = jkq.modifiers
	}
	_spec.Node.Columns = jkq.ctx.Fields
	if len(jkq.ctx.Fields) > 0 {
		_spec.Unique = jkq.ctx.Unique != nil && *jkq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, jkq.driver, _spec)
}

func (jkq *JWTKeyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(jwtkey.Table, jwtkey.Columns, sqlgraph.NewFieldSpec(jwtkey.FieldID, field.TypeInt))
	_spec.From = jkq.sql
	if unique := jkq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if jkq.path != nil {
		_spec.Unique = true
	}
	if fields := jkq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, jwtkey.FieldID)
		for i := range fields {
			if fields[i] != jwtkey.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := jkq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := jkq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := jkq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := jkq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (jkq *JWTKeyQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(jkq.driver.Dialect())
	t1 := builder.Table(jwtkey.Table)
	columns := jkq.ctx.Fields
	if len(columns) == 0 {
		columns = jwtkey.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if jkq.sql != nil {
		selector = jkq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if jkq.ctx.Unique != nil && *jkq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range jkq.modifiers {
		m(selector)
	}
	for _, p := range jkq.predicates {
		p(selector)
	}
	for _, p := range jkq.order {
		p(selector)
	}
	if offset := jkq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := jkq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (jkq *JWTKeyQuery) Modify(modifiers ...func(s *sql.Selector)) *JWTKeySelect {
	jkq.modifiers = append(jkq.modifiers, modifiers...)
	return jkq.Select()
}

// JWTKeyGroupBy is the group-by builder for JWTKey entities.
type JWTKeyGroupBy struct {
	selector
	build *JWTKeyQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (jkgb *JWTKeyGroupBy) Aggregate(fns ...AggregateFunc) *JWTKeyGroupBy {
	jkgb.fns = append(jkgb.fns, fns...)
	return jkgb
}

// Scan applies the selector query and scans the result into the given value.
func (jkgb *JWTKeyGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, jkgb.build.ctx, ent.OpQueryGroupBy)
	if err := jkgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*JWTKeyQuery, *JWTKeyGroupBy](ctx, jkgb.build, jkgb, jkgb.build.inters, v)
}

func (jkgb *JWTKeyGroupBy) sqlScan(ctx context.Context, root *JWTKeyQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(jkgb.fns))
	for _, fn := range jkgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*jkgb.flds)+len(jkgb.fns))
		for _, f := range *jkgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*jkgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := jkgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// JWTKeySelect is the builder for selecting fields of JWTKey entities.
type JWTKeySelect struct {
	*JWTKeyQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (jks *JWTKeySelect) Aggregate(fns ...AggregateFunc) *JWTKeySelect {
	jks.fns = append(jks.fns, fns...)
	return jks
}

// Scan applies the selector query and scans the result into the given value.
func (jks *JWTKeySelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, jks.ctx, ent.OpQuerySelect)
	if err := jks.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*JWTKeyQuery, *JWTKeySelect](ctx, jks.JWTKeyQuery, jks, jks.inters, v)
}

func (jks *JWTKeySelect) sqlScan(ctx context.Context, root *JWTKeyQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(jks.fns))
	for _, fn := range jks.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*jks.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := jks.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (jks *JWTKeySelect) Modify(modifiers ...func(s *sql.Selector)) *JWTKeySelect {
	jks.modifiers = append(jks.modifiers, modifiers...)
	return jks
}
