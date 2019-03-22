package pkg

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/blackbass1988/access_logs_stats/pkg/output"
)

//Sender sends data to output. omg omg omg
type Sender struct {
	filter *Filter

	//ссылка на конфиг приложения
	config *Config

	//настроенный отправлятор, реализации настраиваются в конфиге "outputs"
	output *output.Output

	//мап флоатов с реализацией агрегирующих фунций
	floatData map[string]*Float64Data

	//закешированное кол-во секунд в периоде, указанному в "period" конфигурации
	periodInSeconds float64

	//здесь хранятся числа по полям, указанные в "aggregates" конфигурации
	floatsForAggregates map[string][]float64

	//здесь хранятся счетчики уникальных значений по полям, указанные в "counts" конфигурации
	//хранится по схеме поле.уник_значение.кол-во
	//вывод происходит по схеме - кол-во в 1 секунду
	counts map[string]map[string]uint64

	globalLock sync.Mutex
}

func (s *Sender) resetData() {
	s.globalLock.Lock()
	s.floatsForAggregates = make(map[string][]float64)
	s.floatData = make(map[string]*Float64Data)
	s.counts = make(map[string]map[string]uint64)
	s.globalLock.Unlock()
}

func (s *Sender) appendIfOk(row *RowEntry) (err error) {

	if s.filter.MatchString(row.Raw) {

		for field, val := range row.Fields {

			if _, ok := s.config.Aggregates[field]; ok {
				valFloat, err := strconv.ParseFloat(val, 10)
				checkOrFail(err)
				s.floatsForAggregates[field] = append(s.floatsForAggregates[field], valFloat)
			}

			//в конфиге указано поле, как поле, по которому считаются
			// суммы по уникальным значениям
			if _, ok := s.config.Counts[field]; ok {
				if s.counts[field] == nil {
					s.counts[field] = make(map[string]uint64)
				}
				s.counts[field][val]++
			}
		}
	}

	return err
}

func (s *Sender) sendStats() (err error) {

	s.globalLock.Lock()
	for _, metricsOfField := range s.filter.Items {

		for _, metric := range metricsOfField.Metrics {
			s.appendToOutput(metricsOfField.Field, metric)
		}
	}
	s.globalLock.Unlock()
	s.output.Send()

	return err
}

//NewSender create new sender
func NewSender(filter *Filter, config *Config) (*Sender, error) {
	sender := new(Sender)
	sender.filter = filter
	sender.config = config

	sender.output = new(output.Output)

	if len(filter.Prefix) > 0 {
		sender.output.SetPrefix(filter.Prefix)
	}

	for _, s := range config.Outputs {
		sender.output.Init(s.Type, s.Settings)
	}

	return sender, nil
}

func (s *Sender) appendToOutput(field string, metric string) {
	var (
		key             string
		periodInSeconds float64
		value           string
	)

	key = fmt.Sprintf("%s_%v", field, metric)

	switch {
	case metric == "min":
		value = fmt.Sprintf("%.3f", s.getFloatData(field).Min())
	case metric == "max":
		value = fmt.Sprintf("%.3f", s.getFloatData(field).Max())
	case metric == "len":
		value = fmt.Sprintf("%d", s.getFloatData(field).Len())
	case metric == "avg":
		value = fmt.Sprintf("%.3f", s.getFloatData(field).Avg())
	case metric == "sum":
		value = fmt.Sprintf("%.3f", s.getFloatData(field).Sum())
	case metric == "sum_ps":
		result := s.getFloatData(field).Sum()
		periodInSeconds = s.getPeriodInSeconds()
		if periodInSeconds == 0 {
			result = 0
		} else {
			result = s.getFloatData(field).Sum() / s.getPeriodInSeconds()
		}

		value = fmt.Sprintf("%.3f", result)
	case metric == "ips":
		value = fmt.Sprintf("%.3f", s.getFloatData(field).ItemsPerSeconds(s.getPeriodInSeconds()))
	case strings.Contains(metric, "cent_"):
		cent := strings.Split(metric, "_")
		centFloat, err := strconv.ParseFloat(cent[1], 10)
		checkOrFail(err)
		value = fmt.Sprintf("%.3f", s.getFloatData(field).Percentile(centFloat))
	case metric == "uniq":
		value = fmt.Sprintf("%d", s.getUniqCnt(field))
	case metric == "uniq_ps":
		value = fmt.Sprintf("%.3f", float64(s.getUniqCnt(field))/s.getPeriodInSeconds())
	case strings.Contains(metric, "cps_"):
		value = s.processCps(metric, field)
	case strings.Contains(metric, "percentage_"):
		value = s.processPercentage(metric, field)
	}
	s.output.AddMessage(field, metric, value)
}

func (s *Sender) processCps(metric string, field string) string {
	var (
		ok  bool
		cnt uint64
	)

	metrics := strings.Split(metric, "_")
	metric = metrics[1]
	if _, ok = s.counts[field]; !ok {
		cnt = 0
	} else if cnt, ok = s.counts[field][metric]; !ok {
		cnt = 0
	}
	return fmt.Sprintf("%.3f", float64(cnt)/s.getPeriodInSeconds())
}

func (s *Sender) processPercentage(metric string, field string) string {
	var (
		ok     bool
		cnt    uint64
		result float64
	)

	metrics := strings.Split(metric, "_")
	metric = metrics[1]

	total := s.getTotalCountByField(field)

	result = 0
	if cnt, ok = s.counts[field][metric]; ok && total > 0 {
		result = float64(cnt * 100 / total)
	}
	return fmt.Sprintf("%.3f", result)
}

func (s *Sender) getTotalCountByField(field string) uint64 {
	var (
		ok  bool
		cnt uint64
	)
	if _, ok = s.counts[field]; !ok {
		return 0
	}

	cnt = 0
	for _, c := range s.counts[field] {
		cnt += c
	}

	return cnt
}

func (s *Sender) getUniqCnt(field string) uint64 {
	var (
		cnt uint64
		ok  bool
	)
	cnt = 0
	if _, ok = s.counts[field]; ok {
		cnt = uint64(len(s.counts[field]))
	}
	return cnt
}
func (s *Sender) getFloatData(field string) *Float64Data {
	//кешируем флоатдату
	if _, ok := s.floatData[field]; !ok {
		f := Float64Data(s.floatsForAggregates[field])
		s.floatData[field] = &f
		sort.Sort(s.floatData[field])
	}
	return s.floatData[field]
}

func (s *Sender) getPeriodInSeconds() float64 {
	if s.periodInSeconds == 0 {
		s.periodInSeconds = s.config.Period.Seconds()
	}
	return s.periodInSeconds
}
