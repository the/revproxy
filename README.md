revproxy
========

Simple reverse proxy that allows overriding the `Host` header.

Installation
------------

```
go get github.com/the/revproxy
```

Usage
-----

```
revproxy [OPTIONS] target
  -H	print request headers
  -color
    	format output with colors
  -host string
        override host header
  -port uint
        proxy port (default 8080)
```

Example

```
revproxy -port 9000 -host www.google.ch https://www.google.ch
```
