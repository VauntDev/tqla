package tqla

import (
	"fmt"
	"io"
	"log"
	"text/template"
	"text/template/parse"
)

type sqlTemplate struct {
	text *template.Template
}

func newSqlTemplate(name string, funcs template.FuncMap) *sqlTemplate {
	return &sqlTemplate{
		text: template.New("name").Funcs(funcs),
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
	declaredVars := make(map[string]string)
	for _, tmpl := range t.Templates() {
		formatTree(tmpl.Tree, declaredVars)
	}
}

func formatTree(t *parse.Tree, declared map[string]string) *parse.Tree {
	if t.Root == nil {
		return t
	}
	formatNode(t, t.Root, declared)
	return t
}

func formatNode(t *parse.Tree, n parse.Node, declared map[string]string) {
	switch v := n.(type) {
	case *parse.ActionNode:
		formatNode(t, v.Pipe, declared)
	case *parse.IfNode:
		formatNode(t, v.List, declared)
		formatNode(t, v.ElseList, declared)
	case *parse.RangeNode:
		if len(v.Pipe.Decl) == 2 {
			declared[v.Pipe.Decl[0].String()] = v.Pipe.Decl[0].String()
		}
		formatNode(t, v.List, declared)
		formatNode(t, v.ElseList, declared)
	case *parse.ListNode:
		if v == nil {
			return
		}
		for _, n := range v.Nodes {
			formatNode(t, n, declared)
		}
	case *parse.WithNode:
		formatNode(t, v.List, declared)
		formatNode(t, v.ElseList, declared)
	case *parse.PipeNode:
		if len(v.Decl) > 0 {
			for _, d := range v.Decl {
				declared[d.String()] = d.String()
			}
			// If the pipe sets variables then don't try to format it
			return
		}

		if len(v.Cmds) < 1 {
			return
		}

		cmd := v.Cmds[len(v.Cmds)-1]
		if _, ok := declared[cmd.String()]; ok {
			log.Println("here")
			return
		}

		log.Println("decl", declared)
		if len(cmd.Args) == 1 && cmd.Args[0].Type() == parse.NodeIdentifier && cmd.Args[0].(*parse.IdentifierNode).Ident == "_sql_parser_" {
			return
		}

		v.Cmds = append(v.Cmds, &parse.CommandNode{
			NodeType: parse.NodeCommand,
			Args:     []parse.Node{parse.NewIdentifier("_sql_parser_").SetTree(t).SetPos(cmd.Pos)},
		})
	}
}
