# Haproxyparse

`haproxyparse` is a small utility for parsing haproxy logs into a format compatible with [go-bench](https://github.com/clever/go-bench).

## Usage

```bash
Usage of ./haproxyparse:
  -e="": Only output lines starting before specified date (in RFC3339 format)
  -f="": Format to pull haproxy log attributes into go-bench "extras" param
  -m="": Only output lines with the specified http method (GET,POST,etc)
  -n=false: Use normalized timestamps (all requests evenly spaced one second apart
  -s="": Only output lines starting after specified date (in RFC3339 format)
```

You can either pass an haproxy log to `haproxyparse` via stdin or specify the file path with the `-f` flag.

## Building

`haproxyparse` should build with go 1.3+, with simply:
```bash
$ go build
```

## Testing

Run tests with:
```bash
$ go test
```
