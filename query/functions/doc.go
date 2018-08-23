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

this signature defines the named arguments for the function, along with their types.  There are a few helpers

*/
package functions
