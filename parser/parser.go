package parser

import (
    "fmt"
    "os"
    "io/ioutil"

    "gopkg.in/twik.v1"
    "gopkg.in/twik.v1/ast"
)

type Env struct {
    scope *twik.Scope
    fset *ast.FileSet
}

type Node ast.Node

func NewEnv() *Env {
    fset := twik.NewFileSet()
    scope := twik.NewScope(fset)
    env := &Env{scope, fset}

    // Add our stuff
    env.bootstrap()

    return env
}

func (e *Env) ParseFile(filename string) (Node, error) {
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        return nil, err
    }

    return e.Parse(data)
}

func (e *Env) Parse(data []byte) (Node, error) {
    return twik.Parse(e.fset, "Tune", data)
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

// Sets the package dependencies
// Expects at least one string as argument
//
// Called as:
//   (depends-on <dep1> [dep2...])
func (e *Env) dependsOn(args []interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    var deps []string
    for _, dep := range args {
        deps = append(deps, dep.(string))
    }

    err := e.scope.Set("pkg-dependencies", deps)

    return nil, err
}

// Will change
func (e *Env) cd(args[]interface{}) (interface{}, error) {
    if len(args) != 1 {
        // TODO: return error
        return nil, nil
    }

    // TODO: implement

    return nil, nil
}

// Will change
func (e *Env) shell(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        return nil, nil
    }

    // TODO: implement

    return nil, nil
}

// Will change
func (e *Env) setEnv(args[]interface{}) (interface{}, error) {
    if len(args) != 2 {
        // TODO: return error
        return nil, nil
    }

    // TODO: implement

    return nil, nil
}

// Bootstraps the environment and creates our DSL
func (e *Env) bootstrap() {
    // TODO: maybe switch to naming a param, giving a type,
    // and then automatically create "private" param, getter,
    // and setter
    // TODO: make generic getters and setters (pass types somehow?)

    fields := []struct{
        name string
        initial interface{}
    }{
        {
            "tune-version",
            "",
        },
        {
            "name",
            "",
        },
        {
            "description",
            "",
        },
        {
            "license",
            "",
        },
        {
            "version",
            "",
        },
        {
            "homepage",
            "",
        },
        {
            "url",
            "",
        },
        {
            "install",
            func([]interface{}) (interface{}, error) { return nil, nil },
        },
    }

    for _, field := range fields {
        e.addField(field.name, field.initial)
    }

    // Create symbols list
    symbols := []struct{
        key string
        value interface{}
    }{
        // Internal variables
        {
            "pkg-dependencies",
            []string{},
        },

        // Setters
        {
            "disp",
            e.disp,
        },
        {
            "cd",
            e.cd,
        },
        {
            "shell",
            e.shell,
        },
        {
            "set-env",
            e.setEnv,
        },
        {
            "depends-on",
            e.dependsOn,
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
}

func (e *Env) GetString(name string) (string, error) {
    val, err := e.scope.Get(name)
    if err != nil {
        return "", err
    }

    return val.(string), err
}

func (e *Env) GetStringArray(name string) ([]string, error) {
    val, err := e.scope.Get(name)
    if err != nil {
        return nil, err
    }

    return val.([]string), err
}

// Calls function defined in the tunefile
func (e *Env) Invoke(name string, args []interface{}) (interface{}, error) {
    sym, err := e.scope.Get(name)
    if err != nil {
        return nil, err
    }

    return sym.(func([]interface{})(interface{}, error))(args)
}
