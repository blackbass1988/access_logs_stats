package template_test

import (
	"github.com/blackbass1988/access_logs_stats/pkg/template"
	"testing"
)

func TestGoodTemplate(t *testing.T) {
	field := "count"
	metric := "cps_200"
	expectedString := "count.cps_200"

	tmpl := template.NewTemplate("${field}.${metric}")

	actualString := tmpl.Process(field, metric, nil)

	if expectedString != actualString {
		t.Errorf("String must be >>%s<<, was: >>%s<<", expectedString, actualString)
	}
}

func TestWithTemplateVars(t *testing.T) {
	expectedString := "cps_200[localhost,count]"
	metric := "cps_200"
	field := "count"
	templateVars := make(map[string]string)
	templateVars["hostname"] = "localhost"

	tmpl := template.NewTemplate("${metric}[${hostname},${field}]")

	actualString := tmpl.Process(field, metric, templateVars)

	if expectedString != actualString {
		t.Errorf("String must be [%s], [%s] was", expectedString, actualString)
	}
}
