package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

var (
	method       string
	start        string
	end          string
	normalized   bool
	extrasFormat string
)

func init() {
	flag.StringVar(&method, "m", "", "Only output lines with the specified http method (GET,POST,etc)")
	flag.StringVar(&start, "s", "", "Only output lines starting after specified date (in RFC3339 format)")
	flag.StringVar(&end, "e", "", "Only output lines starting before specified date (in RFC3339 format)")
	flag.BoolVar(&normalized, "n", false, "Use normalized timestamps (all requests evenly spaced one second apart")
	flag.StringVar(&extrasFormat, "f", "", "Format to pull haproxy log attributes into go-bench \"extras\" param")
}

func main() {
	flag.Parse()

	var extrasTemplate *template.Template
	var err error
	if extrasFormat != "" {
		extrasTemplate, err = template.New("extras").Parse(extrasFormat)
		if err != nil {
			log.Fatalf("error parsing extras template: %s", err)
		}
	}

	var startTime time.Time
	var endTime time.Time

	if start != "" {
		startTime = parseTime(start)
	}
	if end != "" {
		endTime = parseTime(end)
	}

	var scanner *bufio.Scanner
	if flag.NArg() > 0 {
		fd, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatalf("error opening input file: %s", err)
		}
		scanner = bufio.NewScanner(fd)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	totalOffset := int64(0)
	first := true
	var lastTime time.Time
	var templateBuffer bytes.Buffer
	for scanner.Scan() {
		parsedline, err := parseLine(scanner.Text())
		if err != nil {
			continue
		}
		if start != "" && parsedline.Time.Before(startTime) {
			continue
		}
		if end != "" && parsedline.Time.After(endTime) {
			os.Exit(0)
		}
		if !parsedline.UriComplete {
			continue
		}
		if method != "" && !equalsIgnoreCase(method, parsedline.HttpMethod) {
			continue
		}

		if first {
			lastTime = parsedline.Time
			first = false
		}
		if normalized {
			totalOffset += int64(1000)
		} else {
			totalOffset += int64(parsedline.Time.Sub(lastTime).Seconds() * 1000.0)
		}

		extras := ""
		if extrasTemplate != nil {
			if err = extrasTemplate.Execute(&templateBuffer, parsedline); err != nil {
				log.Fatalf("error applying template: %s", err)
			}
			extras = templateBuffer.String()
			templateBuffer.Reset()
		}

		fmt.Fprintf(os.Stdout, "%d,%s,%s,%s,%s\n", totalOffset, parsedline.HttpMethod, parsedline.HttpUri, parsedline.AuthHeader, extras)

		lastTime = parsedline.Time
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading input: %s", err)
	}
	os.Exit(0)
}

func equalsIgnoreCase(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func parseTime(t string) time.Time {
	out, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Fatalf("error parsing start time: %s", err)
	}
	return out
}
