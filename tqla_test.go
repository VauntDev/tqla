package tqla

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

type testExample struct {
	FirstName string
	LastName  string
	Page      *testExampePage
}

type testExampePage struct {
	Id    int
	Dttm  string
	Limit int
}

func TestTqla(t *testing.T) {
	// optional functions
	funcs := template.FuncMap{
		"add": func(x int, y int) int {
			return x + y
		},
	}

	testCases := []struct {
		name         string
		templateSql  string
		options      []Option
		data         any
		expectedSql  string
		expectedArgs []any
	}{
		{
			name: "condition column where filter",
			templateSql: `{{ $column:=true }}
				select * from table where v = {{.Value }}
				{{if $column }}
					and c = {{ $column }}
				{{end}}`,
			data: struct {
				Value int
			}{Value: 5},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select * from table where v = $1 and c = $2`,
			expectedArgs: []any{
				5,
				true,
			},
		},
		{
			name: "condition column where filter and custom function",
			templateSql: `{{ $column:=true }}
				select * from table where v = {{ add .Value .Value }}
				{{ if $column }}
					and c = {{ $column }}
				{{ end }}`,
			data: struct {
				Value int
			}{Value: 5},
			options:     []Option{WithPlaceHolder(Dollar), WithFuncMap(funcs)},
			expectedSql: `select * from table where v = $1 and c = $2`,
			expectedArgs: []any{
				10,
				true,
			},
		},
		{
			name: "where in string slice filter",
			templateSql: `
				select * from table where v in ( {{ .Value }} )`,
			data: struct {
				Value []string
			}{Value: []string{"v1", "v2", "v3"}},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select * from table where v in ( $1 )`,
			expectedArgs: []any{
				[]string{"v1", "v2", "v3"},
			},
		},
		{
			name: "complex condition where filter empty value",
			templateSql: `
				select c1,
				c2,
				c3,
				COALESCE(c4,'')
				from d.table as t
				where c1 ilike ({{ .FirstName }})
				{{ if .LastName }}
					and c2 = {{ .LastName }}
				{{ end }}
				{{ if .Page }}
					and (c3,c4) > ({{ .Page.Id }}, {{ .Page.Dttm }})
				{{ end }}`,
			data: &testExample{
				FirstName: "test",
				Page: &testExampePage{
					Id:   1,
					Dttm: "2023-08-08",
				},
			},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select c1, c2, c3, COALESCE(c4,'') from d.table as t where c1 ilike ($1) and (c3,c4) > ($2, $3)`,
			expectedArgs: []any{"test",
				1,
				"2023-08-08",
			},
		},
		{
			name: "complex condition where filter full value",
			templateSql: `
				select c1,
				c2,
				c3,
				COALESCE(c4,'')
				from d.table as t
				where c1 ilike ({{ .FirstName }})
				{{ if .LastName }}
					and c2 = {{ .LastName }}
				{{ end }}
				{{ if .Page }}
					and (c3,c4) > ({{ .Page.Id }}, {{ .Page.Dttm }})
				{{ end }}`,
			data: &testExample{
				FirstName: "test",
				LastName:  "test",
				Page: &testExampePage{
					Id:   1,
					Dttm: "2023-08-08",
				},
			},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select c1, c2, c3, COALESCE(c4,'') from d.table as t where c1 ilike ($1) and c2 = $2 and (c3,c4) > ($3, $4)`,
			expectedArgs: []any{"test",
				"test",
				1,
				"2023-08-08",
			},
		},
		{
			name: "complex condition where filter missing page",
			templateSql: `
				select c1,
				c2,
				c3,
				COALESCE(c4,'')
				from d.table as t
				where c1 ilike ({{ .FirstName }})
				{{ if .LastName }}
					and c2 = {{ .LastName }}
				{{ end }}
				{{ if .Page }}
					and (c3,c4) > ({{ .Page.Id }}, {{ .Page.Dttm }})
				{{ end }}`,
			data: &testExample{
				FirstName: "test",
				LastName:  "test",
			},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select c1, c2, c3, COALESCE(c4,'') from d.table as t where c1 ilike ($1) and c2 = $2`,
			expectedArgs: []any{"test",
				"test",
			},
		},
		{
			name: "complex condition where filter full value with max limit",
			templateSql: `
				{{ $maxLimit := 100 }}
				select c1,
					c2,
					c3,
					COALESCE(c4,'')
				from d.table as t
				where c1 ilike ({{ .FirstName }})
				{{ if .LastName }}
					and c2 = {{ .LastName }}
				{{ end }}
				{{ if .Page }}
					and (c3,c4) > ({{ .Page.Id }}, {{ .Page.Dttm }})
				{{ end }}
				{{ if gt .Page.Limit $maxLimit }}
					LIMIT {{ $maxLimit }}
				{{ else }}
					LIMIT {{ .Page.Limit }}
				{{ end }}`,
			data: &testExample{
				FirstName: "test",
				LastName:  "test",
				Page: &testExampePage{
					Id:    1,
					Dttm:  "2023-08-08",
					Limit: 10,
				},
			},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select c1, c2, c3, COALESCE(c4,'') from d.table as t where c1 ilike ($1) and c2 = $2 and (c3,c4) > ($3, $4) LIMIT $5`,
			expectedArgs: []any{"test",
				"test",
				1,
				"2023-08-08",
				10,
			},
		},
		{
			name: "complex nested query/template",
			templateSql: `
			{{ define "nested_select" }}
				select count(*) as count_v
				from d.table as t
				join d.table_2 as t2 on t.c1 = t2.c2
				where t.c3 ilike {{ .FirstName }} and t.c4 = 'value'
			{{end}}
			{{ $maxLimit := 100 }}
			select c1,
				c2,
				c3,
				COALESCE(c4,''),
				({{ template "nested_select" . }})
			from d.table as t
			where c1 ilike ({{ .FirstName }})
			{{ if .LastName }}
				and c2 = {{ .LastName }}
			{{ end }}
			{{ if .Page }}
				and (c3,c4) > ({{ .Page.Id }}, {{ .Page.Dttm }})
			{{ end }}
			{{ if gt .Page.Limit $maxLimit }}
				LIMIT {{ $maxLimit }}
			{{ else }}
				LIMIT {{ .Page.Limit }}
			{{ end }}`,
			data: &testExample{
				FirstName: "test",
				LastName:  "test",
				Page: &testExampePage{
					Id:    1,
					Dttm:  "2023-08-08",
					Limit: 10,
				},
			},
			options:     []Option{WithPlaceHolder(Dollar)},
			expectedSql: `select c1, c2, c3, COALESCE(c4,''), ( select count(*) as count_v from d.table as t join d.table_2 as t2 on t.c1 = t2.c2 where t.c3 ilike $1 and t.c4 = 'value' ) from d.table as t where c1 ilike ($2) and c2 = $3 and (c3,c4) > ($4, $5) LIMIT $6`,
			expectedArgs: []any{"test",
				"test",
				"test",
				1,
				"2023-08-08",
				10,
			},
		},
		{
			name: "order by",
			templateSql: `
			{{ $len := len .Columns }}
			SELECT * FROM t1 WHERE username = {{ .Username }} 
			{{ if gt $len 0 }}
			ORDER BY
				{{ $pos := 0}}
				{{ range $i, $v := .Columns }}
					{{ if eq $v "username" }}
						username
					{{ end }}
					{{ $pos = add $pos 1 }}
					{{ if eq $pos $len }}{{ else }},{{ end }}
				{{ end }} 
				{{ if eq .Order "ASC" }} 
					ASC
				{{ else }}
					DESC
				{{ end }}
			{{ end }}`,

			data: &struct {
				Username string
				Order    string
				Columns  []string
			}{
				Username: "test",
				Order:    "ASC",
				Columns:  []string{"username"},
			},
			options: []Option{
				WithPlaceHolder(Dollar), WithFuncMap(funcs)},
			expectedSql:  `SELECT * FROM t1 WHERE username = $1 ORDER BY username ASC`,
			expectedArgs: []any{"test"},
		},
	}

	for _, testCase := range testCases {
		t.Run("test", func(t *testing.T) {
			tqla, err := New(testCase.options...)
			if err != nil {
				t.Fatal(err)
			}

			s, a, err := tqla.Compile(
				testCase.templateSql, testCase.data)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, testCase.expectedSql, s)
			assert.Equal(t, testCase.expectedArgs, a)
		})
	}
}
