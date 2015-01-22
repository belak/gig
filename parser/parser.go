package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gopkg.in/twik.v1"
	"gopkg.in/twik.v1/ast"
)

type Env struct {
	scope *twik.Scope
	fset  *ast.FileSet
}

type Node ast.Node

// Returns error for consistency with NewBootstrappedEnv
// Will never actually return an error
func NewEnv() (*Env, error) {
	fset := twik.NewFileSet()
	scope := twik.NewScope(fset)
	env := &Env{scope, fset}

	return env, nil
}

func NewBootstrappedEnv() (*Env, error) {
	env, _ := NewEnv()

	err := env.bootstrap()
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (e *Env) LoadTune(filename string) (Node, error) {
	return e.ParseFile(filename, "tune")
}

func (e *Env) ParseFile(filename, tune string) (Node, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return e.Parse(tune, data)
}

func (e *Env) Parse(tune string, data []byte) (Node, error) {
	return twik.Parse(e.fset, tune, data)
}

func (e *Env) Eval(node Node) (interface{}, error) {
	return e.scope.Eval(node)
}

// Bootstrapping functions

// Adds a field and a setter to the environment
func (e *Env) addField(name string, initial interface{}) {
	e.scope.Create("pkg-"+name, initial)
	e.scope.Create(name, func(args []interface{}) (interface{}, error) {
		if len(args) == 0 {
			// TODO: return error
			return nil, nil
		}

		err := e.scope.Set("pkg-"+name, args[0])

		return nil, err
	})
}

// Returns arguments as a list
// Requires at least one argument
//
// Called as:
//   (list <arg1> [arg2...])
func (e *Env) list(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("list takes at least one argument")
	}

	return args, nil
}

// Concatenates all arguments into one string
// Also expands all environment variables inline
//
// Called as:
//   (str <arg1> [arg2...])
func (e *Env) str(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("str takes at least one argument")
	}

	var arrArgs []interface{}

	if _, ok := args[0].([]interface{}); ok {
		arrArgs = args[0].([]interface{})
	} else {
		arrArgs = args
	}

	// TODO: don't assume all args are strings
	strArgs := []string{}
	for _, arg := range arrArgs {
		expanded := os.ExpandEnv(arg.(string))
		strArgs = append(strArgs, expanded)
	}

	return strings.Join(strArgs, ""), nil
}

// Displays a string to the user
// Requires at least one argument
//
// Called as:
//   (disp <string>)
func (e *Env) disp(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("printf takes a format string")
	}

	fmt.Println(args...)
	return nil, nil
}

func (e *Env) cd(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}

	err := os.Chdir(args[0].(string))

	return nil, err
}

func (e *Env) getPlatform(args []interface{}) (interface{}, error) {
	return runtime.GOOS, nil
}

// Will change
func (e *Env) shell(args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}

	cmd := args[0].(string)
	cmdArgs := []string{}
	if len(args) > 1 {
		for _, arg := range args[1:] {
			cmdArgs = append(cmdArgs, arg.(string))
		}
	}

	// TODO: implement
	out, err := exec.Command(cmd, cmdArgs...).Output()
	if err != nil {
		return nil, err
	}

	return string(out), nil
}

// Bootstraps the environment and creates our DSL
func (e *Env) bootstrap() error {
	// Create symbols list
	symbols := []struct {
		key   string
		value interface{}
	}{
		// Internal variables
		{
			"list",
			e.list,
		},
		{
			"str",
			e.str,
		},
		{
			"disp",
			e.disp,
		},
		{
			"cd",
			e.cd,
		},
		{
			"get-platform",
			e.getPlatform,
		},
		{
			"shell",
			e.shell,
		},
	}

	// Add symbols to scope
	for _, s := range symbols {
		err := e.scope.Create(s.key, s.value)
		if err != nil {
			fmt.Printf("Error bootstrapping: %s\n", err)
			os.Exit(1)
		}
	}

	// Load tunefile bootstrap
	bootstraps := []string{
		"tune-env",
		"default-config",
	}

	for _, strap := range bootstraps {
		node, err := e.ParseFile("bootstrap/"+strap+".tune", strap)
		if err != nil {
			return err
		}

		_, err = e.scope.Eval(node)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Env) GetString(name string) (string, error) {
	val, err := e.scope.Get(name)
	if err != nil {
		return "", err
	}

	return val.(string), err
}

func (e *Env) Get(name string) (interface{}, error) {
	return e.scope.Get(name)
}

func (e *Env) GetStringArray(name string) ([]string, error) {
	val, err := e.scope.Get(name)
	if err != nil {
		return nil, err
	}

	return val.([]string), err
}

func (e *Env) GetList(name string) ([]string, error) {
	val, err := e.scope.Get(name)
	if err != nil {
		return nil, err
	}

	var arrArgs []interface{}

	if _, ok := val.([]interface{}); ok {
		arrArgs = val.([]interface{})
	} else {
		arrArgs = []interface{}{val.(interface{})}
	}

	// TODO: don't assume all args are strings
	strArgs := []string{}
	for _, arg := range arrArgs {
		strArgs = append(strArgs, arg.(string))
	}

	return strArgs, nil
}

// Calls function defined in the tunefile
func (e *Env) Invoke(name string, args []interface{}) (interface{}, error) {
	sym, err := e.scope.Get(name)
	if err != nil {
		return nil, err
	}

	return sym.(func([]interface{}) (interface{}, error))(args)
}
