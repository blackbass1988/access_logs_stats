package output

import (
	"errors"
	"fmt"
	"strings"
)

var requiredFields = [2]string{"field", "metric"}
var varTemplate = "${%s}"

type Template struct {
	template string
}

// Create string from template with input parameters
func (template *Template) Process(field string, metric string, templateVars map[string]string) (error, string) {

	finalString := template.template

	if templateVars == nil {
		templateVars = make(map[string]string)
	}

	templateVars["field"] = field
	templateVars["metric"] = metric

	for f := range templateVars {
		replaceString := fmt.Sprintf(varTemplate, f)

		if !strings.Contains(finalString, replaceString) {
			return errors.New("field \"" + f + "\" not found in template " + template.template), ""
		}

		finalString = strings.ReplaceAll(finalString, replaceString, templateVars[f])
	}

	return nil, finalString
}

func NewTempate(template string) (error, *Template) {
	var err error

	t := new(Template)

	if err = validateTemplate(template); err != nil {
		return err, nil
	}

	t.template = template

	return nil, t
}

func validateTemplate(inputString string) error {

	for _, field := range requiredFields {
		if !strings.Contains(inputString, "${"+field+"}") {
			return errors.New(fmt.Sprintf("\"${%s}\" not found in [%s]", field, inputString))
		}
	}

	return nil
}
