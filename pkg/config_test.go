package pkg_test

import (
	"github.com/blackbass1988/access_logs_stats/pkg"
	"testing"
)

func TestYamlConfig(t *testing.T) {
	filepath := "../config.yaml.example"

	templateVars := make(map[string]string)

	templateVars["fn"] = "foo.txt"

	config, err := pkg.NewConfig(filepath, templateVars)

	if err != nil {
		t.Fatal(err)
	}

	testConfig(t, config)
}

func TestJsonConfig(t *testing.T) {
	filepath := "../config.json.example"

	config, err := pkg.NewConfig(filepath, nil)

	if err != nil {
		t.Fatal(err)
	}
	testConfig(t, config)

}

func testConfig(t *testing.T, config *pkg.Config) {

	if config.InputDsn != "file:foo.txt" {
		t.Error(
			"expected file:foo.txt",
			"actual ",
			config.InputDsn,
		)
	}

	if config.Rex.String() == "" {
		t.Error(
			"expected != \"\"",
			"actual ",
			config.Rex.String(),
		)
	}

	if config.Period.String() != "10s" {
		t.Error(
			"expected 10s",
			"actual ",
			config.Period.String(),
		)
	}

	if len(config.Outputs) != 3 {
		t.Error(
			"expected 3 senders",
			"actual ",
			len(config.Outputs),
		)
	}

	for _, output := range config.Outputs {

		if output.Type == "" {
			t.Error("type must != '' but was ''")
		}

		if output.Settings == nil {
			t.Error("settings was nil")
		}

	}

	if len(config.Filters) != 2 {
		t.Error("filters count expected 2 but was ", len(config.Filters))
		t.FailNow()
	}

	checkFirstFilter(t, config.Filters)
	checkConfig(t, config)
}

func checkFirstFilter(t *testing.T, f []*pkg.Filter) {
	if f[0].String() != ".+" {
		t.Errorf("filter. Expected [.+] . Actual [%s]", f[0].String())
	}

	if f[0].Prefix != "prefix2_" {
		t.Error("filter Prefix. Expected prefix2_ . Actual ", f[0].Prefix)
	}

	if len(f[0].Items) != 2 {
		t.Error("filter Items. Expected 2. Actual ", len(f[0].Items))
	}

	oneItem := f[0].Items[0]
	if oneItem.Field != "code" {
		t.Error("first filter item must be code but was ", oneItem.Field)
	}

	if len(oneItem.Metrics) != 4 {
		t.Error("len(oneItem.Metrics) must be 3 but was ", len(oneItem.Metrics))
	}
}
func checkSecondFilter(t *testing.T, f []*pkg.Filter) {
	if f[1].String() != "/api/v1/" {
		t.Errorf("filter. Expected [/api/v1/] . Actual [%s]", f[1].String())
	}

	if f[1].Prefix != "nginx_requests_" {
		t.Error("filter Prefix. Expected nginx_requests_ . Actual ", f[1].Prefix)
	}

	if len(f[1].Items) != 1 {
		t.Error("filter Items. Expected 1. Actual ", len(f[1].Items))
	}

	firstItem := f[1].Items[0]
	if firstItem.Field != "time" {
		t.Error("first filter item must be code but was ", firstItem.Field)
	}

	if len(firstItem.Metrics) != 1 {
		t.Error("len(firstItem.Metrics) must be 1 but was ", len(firstItem.Metrics))
	}
	if firstItem.Metrics[0] != "prometheus_histogram" {
		t.Error("len(firstItem.Metrics) must be prometheus_histogram but was ", firstItem.Metrics[0])
	}
}

func checkConfig(t *testing.T, config *pkg.Config) {

	if len(config.Counts) != 1 {
		t.Error("config.Counts. Expected 2. Actual ", len(config.Counts))
	}

	if _, ok := config.Counts["code"]; !ok {
		t.Error("config.Counts[code]. Expected code. Actual ", config.Counts["code"])
	}

	if len(config.Aggregates) != 1 {
		t.Error("config.Aggregates. Expected 1. Actual ", len(config.Aggregates))
	}

	if _, ok := config.Aggregates["time"]; !ok {
		t.Error("config.Aggregates[time]. Expected time. Actual ", config.Aggregates["time"])
	}
}
