package tqla

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Placeholder defines a interace that exposes a Format function.
// The Format function is intened for sql placeholder formating (i.e ?,$,:,@)
type Placeholder interface {
	Format(sql string) (string, error)
}

var (
	// Question is a PlaceholderFormat instance that leaves placeholders as
	// question marks.
	Question = questionFormat{}

	// Dollar is a PlaceholderFormat instance that replaces placeholders with
	// dollar-prefixed positional placeholders (e.g. $1, $2, $3).
	Dollar = dollarFormat{}

	// Colon is a PlaceholderFormat instance that replaces placeholders with
	// colon-prefixed positional placeholders (e.g. :1, :2, :3).
	Colon = colonFormat{}

	// AtP is a PlaceholderFormat instance that replaces placeholders with
	// "@p"-prefixed positional placeholders (e.g. @p1, @p2, @p3).
	AtP = atpFormat{}

	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// questionFormat is the default format and should be treated as a pass through
type questionFormat struct{}

func (questionFormat) Format(sql string) (string, error) {
	return sql, nil
}

type dollarFormat struct{}

func (dollarFormat) Format(sql string) (string, error) {
	return format(sql, "$")
}

type colonFormat struct{}

func (colonFormat) Format(sql string) (string, error) {
	return format(sql, ":")
}

type atpFormat struct{}

func (atpFormat) Format(sql string) (string, error) {
	return format(sql, "@p")
}

func format(sql, prefix string) (string, error) {
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}
		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			fmt.Fprintf(buf, "%s%d", prefix, i)
			sql = sql[p+1:]
		}
	}
	buf.WriteString(sql)
	s := whitespaceRegex.ReplaceAllString(strings.TrimSpace(buf.String()), " ")
	return s, nil
}
