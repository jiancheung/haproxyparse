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
	AuthHeader          string
	ResponseHeaders     string
	HttpMethod          string
	HttpUri             string
	UriComplete         bool
	Time                time.Time
}

var LogLineRegex = regexp.MustCompile(
	`^.*haproxy\[(?P<pid>\d+)\]:\s` +
		`(?P<clientip>\d+\.\d+\.\d+\.\d+):` +
		`(?P<clientport>\d+)\s` +
		`\[(?P<date>\S+)\]\s` +
		`(?P<frontend>\S+)\s` +
		`(?P<backend>\S+)\/(?P<server>\S+)\s` +
		`(?P<timesend>[\d-]+)\/(?P<timewait>[\d-]+)\/(?P<timeconnection>[\d-]+)\/(?P<timeresponse>[\d-]+)\/(?P<timetotal>[\d-]+)\s` +
		`(?P<statuscode>[\d-]+)\s` +
		`(?P<bytesread>\d+)\s` +
		`(?P<requestcookie>\S+)\s(?P<responsecookie>\S+)\s` +
		`(?P<terminationstate>\S+)\s` +
		`(?P<activeconnections>\d+)\/(?P<frontendconnections>\d+)\/(?P<backendconnections>\d+)\/(?P<serverconnections>\d+)\/(?P<retries>\d+)\s` +
		`(?P<serverqueue>\d+)\/(?P<backendqueue>\d+)\s` +
		`(\{[^|]*\|[^|]*\|(?P<authheader>.*?)\}\s\{(?P<responseheaders>.*)\}\s){0,1}` +
		`\"(?P<httpmethod>\w+)\s(?P<httpuri>\S+)(\s(?P<httpprotocol>\S+)\")?$`)

const DateFormat = "02/Jan/2006:15:04:05.000"

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
	parsed.AuthHeader = valOrEmpty(matches, "authheader")
	parsed.ResponseHeaders = valOrEmpty(matches, "responseheaders")
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
