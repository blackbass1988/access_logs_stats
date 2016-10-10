package core_test

import (
	"github.com/blackbass1988/access_logs_stats/core"
	"testing"
)

func TestConfig(t *testing.T) {
	filepath := "../config.json.example"

	config, err := core.NewConfig(filepath)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

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

	if len(config.Outputs) != 2 {
		t.Error(
			"expected 2 senders",
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

	if len(config.Filters) != 1 {
		t.Error("filters count expected 1 but was ", len(config.Filters))
		t.FailNow()
	}

	f := config.Filters[0]

	if f.Filter != ".+" {
		t.Error("filter. Expected .+ . Actual ", f.Filter)
	}

	if f.Prefix != "prefix2_" {
		t.Error("filter Prefix. Expected prefix2_ . Actual ", f.Prefix)
	}

	if len(f.Items) != 2 {
		t.Error("filter Items. Expected 2. Actual ", len(f.Items))
	}

	oneItem := f.Items[0]
	if oneItem.Field != "code" {
		t.Error("first filter item must be code but was ", oneItem.Field)
	}

	if len(oneItem.Metrics) != 3 {
		t.Error("len(oneItem.Metrics) must be 3 but was ", len(oneItem.Metrics))
	}

	if len(config.Counts) != 2 {
		t.Error("config.Counts. Expected 2. Actual ", len(config.Counts))
	}

	if _, ok := config.Counts["code"]; !ok {
		t.Error("config.Aggregates[code]. Expected code. Actual ", config.Counts["code"])
	}

	if _, ok := config.Counts["time"]; !ok {
		t.Error("config.Aggregates[time]. Expected time. Actual ", config.Counts["time"])
	}

	if len(config.Aggregates) != 1 {
		t.Error("config.Counts. Expected 1. Actual ", len(config.Aggregates))
	}

	if _, ok := config.Aggregates["time"]; !ok {
		t.Error("config.Aggregates[time]. Expected time. Actual ", config.Aggregates["time"])
	}
}
