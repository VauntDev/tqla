package tqla

import (
	"fmt"
	"text/template"
)

type options struct {
	placeholder Placeholder
	funcs       template.FuncMap
}

type Option interface {
	Apply(*options) error
}

type funcOption struct {
	f func(*options) error
}

func (flo *funcOption) Apply(con *options) error {
	return flo.f(con)
}

func newFuncOption(f func(*options) error) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithPlaceHolder(p Placeholder) Option {
	return newFuncOption(func(o *options) error {
		o.placeholder = p
		return nil
	})
}

func WithFuncMap(funcs template.FuncMap) Option {
	return newFuncOption(func(o *options) error {
		if o.funcs == nil {
			o.funcs = template.FuncMap{}
		}
		for k, v := range funcs {
			if k == "_sql_parser_" {
				return fmt.Errorf("invalid function name, _sql_parser_ is reserved")
			}
			o.funcs[k] = v
		}
		return nil
	})
}

func defaultOptions() *options {
	return &options{
		placeholder: Question,
		funcs:       template.FuncMap{},
	}
}
