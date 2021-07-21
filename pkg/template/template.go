package template

import (
	"fmt"
	"strings"
)

var varTemplate = "${%s}"

type Template struct {
	template string
}

// Create string from template with input parameters
func (t *Template) ProcessTemplate(tempVars map[string]string) string {
	finalString := t.template

	for k, v := range tempVars {
		replaceString := fmt.Sprintf(varTemplate, k)

		finalString = strings.ReplaceAll(finalString, replaceString, v)
	}

	return finalString
}

// Create string from template with input parameters
func (t *Template) Process(field string, metric string, templateVars map[string]string) string {

	tempVars := make(map[string]string)

	if templateVars != nil {

		for k, v := range templateVars {
			tempVars[k] = v
		}
	}

	tempVars["field"] = field
	tempVars["metric"] = metric

	return t.ProcessTemplate(tempVars)
}

func NewTemplate(templateString string) *Template {

	t := new(Template)

	t.template = templateString

	return t
}
