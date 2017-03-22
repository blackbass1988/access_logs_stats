package input

import (
	"testing"
)

type testcaseparser struct {
	msg      string
	priority string
	date     string
	hostname string
	app      string
	message  string
}

var testsparser = []testcaseparser{
	{
		`<5>Oct  6 17:56:56 zzz app1: zzz`,
		"5",
		"Oct  6 17:56:56",
		"zzz",
		"app1",
		"zzz",
	},
	{
		`<149>Oct  7 13:51:20 node2.drom.ru nginx: s.auto.drom.ru s.rdrom.ru 31.173.227.172 - [2016-10-07T13:51:20+10:00] GET "/1/catalog/photos/generations/toyota_passo_g779.jpg?17911" HTTP/1.1 200 4333 "http://www.drom.ru/catalog/toyota/passo/" "Mozilla/5.0 (iPhone; CPU iPhone OS 9_3_5 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13G36 Safari/601.1" 0.000 381 "-" "-" HIT "-" - -56 Safari/602.1" 0.000 644 "-" "-" HIT "-" f87a666AIsfgMIILIrOCyVLlE10KQ0ab -0a7 19665144%3Aabea9d3121a0dba90d9279d31648c8fd`,
		"149",
		"Oct  7 13:51:20",
		"node2.drom.ru",
		"nginx",
		`s.auto.drom.ru s.rdrom.ru 31.173.227.172 - [2016-10-07T13:51:20+10:00] GET "/1/catalog/photos/generations/toyota_passo_g779.jpg?17911" HTTP/1.1 200 4333 "http://www.drom.ru/catalog/toyota/passo/" "Mozilla/5.0 (iPhone; CPU iPhone OS 9_3_5 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13G36 Safari/601.1" 0.000 381 "-" "-" HIT "-" - -56 Safari/602.1" 0.000 644 "-" "-" HIT "-" f87a666AIsfgMIILIrOCyVLlE10KQ0ab -0a7 19665144%3Aabea9d3121a0dba90d9279d31648c8fd`,
	},
}

func TestParseSyslogMsg(t *testing.T) {

	parser, _ := newSyslogParser()

	for _, test := range testsparser {

		syslogM, _ := parser.parseSyslogMsg(test.msg)

		if test.priority != syslogM.Priority {
			t.Error("priority ",
				"expected ", test.priority,
				"actual ", syslogM.Priority,
			)
		}

		if test.date != syslogM.Date {
			t.Error("date ",
				"expected ", test.date,
				"actual ", syslogM.Priority,
			)
		}

		if test.app != syslogM.Application {
			t.Error("Application ",
				"expected ", test.app,
				"actual ", syslogM.Application,
			)
		}

		if test.message != syslogM.Message {
			t.Error("Message ",
				"expected [", test.message, "]",
				"actual [", syslogM.Message, "]",
			)
		}

	}

}
