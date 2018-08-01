package functions

import (
	"fmt"

	"github.com/influxdata/platform/query/compiler"

	"github.com/influxdata/platform/query/interpreter"
	"github.com/pkg/errors"

	"github.com/influxdata/platform/query"
	"github.com/influxdata/platform/query/execute"
	"github.com/influxdata/platform/query/plan"
	"github.com/influxdata/platform/query/semantic"
	"github.com/influxdata/platform/query/values"
)

const RenameKind = "rename"
const DropKind = "drop"
const KeepKind = "keep"

type RenameOpSpec struct {
	RenameCols map[string]string            `json:"columns"`
	RenameFn   *semantic.FunctionExpression `json:"fn"`
}

type DropOpSpec struct {
	DropCols      []string                     `json:"columns"`
	DropPredicate *semantic.FunctionExpression `json:"fn"`
}

type KeepOpSpec struct {
	KeepCols      []string                     `json:"columns"`
	KeepPredicate *semantic.FunctionExpression `json:"fn"`
}

var renameSignature = query.DefaultFunctionSignature()
var dropSignature = query.DefaultFunctionSignature()
var keepSignature = query.DefaultFunctionSignature()

func init() {
	renameSignature.Params["columns"] = semantic.Object
	renameSignature.Params["fn"] = semantic.Function

	query.RegisterFunction(RenameKind, createRenameOpSpec, renameSignature)
	query.RegisterOpSpec(RenameKind, newRenameOp)
	plan.RegisterProcedureSpec(RenameKind, newRenameProcedure, RenameKind)

	dropSignature.Params["columns"] = semantic.NewArrayType(semantic.String)
	dropSignature.Params["fn"] = semantic.Function

	query.RegisterFunction(DropKind, createDropOpSpec, dropSignature)
	query.RegisterOpSpec(DropKind, newDropOp)
	plan.RegisterProcedureSpec(DropKind, newDropProcedure, DropKind)

	keepSignature.Params["columns"] = semantic.Object
	keepSignature.Params["fn"] = semantic.Function

	query.RegisterFunction(KeepKind, createKeepOpSpec, keepSignature)
	query.RegisterOpSpec(KeepKind, newKeepOp)
	plan.RegisterProcedureSpec(KeepKind, newDropProcedure, KeepKind)

	execute.RegisterTransformation(RenameKind, createRenameDropTransformation)
	execute.RegisterTransformation(DropKind, createRenameDropTransformation)
	execute.RegisterTransformation(KeepKind, createRenameDropTransformation)
}

func createRenameOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}
	var cols values.Object
	if c, ok, err := args.GetObject("columns"); err != nil {
		return nil, err
	} else if ok {
		cols = c
	}

	var renameFn *semantic.FunctionExpression
	if f, ok, err := args.GetFunction("fn"); err != nil {
		return nil, err
	} else if ok {
		if fn, err := interpreter.ResolveFunction(f); err != nil {
			return nil, err
		} else {
			renameFn = fn
		}
	}

	if cols == nil && renameFn == nil {
		return nil, errors.New("rename error: neither column list nor map function provided")
	}

	if cols != nil && renameFn != nil {
		return nil, errors.New("rename error: both column list and map function provided")
	}

	spec := &RenameOpSpec{
		RenameFn: renameFn,
	}

	if cols != nil {
		var err error
		renameCols := make(map[string]string, cols.Len())
		// Check types of object values manually
		cols.Range(func(name string, v values.Value) {
			if err != nil {
				return
			}
			if v.Type() != semantic.String {
				err = fmt.Errorf("rename error: columns object contains non-string value of type %s", v.Type())
				return
			}
			renameCols[name] = v.Str()
		})
		if err != nil {
			return nil, err
		}
		spec.RenameCols = renameCols
	}

	return spec, nil
}

func createDropOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	var cols values.Array
	if c, ok, err := args.GetArray("columns", semantic.String); err != nil {
		return nil, err
	} else if ok {
		cols = c
	}

	var dropPredicate *semantic.FunctionExpression
	if f, ok, err := args.GetFunction("fn"); err != nil {
		return nil, err
	} else if ok {
		fn, err := interpreter.ResolveFunction(f)
		if err != nil {
			return nil, err
		}

		dropPredicate = fn
	}

	if cols == nil && dropPredicate == nil {
		return nil, errors.New("drop error: neither column list nor predicate function provided")
	}

	if cols != nil && dropPredicate != nil {
		return nil, errors.New("drop error: both column list and predicate provided")
	}

	var dropCols []string
	var err error
	if cols != nil {
		dropCols, err = interpreter.ToStringArray(cols)
		if err != nil {
			return nil, err
		}
	}

	return &DropOpSpec{
		DropCols:      dropCols,
		DropPredicate: dropPredicate,
	}, nil
}

func createKeepOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	var cols values.Array
	if c, ok, err := args.GetArray("columns", semantic.String); err != nil {
		return nil, err
	} else if ok {
		cols = c
	}

	var keepPredicate *semantic.FunctionExpression
	if f, ok, err := args.GetFunction("fn"); err != nil {
		return nil, err
	} else if ok {
		fn, err := interpreter.ResolveFunction(f)
		if err != nil {
			return nil, err
		}

		keepPredicate = fn
	}

	if cols == nil && keepPredicate == nil {
		return nil, errors.New("keep error: neither column list nor predicate function provided")
	}

	if cols != nil && keepPredicate != nil {
		return nil, errors.New("keep error: both column list and predicate provided")
	}

	var keepCols []string
	var err error
	if cols != nil {
		keepCols, err = interpreter.ToStringArray(cols)
		if err != nil {
			return nil, err
		}
	}

	return &KeepOpSpec{
		KeepCols:      keepCols,
		KeepPredicate: keepPredicate,
	}, nil
}

func newRenameOp() query.OperationSpec {
	return new(RenameOpSpec)
}

func (s *RenameOpSpec) Kind() query.OperationKind {
	return RenameKind
}

func newDropOp() query.OperationSpec {
	return new(DropOpSpec)
}

func (s *DropOpSpec) Kind() query.OperationKind {
	return DropKind
}

func newKeepOp() query.OperationSpec {
	return new(KeepOpSpec)
}

func (s *KeepOpSpec) Kind() query.OperationKind {
	return KeepKind
}

type RenameDropProcedureSpec struct {
	RenameCols map[string]string
	RenameFn   *semantic.FunctionExpression
	// The same field is used for both columns to drop and columns to keep
	DropKeepCols map[string]bool
	// the same field is used for the drop predicate and the keep predicate
	DropKeepPredicate *semantic.FunctionExpression
	// Denotes whether we're going to do a drop or a keep
	KeepSpecified bool
}

func newRenameProcedure(qs query.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	s, ok := qs.(*RenameOpSpec)

	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	var renameCols map[string]string
	if s.RenameCols != nil {
		renameCols = s.RenameCols
	}

	return &RenameDropProcedureSpec{
		RenameCols: renameCols,
		RenameFn:   s.RenameFn,
	}, nil
}

func (s *RenameDropProcedureSpec) Kind() plan.ProcedureKind {
	return RenameKind
}

func (s *RenameDropProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(RenameDropProcedureSpec)
	ns.RenameCols = s.RenameCols
	ns.DropKeepCols = s.DropKeepCols
	ns.RenameFn = s.RenameFn
	ns.DropKeepPredicate = s.DropKeepPredicate
	return ns
}

func (s *RenameDropProcedureSpec) PushDownRules() []plan.PushDownProcedureSpec {
	return nil
}

func (s *RenameDropProcedureSpec) PushDown() (root *plan.Procedure, dup func() *plan.Procedure) {
	return nil, nil
}

// Keep and Drop are inverses, so they share the same procedure constructor
func newDropProcedure(qs query.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	pr := &RenameDropProcedureSpec{}
	switch s := qs.(type) {
	case *DropOpSpec:
		if s.DropCols != nil {
			pr.DropKeepCols = toStringSet(s.DropCols)
		}
		pr.DropKeepPredicate = s.DropPredicate
	case *KeepOpSpec:
		// Flip use of dropCols field from drop to keep
		pr.KeepSpecified = true
		if s.KeepCols != nil {
			pr.DropKeepCols = toStringSet(s.KeepCols)
		}
		pr.DropKeepPredicate = s.KeepPredicate

	default:
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}
	return pr, nil
}

