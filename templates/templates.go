package templates

import (
	"text/template"
)

var tmpl *template.Template

func Get() *template.Template {
	return tmpl
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*"))
}
