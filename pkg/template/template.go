package template

import (
	"errors"
	"fmt"
	"strings"
)

var varTemplate = "${%s}"

type Template struct {
	template string
}

// Create string from template with input parameters
func (t *Template) ProcessTemplate(tempVars map[string]string, errorOnNoField bool) (error, string) {
	finalString := t.template

	for k, v := range tempVars {
		replaceString := fmt.Sprintf(varTemplate, k)

		if errorOnNoField && !strings.Contains(finalString, replaceString) {
			return errors.New("field \"" + k + "\" not found in t " + t.template), ""
		}

		finalString = strings.ReplaceAll(finalString, replaceString, v)
	}

	return nil, finalString
}

// Create string from template with input parameters
func (t *Template) Process(field string, metric string, templateVars map[string]string) (error, string) {

	tempVars := make(map[string]string)

	if templateVars != nil {

		for k, v := range templateVars {
			tempVars[k] = v
		}
	}

	tempVars["field"] = field
	tempVars["metric"] = metric

	return t.ProcessTemplate(tempVars, true)
}

func NewTempate(templateString string) (error, *Template) {

	t := new(Template)

	t.template = templateString

	return nil, t
}
