package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO more tests

type valOrEmptyTest struct {
	from     map[string]string
	val      string
	expected string
}

var valOrEmptyTests = []valOrEmptyTest{
	{
		from:     map[string]string{"a": "b"},
		val:      "a",
		expected: "b",
	},
	{
		from:     map[string]string{"a": "b"},
		val:      "c",
		expected: "",
	},
}

func TestValOrEmpty(t *testing.T) {
	for _, testCase := range valOrEmptyTests {
		actual := valOrEmpty(testCase.from, testCase.val)
		assert.Equal(t, testCase.expected, actual)
	}
}

type getMatchedSubexprTest struct {
	regex    *regexp.Regexp
	toMatch  string
	expected map[string]string
}

var getMatchedSubexprTests = []getMatchedSubexprTest{
	{
		regex:    regexp.MustCompile(`^(?P<first>\w+)\s(?P<second>\w+)$`),
		toMatch:  "first second",
		expected: map[string]string{"first": "first", "second": "second"},
	},
}

func TestGetMatchedSubexpr(t *testing.T) {
	for _, testCase := range getMatchedSubexprTests {
		actual := getMatchedSubexpr(testCase.regex, testCase.toMatch)
		assert.Equal(t, testCase.expected, actual)
	}
}

type parseLineTest struct {
	line     string
	expected HAProxyLogLine
}

func parseTimeIgnoreError(s string) time.Time {
	t, _ := time.Parse(DateFormat, s)
	return t
}

var parseLineTests = []parseLineTest{
	//	{
	//		line: `Jan  1 00:00:00 local0 haproxy[2072]: 1.2.3.4:5555 [02/Dec/2014:16:41:06.226] front back/server 0/0/0/7/7 200 359 - - ---- 72/72/1/1/0 0/0 {|ELB-HealthChecker/1.0|Basic something} {|} "GET /elb/check HTTP/1.1"`,
	//		expected: HAProxyLogLine{
	//			ProcessId:           "2072",
	//			ClientIp:            "1.2.3.4",
	//			ClientPort:          "5555",
	//			Date:                "02/Dec/2014:16:41:06.226",
	//			FrontendName:        "front",
	//			BackendName:         "back",
	//			ServerName:          "server",
	//			SendTime:            "0",
	//			WaitTime:            "0",
	//			ConnectionTime:      "0",
	//			ResponseTime:        "7",
	//			TotalTime:           "7",
	//			StatusCode:          "200",
	//			BytesRead:           "359",
	//			RequestCookie:       "-",
	//			ResponseCookie:      "-",
	//			TerminationState:    "----",
	//			ActiveConnections:   "72",
	//			FrontendConnections: "72",
	//			BackendConnections:  "1",
	//			ServerConnections:   "1",
	//			Retries:             "0",
	//			ServerQueue:         "0",
	//			BackendQueue:        "0",
	//			AuthHeader:          "Basic something",
	//			ResponseHeaders:     "|",
	//			HttpMethod:          "GET",
	//			HttpUri:             "/elb/check",
	//			UriComplete:         true,
	//			Time:                parseTimeIgnoreError("02/Dec/2014:16:41:06.226"),
	//		},
	//	},
	//	{
	//		line: `Dec  2 16:41:06 local0 haproxy[2072]: 1.2.3.4:5555 [02/Dec/2014:16:41:06.226] front back/server 0/0/0/7/7 200 359 - - ---- 72/72/1/1/0 0/0 {|ELB-HealthChecker/1.0|} {|} "GET /elb/check`,
	//		expected: HAProxyLogLine{
	//			ProcessId:           "2072",
	//			ClientIp:            "1.2.3.4",
	//			ClientPort:          "5555",
	//			Date:                "02/Dec/2014:16:41:06.226",
	//			FrontendName:        "front",
	//			BackendName:         "back",
	//			ServerName:          "server",
	//			SendTime:            "0",
	//			WaitTime:            "0",
	//			ConnectionTime:      "0",
	//			ResponseTime:        "7",
	//			TotalTime:           "7",
	//			StatusCode:          "200",
	//			BytesRead:           "359",
	//			RequestCookie:       "-",
	//			ResponseCookie:      "-",
	//			TerminationState:    "----",
	//			ActiveConnections:   "72",
	//			FrontendConnections: "72",
	//			BackendConnections:  "1",
	//			ServerConnections:   "1",
	//			Retries:             "0",
	//			ServerQueue:         "0",
	//			BackendQueue:        "0",
	//			AuthHeader:          "",
	//			ResponseHeaders:     "|",
	//			HttpMethod:          "GET",
	//			HttpUri:             "/elb/check",
	//			UriComplete:         false,
	//			Time:                parseTimeIgnoreError("02/Dec/2014:16:41:06.226"),
	//		},
	//	},
	{
		line: `Jan 1 00:00:00 local0 haproxy[2344]: 1.2.3.4 - - [02/Dec/2014:16:41:32 +0000] "GET /elb/check HTTP/1.1" 200 359 "" "" 5555 073 "front" "back" "server" 4995 0 0 5 5000 ---- 72 72 1 1 0 0 0 "" "" "" "ELB-HealthChecker/1.0" "Basic something"`,
		expected: HAProxyLogLine{
			ProcessId:           "2344",
			ClientIp:            "1.2.3.4",
			ClientPort:          "5555",
			Date:                "02/Dec/2014:16:41:32 +0000",
			FrontendName:        "front",
			BackendName:         "back",
			ServerName:          "server",
			SendTime:            "4995",
			WaitTime:            "0",
			ConnectionTime:      "0",
			ResponseTime:        "5",
			TotalTime:           "5000",
			StatusCode:          "200",
			BytesRead:           "359",
			RequestCookie:       "",
			ResponseCookie:      "",
			TerminationState:    "----",
			ActiveConnections:   "72",
			FrontendConnections: "72",
			BackendConnections:  "1",
			ServerConnections:   "1",
			Retries:             "0",
			ServerQueue:         "0",
			BackendQueue:        "0",
			XForward:            "",
			UserAgent:           "ELB-HealthChecker/1.0",
			AuthHeader:          "Basic something",
			HttpMethod:          "GET",
			HttpUri:             "/elb/check",
			UriComplete:         true,
			Time:                parseTimeIgnoreError("02/Dec/2014:16:41:32 +0000"),
		},
	},
}

func TestParseLine(t *testing.T) {
	for _, testCase := range parseLineTests {
		actual, err := parseLine(testCase.line)
		assert.Nil(t, err)
		assert.Equal(t, testCase.expected, actual)
	}
}
