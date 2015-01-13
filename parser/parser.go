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

// Will change
// Sets the install function that is called to install the package
// Expects one function as an argument
func (e *Env) install(args []interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-install", args[0])

    return nil, err
}

// Sets the tunefile version
// Expects one float as argument
//
// Called as:
//   (tune-version <version>)
func (e *Env) setTuneVersion(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("tune-version", args[0])

    return nil, err
}

// Sets the package name
// Expects one string as argument
//
// Called as:
//   (name <name>)
func (e *Env) setName(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-name", args[0])

    return nil, err
}

// Sets the package description
// Expects one string as argument
//
// Called as:
//   (description <description>)
func (e *Env) setDescription(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-description", args[0])

    return nil, err
}

// Sets the package license
// Expects one string as argument
//
// Called as:
//   (license <license>)
func (e *Env) setLicense(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-license", args[0])

    return nil, err
}

// Sets the package version
// Expects one float as argument
//
// Called as:
//   (version <version>)
func (e *Env) setVersion(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-version", args[0])

    return nil, err
}

// Sets the package homepage
// Expects one string as argument
//
// Called as:
//   (homepage <url>)
func (e *Env) setHomepage(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-homepage", args[0])

    return nil, err
}

// Sets the package source URL
// Expects one string as argument
//
// Called as:
//   (url <url>)
func (e *Env) setUrl(args[]interface{}) (interface{}, error) {
    if len(args) == 0 {
        // TODO: return error
        return nil, nil
    }

    err := e.scope.Set("pkg-url", args[0])

    return nil, err
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

    // Create symbols list
    symbols := []struct{
        key string
        value interface{}
    }{
        // Internal variables
        {
            "pkg-name",
            "",
        },
        {
            "pkg-description",
            "",
        },
        {
            "pkg-license",
            "",
        },
        {
            "pkg-version",
            0.0,
        },
        {
            "pkg-homepage",
            "",
        },
        {
            "pkg-url",
            "",
        },
        {
            "pkg-install",
            func([]interface{}) (interface{}, error) { return nil, nil },
        },
        {
            "pkg-dependencies",
            []string{},
        },

        // Setters
        {
            "tune-version",
            e.setTuneVersion,
        },
        {
            "name",
            e.setName,
        },
        {
            "description",
            e.setDescription,
        },
        {
            "license",
            e.setLicense,
        },
        {
            "version",
            e.setVersion,
        },
        {
            "homepage",
            e.setHomepage,
        },
        {
            "url",
            e.setUrl,
        },
        {
            "install",
            e.install,
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
