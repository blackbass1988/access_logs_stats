package template_test

import (
	"github.com/blackbass1988/access_logs_stats/pkg/template"
	"testing"
)

func TestGoodTemplate(t *testing.T) {
	field := "count"
	metric := "cps_200"
	expectedString := "count.cps_200"

	err, tmpl := template.NewTempate("${field}.${metric}")
	if err != nil {
		t.Errorf("Error must be nil, was: [%s]", err.Error())
		t.FailNow()
	}

	err, actualString := tmpl.Process(field, metric, nil)

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

	err, tmpl := template.NewTempate("${metric}[${hostname},${field}]")
	if err != nil {
		t.Errorf("Error must be nil, [%s] was", err)
	}

	err, actualString := tmpl.Process(field, metric, templateVars)

	if expectedString != actualString {
		t.Errorf("String must be [%s], [%s] was", expectedString, actualString)
	}
}
