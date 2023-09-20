package tqla

type sqlParser struct {
	args []any
}

func newSqlParser() *sqlParser {
	return &sqlParser{
		args: make([]any, 0),
	}
}

func (s *sqlParser) parsefunc(arg any) string {
	if s.args == nil {
		s.args = make([]any, 0)
	}
	s.args = append(s.args, arg)
	return "?"
}
