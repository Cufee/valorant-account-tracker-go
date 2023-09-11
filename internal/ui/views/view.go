package views

import (
	"bytes"
	t "text/template"
)

type View struct {
	template *t.Template
}

func (v *View) Render(props any) (string, error) {
	var tpl bytes.Buffer
	if err := v.template.Execute(&tpl, props); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func newView(name, template string, children ...View) (View, error) {
	parsedTemplate, err := t.New(name).Parse(template)
	for _, child := range children {
		_, err := parsedTemplate.AddParseTree(child.template.Name(), child.template.Tree)
		if err != nil {
			return View{}, err
		}
	}
	if err != nil {
		return View{}, err
	}

	return View{
		template: parsedTemplate,
	}, nil
}
