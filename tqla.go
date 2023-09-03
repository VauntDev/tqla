package tqla

import (
	"bytes"
	"fmt"
	"text/template"
)

type tqla struct {
	tmpl        *sqlTemplate
	parser      *sqlParser
	placeholder Placeholder
}

func New(options ...Option) (*tqla, error) {
	opts := defaultOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	parser := newSqlParser()
	funcs := template.FuncMap{
		"_sql_parser_": parser.parsefunc,
	}

	for k, v := range opts.funcs {
		if k == "_sql_parser_" {
			return nil, fmt.Errorf("invalid function name, _sql_parser_ is reserved")
		}
		funcs[k] = v
	}

	tmpl := newSqlTemplate("tqla", funcs)
	return &tqla{tmpl: tmpl,
		parser:      parser,
		placeholder: opts.placeholder,
	}, nil
}

func (t *tqla) Compile(statement string, data any) (string, []any, error) {
	if err := t.tmpl.parse(statement); err != nil {
		return "", nil, err
	}

	b := &bytes.Buffer{}
	if err := t.tmpl.execute(b, data); err != nil {
		return "", nil, err
	}

	sql, err := t.placeholder.Format(b.String())
	if err != nil {
		return "", nil, err
	}

	return sql, t.parser.args, nil
}