func toStringSet(arr []string) map[string]bool {
	if arr == nil {
		return nil
	}
	ret := make(map[string]bool, len(arr))
	for _, s := range arr {
		ret[s] = true
	}
	return ret
}

func createRenameDropTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)

	t, err := NewRenameDropTransformation(d, cache, spec)
	if err != nil {
		return nil, nil, err
	}
	return t, d, nil
}

type renameDropTransformation struct {
	d                 execute.Dataset
	cache             execute.TableBuilderCache
	renameCols        map[string]string
	renameFn          compiler.Func
	renameScope       compiler.Scope
	renameColParam    string
	dropKeepCols      map[string]bool
	dropExclusiveCols map[string]bool
	keepSpecified     bool
	dropKeepPredicate compiler.Func
	dropKeepScope     compiler.Scope
	dropKeepColParam  string
}

func newFunc(fn *semantic.FunctionExpression, types [2]semantic.Type) (compiler.Func, string, error) {
	scope, decls := query.BuiltIns()
	compileCache := compiler.NewCompilationCache(fn, scope, decls)
	if len(fn.Params) != 1 {
		return nil, "", fmt.Errorf("function should only have a single parameter, got %d", len(fn.Params))
	}
	paramName := fn.Params[0].Key.Name

	compiled, err := compileCache.Compile(map[string]semantic.Type{
		paramName: types[0],
	})
	if err != nil {
		return nil, "", err
	}

	if compiled.Type() != types[1] {
		return nil, "", fmt.Errorf("provided function does not evaluate to type %s", types[1].Kind())
	}

	return compiled, paramName, nil
}

func NewRenameDropTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec plan.ProcedureSpec) (*renameDropTransformation, error) {
	s, ok := spec.(*RenameDropProcedureSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", spec)
	}

	var renameMapFn compiler.Func
	var renameScope compiler.Scope
	var renameColParam string
	if s.RenameFn != nil {
		compiledFn, param, err := newFunc(s.RenameFn, [2]semantic.Type{semantic.String, semantic.String})
		if err != nil {
			return nil, err
		}
		renameMapFn = compiledFn
		renameColParam = param
		// Scope for calling the rename map function
		renameScope = make(map[string]values.Value, 1)
	}

	var dropKeepPredicate compiler.Func
	var dropKeepColParam string
	var dropKeepScope compiler.Scope
	if s.DropKeepPredicate != nil {
		compiledFn, param, err := newFunc(s.DropKeepPredicate, [2]semantic.Type{semantic.String, semantic.Bool})
		if err != nil {
			return nil, err
		}

		dropKeepPredicate = compiledFn
		dropKeepColParam = param
		// Scope for calling the drop/keep predicate function
		dropKeepScope = make(map[string]values.Value, 1)
	}

	return &renameDropTransformation{
		d:                 d,
		cache:             cache,
		renameCols:        s.RenameCols,
		renameFn:          renameMapFn,
		renameScope:       renameScope,
		renameColParam:    renameColParam,
		dropKeepCols:      s.DropKeepCols,
		dropKeepPredicate: dropKeepPredicate,
		dropKeepScope:     dropKeepScope,
		dropKeepColParam:  dropKeepColParam,
		keepSpecified:     s.KeepSpecified,
	}, nil
}

func (t *renameDropTransformation) keepToDropCols(tbl query.Table) error {
	cols := tbl.Cols()

	// Check to make sure we aren't trying to `drop` or `keep` a column which
	// is not present in the table.
	if t.dropKeepCols != nil {
		for c := range t.dropKeepCols {
			if execute.ColIdx(c, cols) < 0 {
				if t.keepSpecified {
					return fmt.Errorf(`keep error: column "%s" doesn't exist`, c)
				}
				return fmt.Errorf(`drop error: column "%s" doesn't exist`, c)
			}
		}
	}

	// If `keepSpecified` is true, i.e., we want to exclusively keep the columns listed in dropKeepCols
	// as opposed to exclusively dropping them. So in that case, we invert the dropKeepCols map and store it
	// in exclusiveDropCols; exclusiveDropCols may be changed with each call to `Process`, but
	// `dropKeepCols` will not be.
	if t.keepSpecified && t.dropKeepCols != nil {
		exclusiveDropCols := make(map[string]bool, len(tbl.Cols()))
		for _, c := range tbl.Cols() {
			if _, ok := t.dropKeepCols[c.Label]; !ok {
				exclusiveDropCols[c.Label] = true
			}
		}
		t.dropExclusiveCols = exclusiveDropCols
	} else if t.dropKeepCols != nil {
		t.dropExclusiveCols = t.dropKeepCols
	}

	return nil
}

