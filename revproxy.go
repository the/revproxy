package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func usage() {
	fmt.Println("Usage: revproxy [OPTIONS] target")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	port := flag.Uint("port", 8080, "proxy port")
	host := flag.String("host", "", "override host header")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	target, err := url.Parse(flag.Arg(0))
	if err != nil || target.Host == "" {
		log.Panicf("invalid target URL: %s", target)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		serving := r.URL.Path
		if r.URL.RawQuery != "" {
			serving += "?" + r.URL.RawQuery
		}
		log.Printf("%s %s", r.Method, serving)

		director(r)
		if *host != "" {
			r.Host = *host
		}
	}

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Printf("proxying %s at %s", target, addr)
	http.Handle("/", proxy)
	http.ListenAndServe(addr, nil)
}
