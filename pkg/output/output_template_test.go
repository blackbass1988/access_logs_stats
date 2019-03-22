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
}


func TestGoodTemplate(t *testing.T) {
	metric := "cps_200"
	field := "count"
	expectedString := "count.cps_200"

	err, template := output.NewTempate("${field}.${metric}")
	if err != nil {
		t.Errorf("Error must be nil, [%s] was", err)
	}

	err, actualString := template.Process(field, metric, nil)

	if expectedString != actualString {
		t.Errorf("String must be [%s], [%s] was", expectedString, actualString)
	}
}


func TestGoodTemplateWithoutPayloadButInTemplateItWas(t *testing.T) {
	metric := "cps_200"
	field := "count"

	err, template := output.NewTempate("${field}.${metric}:${payload[field]}")

	err, _ = template.Process(field, metric, nil)

	if err == nil {
		t.Errorf("Must be error, but was nil")
	}
}

func TestWithPayload(t *testing.T) {
	expectedString := "cps_200[localhost,count]"
	metric := "cps_200"
	field := "count"
	payload := make(map[string]string)
	payload["host"] = "localhost"


	err, template := output.NewTempate ("${metric}[${payload[host]},${field}]")
	if err != nil {
		t.Errorf("Error must be nil, [%s] was", err)
	}

	err, actualString := template.Process(field, metric, payload)

	if expectedString != actualString {
		t.Errorf("String must be [%s], [%s] was", expectedString, actualString)
	}

}
