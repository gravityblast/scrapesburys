package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		url           string
		concurrency   int
		logLevelError bool
		logLevelInfo  bool
		logLevelDebug bool
	)

	flag.IntVar(&concurrency, "c", 1, "concurreny, min value is 1")
	flag.StringVar(&url, "u", "", "the main URL")
	flag.BoolVar(&logLevelError, "v", false, "log level error")
	flag.BoolVar(&logLevelInfo, "vv", false, "log level info")
	flag.BoolVar(&logLevelDebug, "vvv", false, "log level debug")
	flag.Parse()

	if url == "" {
		fmt.Printf("-u flag is mandatory\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if concurrency < 1 {
		fmt.Printf("invalid concurrency value\n\n")
		flag.Usage()
		os.Exit(1)
	}

	s, err := NewProductsScraper(url)
	if err != nil {
		log.Fatal(err)
	}

	var logLevel int

	switch {
	case logLevelDebug:
		logLevel = LogLevelDebug
	case logLevelInfo:
		logLevel = LogLevelInfo
	case logLevelError:
		logLevel = LogLevelError
	}

	appLog := newAppLogger(logLevel)

	res, err := s.Scrape(appLog, concurrency)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(res)
}
