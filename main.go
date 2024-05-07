package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Server(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type loadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *loadBalancer {
	return &loadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)

	handleError(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func (lb *loadBalancer) getNextAvailableServer() Server {}

func (lb *loadBalancer) serverProxy(rw http.ResponseWriter, r *http.Request) {}

func main() {
	fmt.Println("Hello world from load balancer")
	servers := []Server{
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https:///www.facebook.com"),
		newSimpleServer("https://www.bing.com"),
	}

	lb := NewLoadBalancer("8000", servers)

	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serverProxy(rw, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Server is running at the port: %s", lb.port)

	http.ListenAndServe(":"+lb.port, nil)
}
