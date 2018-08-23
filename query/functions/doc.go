/*
Package functions is a collection of built-in functions that are callable in the flux query processor.  While flux may
be extended at runtime by writing function expressions, there are some limitations for which a built-in function is
necessary, such as the need for custom data structures, stateful procedures, complex looping and branching, and
connection to external services.  Another reason for implementing a built-in function is to provide a function that is
broadly applicaple for many users (e.g., sum() or max()).

The functions package is rarely accessed as a direct API.  Rather, the query processing engine registers various
interfaces implemented within the functions package and executes them generically using an API that is common to all
functions.  The registration process is executed by running the init() function in each function file, and is then finalized
by importing package builtin, which itself imports the functions package and runs a final setup routine that finalizes
installation of the builtin functions to the query processor.

Because of this design, a built-in function implementation consists of a bundle of different interface implementations
that are required at various phases of query execution.  It's important that each of the interfaces is completely
implemented and correctly registered with the query processor so that the function can be properly initialized, planned,
and executed in the context of a whole query.  The remainder of this documentation details the different interfaces that
must be implemented in order to properly implement a built-in function as part of the query processor.

Query Package Registrations

The query package is responsible for verifying the syntax, parsing arguments and generating internal function representations
for each built-in. The following registrations are defined in query/compile.go:

	// RegisterFunction adds a new builtin top level function.
	// name: the name of the function as it would be called
	// c: a function reference with the signature func(args Arguments, a *Administration) (OperationSpec, error)
	// sig: a function signature type that specifies the names and types of each argument for the function
	func RegisterFunction(name string, c CreateOperationSpec, sig semantic.FunctionSignature)

	// RegisterFunctionWithSideEffect adds a new builtin top level function that produces side effects.
	// For example, the builtin functions yield(), toKafka(), and toHTTP() all produce side effects.
	// name: the name of the function as it would be called
	// c: a function reference with the signature func(args Arguments, a *Administration) (OperationSpec, error)
	// sig: a function signature type that specifies the names and types of each argument for the function
	func RegisterFunctionWithSideEffect(name string, c CreateOperationSpec, sig semantic.FunctionSignature)

	// RegisterOpSpec registers an operation spec with a given kind.
	// k: a label that uniquely identifies this operation. If the kind has already been registered the call panics.
	// c: a function reference that creates a new, default-initialized opSpec for the given kind.
	func RegisterOpSpec(k OperationKind, c NewOperationSpec)

	// RegisterBuiltIn adds any variable declarations written in flux script to the builtin scope.
	func RegisterBuiltIn(name, script string)

There are several types that you must implement to properly register your new function in the query package:

semantic.FunctionSignature:

this signature defines the named arguments for the function, along with their types.  All builtin functions provide
an instance of the same type, semantic.FunctionSignature and it's the implementor's job to populate its members with
the correct values for the new function.  In many cases, it's simplest to start with a default signature and add custom
parameters to it:

	var covarianceSignature = query.DefaultFunctionSignature()
	func init() {
		covarianceSignature.Params["pearsonr"] = semantic.Bool
		covarianceSignature.Params["valueDst"] = semantic.String
		...
	}

Alternatively, a full instantiation can be used:

	var joinSignature = semantic.FunctionSignature{
		Params: map[string]semantic.Type{
			"tables": semantic.Object,
			"on":     semantic.NewArrayType(semantic.String),
			"method": semantic.String,
		},
		ReturnType:   query.TableObjectType,
		PipeArgument: "tables",
	}

The second type that must be implemented is the OpSpec.  In this case, the query package expects a custom implementation
of the query.OperationSpec type.  Only the Kind() function is in the interface, but two helper functions must also be
implemented and registered with the query package, a create function and a new function.
These are best understood by example, as for the built-in covariance function:

	// define the fields that will be needed to compute this function.
	type CovarianceOpSpec struct {
		PearsonCorrelation bool   `json:"pearsonr"`
		ValueDst           string `json:"valueDst"`
		// arguments that are common to all Aggregate functions are defined in a shared type
		execute.AggregateConfig
	}

	const CovarianceKind = "covariance"
	func createCovarianceOpSpec(args query.Arguments, a *query.Administration) (query.OperationSpec, error) {
		if err := a.AddParentFromArgs(args); err != nil {
			return nil, err
		}

		spec := new(CovarianceOpSpec)
		pearsonr, ok, err := args.GetBool("pearsonr")
		if err != nil {
			return nil, err
		} else if ok {
			spec.PearsonCorrelation = pearsonr
		}

		label, ok, err := args.GetString("valueDst")
		if err != nil {
			return nil, err
		} else if ok {
			spec.ValueDst = label
		} else {
			spec.ValueDst = execute.DefaultValueColLabel
		}

		if err := spec.AggregateConfig.ReadArgs(args); err != nil {
			return nil, err
		}
		if len(spec.Columns) != 2 {
			return nil, errors.New("must provide exactly two columns")
		}
		return spec, nil
	}

	func newCovarianceOp() query.OperationSpec {
		return new(CovarianceOpSpec)
	}

	func (s *CovarianceOpSpec) Kind() query.OperationKind {
		return CovarianceKind
	}

To summarize this section, a new function implementation must instantiate and register a query.Signature, define and
register an implementation of a query.OperationSpec, and implement the Kind(), createXXOpSpec and newXXOpSpec functions.
The createXXOpSpec function is registered so that the query processor knows how to parse the function arguments and create
a proper OperationSpec representation of the function call.  The newXXOp function is registered to aid the query processor in
allocating a new OperationSpec struct of the proper type, and the Kind() function is implemented to inform the system
about the OperationSpec's actual type when the type would otherwise be ambiguous.

The end result of these registrations is that the query processor has the information it needs to create an internal
representation of a function call that can then be consumed by the logical planner that will produce an execution plan for
a complete query.

*/
package functions
