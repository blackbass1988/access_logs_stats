package output_test

import (
	"github.com/blackbass1988/access_logs_stats/pkg/output"
	"testing"
)

func TestTemplateWithoutVars(t *testing.T) {
	err, _ := output.NewTempate ("abcdef")
	if err == nil {
		t.Errorf("Must be error, but was nil")
	}
}

func TestTemplateWithoutRequiredVars(t *testing.T) {
	err, _ := output.NewTempate ("${field}.hello.world")
	if err == nil {
		t.Errorf("Must be error, but was nil")
	}

	if err != nil && err.Error() != "\"${metric}\" not found in [${field}.hello.world]" {
		t.Errorf("Invalid error: [%s]", err)
	}

}


func TestGoodTemplate(t *testing.T) {
	field := "count"
	metric := "cps_200"
	expectedString := "count.cps_200"

	err, template := output.NewTempate("${field}.${metric}")
	if err != nil {
		t.Errorf("Error must be nil, was: [%s]", err.Error())
		t.FailNow()
	}

	err, actualString := template.Process(field, metric, nil)

	if expectedString != actualString {
		t.Errorf("String must be >>%s<<, was: >>%s<<", expectedString, actualString)
	}
}

func TestWithPayload(t *testing.T) {
	expectedString := "cps_200[localhost,count]"
	metric := "cps_200"
	field := "count"
	payload := make(map[string]string)
	payload["hostname"] = "localhost"


	err, template := output.NewTempate ("${metric}[${hostname},${field}]")
	if err != nil {
		t.Errorf("Error must be nil, [%s] was", err)
	}

	err, actualString := template.Process(field, metric, payload)

	if expectedString != actualString {
		t.Errorf("String must be [%s], [%s] was", expectedString, actualString)
	}
}
