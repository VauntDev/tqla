package tqla

import (
	"bytes"
	"text/template"
)

type tqla struct {
	tmpl   *sqlTemplate
	parser *sqlParser
	opts   *options
}

func New(options ...Option) (*tqla, error) {
	opts := defaultOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	placeHolder := newSqlParser()
	funcs := template.FuncMap{
		"_placeholder_": placeHolder.parsefunc,
	}

	tmpl := newSqlTemplate("tqla", funcs)
	return &tqla{tmpl: tmpl,
		parser: placeHolder,
		opts:   opts,
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

	sql, err := t.opts.placeholder.Format(b.String())
	if err != nil {
		return "", nil, err
	}

	return sql, t.parser.args, nil
}
