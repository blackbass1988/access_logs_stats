package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type OutputConfig struct {
	Type     string            `json:"type"`
	Settings map[string]string `json:"settings"`
}

type configJson struct {
	InputDsn   string   `json:"input"`
	Regexp     string   `json:"regexp"`
	Period     string   `json:"period"`
	Counts     []string `json:"counts"`
	Aggregates []string `json:"aggregates"`

	Filters []*Filter       `json:"filters"`
	Outputs []*OutputConfig `json:"output"`
}

type Config struct {
	InputDsn string

	ExitAfterOneTick bool

	Counts     map[string]bool
	Aggregates map[string]bool

	Outputs []*OutputConfig
	Rex     *regexp.Regexp
	Period  time.Duration
	Filters []*Filter
}

func NewConfig(filepath string) (config Config, err error) {
	configJson := new(configJson)
	config.Aggregates = make(map[string]bool)
	config.Counts = make(map[string]bool)

	//we need to lock file in processlist for restore by file descriptor if delete in runtime
	_, err = os.Open(filepath)
	if err != nil {
		return config, err
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &configJson)
	if err != nil {
		return config, err
	}

	config.InputDsn = configJson.InputDsn
	config.Period, err = time.ParseDuration(configJson.Period)
	if err != nil {
		return config, err
	}

	config.Rex, err = regexp.Compile(configJson.Regexp)
	if err != nil {
		return config, err
	}

	config.Outputs = configJson.Outputs

	for _, el := range configJson.Counts {
		config.Counts[el] = true
	}

	for _, el := range configJson.Aggregates {
		config.Aggregates[el] = true
	}

	for _, f := range configJson.Filters {
		//f.FilterRex, err = regexp.Compile(f.matcher)

		for _, filterItem := range f.Items {
			for _, metric := range filterItem.Metrics {

				switch {
				case metric == "min", metric == "max", metric == "len", metric == "avg",
					metric == "sum", metric == "sum_ps", metric == "ips", strings.Contains(metric, "cent_"):

					if !config.Aggregates[filterItem.Field] {
						err = errors.New(
							fmt.Sprintf("field \"%s\" must in in \"aggregates\" section"+
								" because you want metric \"%s\"",
								filterItem.Field, metric))
					}
				case metric == "uniq", metric == "uniq_ps", strings.Contains(metric, "cps_"), strings.Contains(metric, "percentage_"):
					if !config.Counts[filterItem.Field] {
						err = errors.New(
							fmt.Sprintf("field \"%s\" must in in \"counts\" section "+
								"because you want metric \"%s\"",
								filterItem.Field, metric))
					}
				}

			}
		}

		check(err)
		config.Filters = append(config.Filters, f)
	}

	if len(config.Filters) == 0 {
		err = errFiltersNotSet
	}

	if len(config.Outputs) == 0 {
		return config, errOutputNotSet
	}

	return config, err
}
