package input

import (
	"testing"
)

type testcase struct {
	dsn                 string
	expectedProtocol    string
	expectedListen      string
	expectedApplication string
	expectedError       error
}

var tests = []testcase{
	{
		"syslog:udp:binding_ip:binding_port/application",
		"udp",
		"binding_ip:binding_port",
		"application",
		nil,
	},
	{
		"syslog:tcp:127.0.0.1:1234/application",
		"tcp",
		"127.0.0.1:1234",
		"application",
		nil,
	},
	{
		"syslog:tcp::1234/application",
		"tcp",
		":1234",
		"application",
		nil,
	},

	{
		"syslog:zoo::1234/application",
		"zoo",
		":1234",
		"application",
		ErrorUnknownProtocol,
	},
	{
		"syslog:tcp::1234",
		"tcp",
		":1234",
		"",
		ErrorIncorrectDSN,
	},
	{
		"syslog:foo.txt",
		"",
		":1234",
		"",
		ErrorIncorrectDSN,
	},
}

func TestParseDsn(t *testing.T) {

	for _, tcase := range tests {
		actualProtocol, actualListen, actualApplication, actualErr := parseSyslogDsn(tcase.dsn)

		if actualErr != tcase.expectedError {
			t.Error("expected err ", tcase.expectedError,
				"actual err", actualErr, tcase.dsn)
			continue
		}

		if actualErr == tcase.expectedError && actualErr != nil {
			continue
		}

		if actualApplication != tcase.expectedApplication {
			t.Error("expectedApplication must ", tcase.expectedApplication,
				"actualApplication was ", actualApplication, tcase.dsn)
		}

		if actualProtocol != tcase.expectedProtocol {
			t.Error("expectedProtocol must ", tcase.expectedProtocol,
				"actualProtocol was ", actualProtocol, tcase.dsn, actualErr)
		}

		if actualListen != tcase.expectedListen {
			t.Error("expectedListen must ", tcase.expectedListen,
				"actualListen was ", actualProtocol, tcase.dsn)
		}

	}
}
