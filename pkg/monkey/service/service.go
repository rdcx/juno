package service

import (
	"encoding/json"
	"fmt"
	"juno/pkg/monkey/evaluator"
	"juno/pkg/monkey/lexer"
	"juno/pkg/monkey/object"
	"juno/pkg/monkey/parser"
	"os"

	_ "embed"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

//go:embed "data/stdlib.mon"
var stdlib string

func searchHtmlFun(args ...object.Object) object.Object {
	return &object.String{Value: "searchHtml"}
}

// Execute the supplied string as a program.
func (s *Service) Execute(input string) ([]byte, error) {

	env := object.NewEnvironment()
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		os.Exit(1)
	}

	// Register a function called version()
	// that the script can call.
	evaluator.RegisterBuiltin("searchHtml",
		func(env *object.Environment, args ...object.Object) object.Object {
			return (searchHtmlFun(args...))
		})

	//
	//  Parse and evaluate our standard-library.
	//
	initL := lexer.New(stdlib)
	initP := parser.New(initL)
	initProg := initP.ParseProgram()
	evaluator.Eval(initProg, env)

	//
	//  Now evaluate the code the user wanted to load.
	//
	//  Note that here our environment will still contain
	// the code we just loaded from our data-resource
	//
	//  (i.e. Our monkey-based standard library.)
	//
	out := evaluator.Eval(program, env)

	if out.Type() == object.ERROR_OBJ {
		return nil, fmt.Errorf(out.Inspect())
	}

	json, err := json.Marshal(out.ToInterface())

	if err != nil {
		return nil, err
	}

	return json, nil
}
