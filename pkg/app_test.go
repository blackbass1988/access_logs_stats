package pkg_test

import (
	"log"
	"testing"

	"github.com/blackbass1988/access_logs_stats/pkg"
	"github.com/blackbass1988/access_logs_stats/pkg/re"
)

type testcase struct {
	regexp    string
	rawString string
	fields    map[string]string
}

var tests = []testcase{
	{
		`.+HTTP/\d.?\d?\s(?P<code>\d+)[^"]+"[^"]*" "[^"]*" (?P<time>\d{1,}\.\d{3})`,
		"Sep 18 13:47:55 squid1.drom.ru nginx: s.auto.drom.ru s.auto.drom.ru 192.168.200.86 - " +
			"[2016-09-18T13:47:55+10:00] GET \"/4/sales/photos/15844/15843967/129148731.jpg\" HTTP/1.1 200 85972 " +
			"\"-\" \"Go-http-client/1.1\" 0.132 159 \"192.168.22.171:7483\" \"0.132\" MISS \"-\" - -",
		map[string]string{"code": "200", "time": "0.132"},
	},
	{
		`.+HTTP/\d.?\d?\s(?P<code>\d+)\s(?P<bytes>\d+)[^"]+"[^"]*" "[^"]*" (?P<time>\d{1,}\.\d{3})`,
		`Sep 19 20:22:41 node4.drom.ru nginx: s.auto.drom.ru s.drom.ru 85.26.241.25 - [2016-09-19T20:22:41+10:00] GET "/1/reviews/photos/mitsubishi/gto/gen2154x2_1519_0.jpeg" HTTP/1.1 200 98389 "http://www.drom.ru/reviews/mitsubishi/gto/" "Mozilla/5.0 (Linux; Android 5.0.2; SM-A300F Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.81 Mobile Safari/537.36" 47.400 722 "192.168.200.86:80" "0.000" MISS "-" b593b3b588365341b5554d0486c6b880 -`,
		map[string]string{"code": "200", "time": "47.400", "bytes": "98389"},
	},
	{
		`.+HTTP/\d.?\d?\s(?P<code>\d+)[^"]+"[^"]*" "[^"]*" (?P<time>\d{1,}\.\d{3})`,
		`Sep 18 15:46:11 squid1.drom.ru nginx: s.auto.drom.ru s.auto.drom.ru 95.58.173.133 - ` +
			`[2016-09-18T15:46:11+10:00] GET "/1/sales/photos/12157/12156267/ttn_160_90971946.jpg" ` +
			`HTTP/1.1 403 78 "http://novosibirsk.drom.ru/toyota/land_cruiser/12156267.html" "Mozilla/5.0 ` +
			`(Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36" ` +
			`0.000 747 "192.168.22.171:7482" "0.000" MISS "-" ec5801arPcHw3p4AHrwnUotLtUZjQ0a4 -`,
		map[string]string{"code": "403", "time": "0.000"},
	},
	{
		`(?P<val1>\S{32,}) (\S+)$`,
		`Oct 11 19:50:41 node3.drom.ru nginx: asterisk.drom.ru habarovsk.drom.ru 62.249.146.88 - [2016-10-11T19:50:41+10:00] GET "/nissan/cedric/23791224.html" HTTP/1.1 200 7632 "http://auto.drom.ru/region27/page2/?minyear=1995&minprice=100000&maxprice=180000&order=enterdate&order_d=desc&go_search=2" "Mozilla/5.0 (Linux; Android 4.4.2; ASUS_T00I Build/KVT49L) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.124 Mobile Safari/537.36" 0.600 1283 "192.168.200.184:9000" "0.600" - "-" 6b80e2cLttrt7c0ZBRRMFObguMU%2BA0a0 15306023%3A2ef31a9096b7635c9457b96a5389cee8`,
		map[string]string{"val1": "6b80e2cLttrt7c0ZBRRMFObguMU%2BA0a0"},
	},
}

func TestRegexp(t *testing.T) {

	for _, test := range tests {

		rex := re.MustCompile(test.regexp)

		actualRow, err := pkg.NewRow(test.rawString, rex)
		check_test(err, t)

		if actualRow.Raw != test.rawString {
			t.Errorf("expected [%s] actual [%s]", test.rawString, actualRow.Raw)
		}

		log.Print(actualRow.Fields)

		for actualField, actualValue := range actualRow.Fields {
			if test.fields[actualField] != actualValue {
				t.Errorf("field [%s] expected [%s] actual [%s]", actualField, test.fields[actualField], actualValue)
			}
		}
	}
}

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
