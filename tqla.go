package tqla

import (
	"bytes"
	"maps"
	"text/template"
)

type tqla struct {
	placeholder Placeholder
	funcs       template.FuncMap
}

// New creates a new tqla instance that can be used to generate dynamic queries via sql text/templating
func New(options ...Option) (*tqla, error) {
	opts := defaultOptions()
	for _, opt := range options {
		if err := opt.Apply(opts); err != nil {
			return nil, err
		}
	}

	return &tqla{
		placeholder: opts.placeholder,
		funcs:       opts.funcs,
	}, nil
}

// Compile accepts a string sql template and a data object that
// is used to build the db args. A finalized sql statement is returned with placeholders where dynamic
// arguments are passed as well as a slice of args. Returns a error if the sql template cannot be parsed. Compile
// is safe to call multiple times.
func (t *tqla) Compile(statement string, data any) (string, []any, error) {

	parser := newSqlParser()

	funcs := maps.Clone(t.funcs)
	funcs["_sql_parser_"] = parser.parsefunc

	tmpl := newSqlTemplate("tqla", funcs)

	if err := tmpl.parse(statement); err != nil {
		return "", nil, err
	}

	b := &bytes.Buffer{}
	if err := tmpl.execute(b, data); err != nil {
		return "", nil, err
	}

	sql, err := t.placeholder.Format(b.String())
	if err != nil {
		return "", nil, err
	}

	return sql, parser.args, nil
}