func (t *renameDropTransformation) checkColumnReferences(tbl query.Table) error {
	cols := tbl.Cols()
	if t.renameCols != nil {
		for c := range t.renameCols {
			if execute.ColIdx(c, cols) < 0 {
				return fmt.Errorf(`rename error: column "%s" doesn't exist`, c)
			}
		}
	}

	if t.dropKeepCols != nil && t.renameCols != nil {
		for k := range t.renameCols {
			if _, ok := t.dropKeepCols[k]; ok {
				return fmt.Errorf(`rename error: cannot rename column "%s" which is marked for drop`, k)
			}
		}
	}

	return nil
}

func (t *renameDropTransformation) shouldDrop(col string) (bool, error) {
	t.dropKeepScope[t.dropKeepColParam] = values.NewStringValue(col)
	if shouldDrop, err := t.dropKeepPredicate.EvalBool(t.dropKeepScope); err != nil {
		return false, err
	} else if t.keepSpecified {
		return !shouldDrop, nil
	} else {
		return shouldDrop, nil
	}
}

func (t *renameDropTransformation) shouldDropCol(col string) (bool, error) {
	if t.dropExclusiveCols != nil {
		if _, exists := t.dropExclusiveCols[col]; exists {
			return true, nil
		}
	} else if t.dropKeepPredicate != nil {
		return t.shouldDrop(col)
	}
	return false, nil
}

func (t *renameDropTransformation) renameCol(col *query.ColMeta) error {
	if col == nil {
		return errors.New("rename error: cannot rename nil column")
	}
	if t.renameCols != nil {
		if newName, ok := t.renameCols[col.Label]; ok {
			col.Label = newName
		}
	} else if t.renameFn != nil {
		t.renameScope[t.renameColParam] = values.NewStringValue(col.Label)
		newName, err := t.renameFn.EvalString(t.renameScope)
		if err != nil {
			return err
		}
		col.Label = newName
	}
	return nil
}

func (t *renameDropTransformation) Process(id execute.DatasetID, tbl query.Table) error {
	// Adjust internal column list depending
	// on whether this is a keep or a drop operation
	if err := t.keepToDropCols(tbl); err != nil {
		return err
	}

	// Check to make sure we don't have duplicates between the drop and rename column lists,
	// or whether we're rename a column that doesn't exist.
	if err := t.checkColumnReferences(tbl); err != nil {
		return err
	}

	keyCols := make([]query.ColMeta, 0, len(tbl.Cols()))
	keyValues := make([]values.Value, 0, len(tbl.Cols()))
	builderCols := make([]query.ColMeta, 0, len(tbl.Cols()))
	// If we remove columns, column indices will be different between the
	// builder and the table - we need to keep track
	colMap := make([]int, 0, len(tbl.Cols()))

	for i, c := range tbl.Cols() {
		if shouldDrop, err := t.shouldDropCol(c.Label); err != nil {
			return err
		} else if shouldDrop {
			continue
		}

		keyIdx := execute.ColIdx(c.Label, tbl.Key().Cols())
		keyed := keyIdx >= 0

		if err := t.renameCol(&c); err != nil {
			return err
		}

		if keyed {
			keyCols = append(keyCols, c)
			keyValues = append(keyValues, tbl.Key().Value(keyIdx))
		}

		colMap = append(colMap, i)
		builderCols = append(builderCols, c)
	}

	key := execute.NewGroupKey(keyCols, keyValues)
	builder, created := t.cache.TableBuilder(key)
	if created {
		for _, c := range builderCols {
			builder.AddCol(c)
		}
	}

	err := tbl.Do(func(cr query.ColReader) error {
		for i := 0; i < cr.Len(); i++ {
			execute.AppendMappedRecord(i, cr, builder, colMap)
		}
		return nil
	})

	return err
}

func (t *renameDropTransformation) RetractTable(id execute.DatasetID, key query.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *renameDropTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}

func (t *renameDropTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}

func (t *renameDropTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
