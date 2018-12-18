package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const (
	formatReset  = "\033[0m"
	formatBold   = "\033[1m"
	formatCyan   = "\033[36m"
	formatGray   = "\033[90m"
	formatBlack  = "\033[30m"
	formatYellow = "\033[33m"
)

var renderColors *bool

func usage() {
	fmt.Println("Usage: revproxy [OPTIONS] target")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	port := flag.Uint("port", 8080, "proxy port")
	host := flag.String("host", "", "override host header")
	printHeaders := flag.Bool("H", false, "print request headers")
	renderColors = flag.Bool("color", false, "format output with colors")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	target, err := url.Parse(flag.Arg(0))
	if err != nil || target.Host == "" {
		log.Panicf("invalid target URL: %s", target)
	}

	if *host == "" {
		*host = target.Host
	}

	questionMark := format("?", formatBlack)
	equals := format("=", formatGray)
	ampersand := format("&", formatBlack)

	proxy := httputil.NewSingleHostReverseProxy(target)
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		serving := format(r.URL.Path, formatYellow)
		if r.URL.RawQuery != "" {
			serving += questionMark
			var kvs []string
			for name, values := range r.URL.Query() {
				for _, value := range values {
					kvs = append(kvs, format(name, formatCyan)+equals+value)
				}
			}
			serving += strings.Join(kvs, ampersand)
		}
		log.Printf("%s %s", format(r.Method, formatBold), serving)

		if *printHeaders {
			for name, values := range r.Header {
				value := strings.Join(values, "; ")
				log.Printf("%s: %s", format(name, formatGray), value)
			}
		}

		director(r)
		r.Host = *host
	}

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Printf("proxying %s at %s", target, addr)
	http.Handle("/", proxy)
	http.ListenAndServe(addr, nil)
}

func format(value, format string) string {
	if *renderColors {
		return format + value + formatReset
	}

	return value
}
