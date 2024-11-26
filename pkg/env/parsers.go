package env

import (
	"errors"
	"github.com/beetbasket/program/pkg/env/internal/parsers"
	"os"
	"reflect"
)

func init() {
	RegisterParser[parsers.String]()
	RegisterParser[parsers.AbsFile]()
	RegisterParser[parsers.Args]()
	RegisterParser[parsers.Port]()
}

type Parser[T any] interface {
	Name() string
	Parse(string) (T, error)
}

var parserMap = map[string]func(envKey string, defaultValue *string) (reflect.Value, error){}

var ErrNotFound = errors.New("env value not found")

func RegisterParser[P Parser[T], T any]() {
	parserMap[parserName[P, T]()] = parse[P, T]
}

func parse[P Parser[T], T any](envKey string, defaultValue *string) (reflect.Value, error) {
	var v string
	if ev, ok := os.LookupEnv(envKey); ok {
		v = ev
	} else if defaultValue != nil {
		v = *defaultValue
	} else {
		return reflect.Value{}, ErrNotFound
	}

	pv, err := (*new(P)).Parse(v)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(pv), nil
}

func parserName[P Parser[T], T any]() string {
	var p P
	return p.Name()
}