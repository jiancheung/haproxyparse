package main

import (
	"log"
	"regexp"
	"strings"
	"time"
)

// HAProxyLine represents a single line of the HAProxy log, and all fields it contains.
// All fields are captured as strings.
type HAProxyLogLine struct {
	ProcessId           string
	ClientIp            string
	ClientPort          string
	Date                string
	FrontendName        string
	BackendName         string
	ServerName          string
	SendTime            string
	WaitTime            string
	ConnectionTime      string
	ResponseTime        string
	TotalTime           string
	StatusCode          string
	BytesRead           string
	RequestCookie       string
	ResponseCookie      string
	TerminationState    string
	ActiveConnections   string
	FrontendConnections string
	BackendConnections  string
	ServerConnections   string
	Retries             string
	ServerQueue         string
	BackendQueue        string
	XForward            string
	UserAgent           string
	AuthHeader          string
	HttpMethod          string
	HttpUri             string
	UriComplete         bool
	Time                time.Time
}

//http://blog.haproxy.com/2012/10/29/haproxy-log-customization/
var LogLineRegex = regexp.MustCompile(
	`^.*haproxy\[(?P<pid>\d+)\]:\s` +
		`(?P<clientip>\d+\.\d+\.\d+\.\d+)\s` + //%Ci
		`-\s-\s` +
		`\[(?P<date>\S+\s\S+)\]\s` + // [%t]
		`\"(?P<httpmethod>\w+)\s(?P<httpuri>\S+)(\s(?P<httpprotocol>\S+)\")?\s` + //%r
		`(?P<statuscode>[\d-]+)\s` + //%st
		`(?P<bytesread>\d+)\s` + //%B
		`""\s""\s` +
		`(?P<clientport>\d+)\s` + //%Cp
		`(?P<milliseconds>\d+)\s` + //%ms
		`"(?P<frontend>\S*?)"\s` + //%ft
		`"(?P<backend>\S*?)"\s` + //%b
		`"(?P<server>\S*?)"\s` + //%s
		`(?P<timesend>[\d-]+)\s(?P<timewait>[\d-]+)\s(?P<timeconnection>[\d-]+)\s(?P<timeresponse>[\d-]+)\s(?P<timetotal>[\d-]+)\s` + //%Tq %Tw %Tc %Tr %Tt
		`(?P<terminationstate>\S+)\s` + //%tsc
		`(?P<activeconnections>\d+)\s(?P<frontendconnections>\d+)\s(?P<backendconnections>\d+)\s(?P<serverconnections>\d+)\s(?P<retries>\d+)\s` + //%ac %fc %bc %sc %rc
		`(?P<serverqueue>\d+)\s(?P<backendqueue>\d+)\s` + //%sq %bq
		`"(?P<requestcookie>\S*?)"\s"(?P<responsecookie>\S*?)"\s` + //%cc %cs
		`"(?P<xforward>.*?)"\s"(?P<useragent>.*?)"\s"(?P<authheader>.*?)"` + // custom
		`$`)

const DateFormat = "02/Jan/2006:15:04:05 -0700"

func parseLine(logLine string) (HAProxyLogLine, error) {
	matches := getMatchedSubexpr(LogLineRegex, logLine)
	parsed := HAProxyLogLine{}
	// gotta be an easier way...
	parsed.ProcessId = valOrEmpty(matches, "pid")
	parsed.ClientIp = valOrEmpty(matches, "clientip")
	parsed.ClientPort = valOrEmpty(matches, "clientport")
	parsed.Date = valOrEmpty(matches, "date")
	parsed.FrontendName = valOrEmpty(matches, "frontend")
	parsed.BackendName = valOrEmpty(matches, "backend")
	parsed.ServerName = valOrEmpty(matches, "server")
	parsed.SendTime = valOrEmpty(matches, "timesend")
	parsed.WaitTime = valOrEmpty(matches, "timewait")
	parsed.ConnectionTime = valOrEmpty(matches, "timeconnection")
	parsed.ResponseTime = valOrEmpty(matches, "timeresponse")
	parsed.TotalTime = valOrEmpty(matches, "timetotal")
	parsed.StatusCode = valOrEmpty(matches, "statuscode")
	parsed.BytesRead = valOrEmpty(matches, "bytesread")
	parsed.RequestCookie = valOrEmpty(matches, "requestcookie")
	parsed.ResponseCookie = valOrEmpty(matches, "responsecookie")
	parsed.TerminationState = valOrEmpty(matches, "terminationstate")
	parsed.ActiveConnections = valOrEmpty(matches, "activeconnections")
	parsed.FrontendConnections = valOrEmpty(matches, "frontendconnections")
	parsed.BackendConnections = valOrEmpty(matches, "backendconnections")
	parsed.ServerConnections = valOrEmpty(matches, "serverconnections")
	parsed.Retries = valOrEmpty(matches, "retries")
	parsed.ServerQueue = valOrEmpty(matches, "serverqueue")
	parsed.BackendQueue = valOrEmpty(matches, "backendqueue")
	parsed.XForward = valOrEmpty(matches, "xforward")
	parsed.UserAgent = valOrEmpty(matches, "useragent")
	parsed.AuthHeader = valOrEmpty(matches, "authheader")
	parsed.HttpMethod = valOrEmpty(matches, "httpmethod")
	parsed.HttpUri = valOrEmpty(matches, "httpuri")
	// if we got the http protocol, then the uri was complete,
	// otherwise it was truncated
	if valOrEmpty(matches, "httpprotocol") != "" {
		parsed.UriComplete = true
	}

	// parse date format
	parsetime, err := time.Parse(DateFormat, parsed.Date)
	if err != nil {
		return HAProxyLogLine{}, err
		log.Printf("unable to parse time: %s, original line: %s, parsed: %v, matches: %v\n", err, logLine, parsed, matches)
	}
	parsed.Time = parsetime

	// escape commas
	if parsed.UriComplete {
		parsed.HttpUri = strings.Replace(parsed.HttpUri, ",", "%2c", -1)
	}

	return parsed, nil
}

func valOrEmpty(from map[string]string, val string) string {
	v, ok := from[val]
	if !ok {
		return ""
	}
	return v
}

func getMatchedSubexpr(re *regexp.Regexp, toMatch string) map[string]string {
	results := map[string]string{}
	matches := re.FindStringSubmatch(toMatch)
	if matches == nil {
		log.Println("NIL OH NOES")
		return results
	}
	names := re.SubexpNames()
	for ix, match := range matches {
		if ix == 0 || names[ix] == "" {
			continue
		}
		results[names[ix]] = match
	}
	return results
}
