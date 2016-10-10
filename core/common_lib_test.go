package core_test

import (
	"testing"
	//"regexp"
	//"github.com/blackbass1988/access_logs_stats/core"
	//"time"
	"log"
)

type testcase struct {
	rawString            string
	expectedDate         string
	expectedStatus       int
	expectedResponseTime float64
}

////todo reimplement this test to new system
//var tests = []testcase{
//	{
//		"Sep 18 13:47:55 squid1.drom.ru nginx: s.auto.drom.ru s.auto.drom.ru 192.168.200.86 - " +
//			"[2016-09-18T13:47:55+10:00] GET \"/4/sales/photos/15844/15843967/129148731.jpg\" HTTP/1.1 200 85972 " +
//			"\"-\" \"Go-http-client/1.1\" 0.132 159 \"192.168.22.171:7483\" \"0.132\" MISS \"-\" - -",
//		"Sep 18 13:47:55",
//		200,
//		0.132,
//	},
//	{`Sep 19 20:22:41 node4.drom.ru nginx: s.auto.drom.ru s.drom.ru 85.26.241.25 - [2016-09-19T20:22:41+10:00] GET "/1/reviews/photos/mitsubishi/gto/gen2154x2_1519_0.jpeg" HTTP/1.1 200 98389 "http://www.drom.ru/reviews/mitsubishi/gto/" "Mozilla/5.0 (Linux; Android 5.0.2; SM-A300F Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.81 Mobile Safari/537.36" 47.400 722 "192.168.200.86:80" "0.000" MISS "-" b593b3b588365341b5554d0486c6b880 -`,
//		"Sep 19 20:22:41",
//		200,
//		47.4,
//	},
//	{
//		`Sep 18 15:46:11 squid1.drom.ru nginx: s.auto.drom.ru s.auto.drom.ru 95.58.173.133 - ` +
//			`[2016-09-18T15:46:11+10:00] GET "/1/sales/photos/12157/12156267/ttn_160_90971946.jpg" ` +
//			`HTTP/1.1 403 78 "http://novosibirsk.drom.ru/toyota/land_cruiser/12156267.html" "Mozilla/5.0 ` +
//			`(Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36" ` +
//			`0.000 747 "192.168.22.171:7482" "0.000" MISS "-" ec5801arPcHw3p4AHrwnUotLtUZjQ0a4 -`,
//		"Sep 18 15:46:11",
//		403,
//		0.0,
//	},
//	{
//		`Sep 19 21:22:32 squid2.drom.ru nginx: s.auto.drom.ru s.auto.drom.ru 176.59.141.154 - [2016-09-19T21:22:32+10:00] GET "/i24200/s/photos/23075/23074726/180719837.jpg" HTTP/1.1 200 138011 "http://abakan.drom.ru/lada/2106/23074726.html" "" 0.000 576 "-" "-" HIT "-" b3524a3l3GyZI0FjSkdSaq%2FXP9l%2Bw0ac -`,
//		"Sep 19 21:22:32",
//		200,
//		0.0,
//	},
//}
//
//func TestRegexp(t *testing.T) {
//	rex, err := regexp.Compile(`(?P<date>\w+\s\w+\s[0-9:]+).+HTTP/\d.?\d?\s(?P<code>\d+)[^"]+"[^"]*" "[^"]*" (?P<time>\d{1,}\.\d{3})`)
//	check_test(err, t)
//
//	for _, test := range tests {
//
//		actualRow, err := core.NewRow(rex, test.rawString)
//		check_test(err, t)
//
//		expectedDate, err := time.Parse(time.Stamp, test.expectedDate)
//
//		check_test(err, t)
//
//		if !expectedDate.Equal(actualRow.Date) {
//			t.Error(
//				"expected date ", expectedDate,
//				"actual", actualRow.Date,
//			)
//		}
//		if test.expectedStatus != actualRow.StatusCode {
//			t.Error(
//				"expected status ", test.expectedStatus,
//				"actual", actualRow.StatusCode,
//			)
//		}
//
//		if test.expectedResponseTime != actualRow.ResponseTime {
//			t.Error(
//				"expected response time ", test.expectedResponseTime,
//				"actual", actualRow.ResponseTime,
//			)
//		}
//
//	}
//}

func TestStatusCodes(t *testing.T) {
	statusCodes := map[int]uint64{}
	statusCodes[200]++
	statusCodes[200]++
	statusCodes[200]++
	statusCodes[400]++

	for code, cnt := range statusCodes {

		if code == 200 && cnt != 3 {
			t.Error(code, cnt)
		}

		if code == 400 && cnt != 1 {
			t.Error(code, cnt)
		}
	}
}

func check_test(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestFoo(t *testing.T) {
	var counts map[string]map[string]uint64
	counts = make(map[string]map[string]uint64)
	counts["zz"] = make(map[string]uint64)
	counts["zz"]["zzzz"]++
	log.Println(counts)
}
