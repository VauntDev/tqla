package tqla

import (
	"fmt"
	"io"
	"text/template"
	"text/template/parse"
)

type sqlTemplate struct {
	text *template.Template
}

func newSqlTemplate(name string, funcs template.FuncMap) *sqlTemplate {
	return &sqlTemplate{
		text: template.New(name).Funcs(funcs),
	}
}

func (st *sqlTemplate) parse(text string) error {
	if st.text == nil {
		return fmt.Errorf("nil template")
	}
	tmpl, err := st.text.Parse(text)
	if err != nil {
		return err
	}
	formatTemplate(tmpl)
	return nil
}

func (st *sqlTemplate) execute(wr io.Writer, data any) error {
	return st.text.Execute(wr, data)
}

// formatTemplate formats all the templates defined in a template.
func formatTemplate(t *template.Template) {
	for _, tmpl := range t.Templates() {
		formatTree(tmpl.Tree)
	}
}

func formatTree(t *parse.Tree) *parse.Tree {
	if t.Root == nil {
		return t
	}
	formatNode(t, t.Root)
	return t
}

func formatNode(t *parse.Tree, n parse.Node) {
	switch v := n.(type) {
	case *parse.ActionNode:
		formatNode(t, v.Pipe)
	case *parse.IfNode:
		formatNode(t, v.List)
		formatNode(t, v.ElseList)
	case *parse.RangeNode:
		formatNode(t, v.List)
		formatNode(t, v.ElseList)
	case *parse.ListNode:
		if v == nil {
			return
		}
		for _, n := range v.Nodes {
			formatNode(t, n)
		}
	case *parse.WithNode:
		formatNode(t, v.List)
		formatNode(t, v.ElseList)
	case *parse.PipeNode:
		if len(v.Decl) > 0 {
			// If the pipe sets variables then don't try to format it
			return
		}
		if len(v.Cmds) < 1 {
			return
		}
		cmd := v.Cmds[len(v.Cmds)-1]
		if len(cmd.Args) == 1 && cmd.Args[0].Type() == parse.NodeIdentifier && cmd.Args[0].(*parse.IdentifierNode).Ident == "_sql_parser_" {
			return
		}
		v.Cmds = append(v.Cmds, &parse.CommandNode{
			NodeType: parse.NodeCommand,
			Args:     []parse.Node{parse.NewIdentifier("_sql_parser_").SetTree(t).SetPos(cmd.Pos)},
		})
	}
}
