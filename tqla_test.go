package tqla

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTqla(t *testing.T) {
	testCases := []struct {
		templateSql  string
		placeholder  Placeholder
		data         any
		expectedSql  string
		expectedArgs []any
	}{
		{
			templateSql: `{{$column:=true -}}
			select * from table where v = {{.Value }} 
			{{if $column -}}
			and c = {{$column -}}
			{{end}}`,
			data: struct {
				Value int
			}{Value: 5},
			placeholder: Dollar,
			expectedSql: `select * from table where v = $1 
			and c = $2`,
			expectedArgs: []any{
				5,
				true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run("test", func(t *testing.T) {
			tqla, err := New(WithPlaceHolder(testCase.placeholder))
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
