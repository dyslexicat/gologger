package main

import "fmt"

// Request is...
type Request struct {
	method      string
	requestPath string
}

//ipLogs is...
type ipLogs struct {
	requests     map[Request]int
	browsers     map[string]int
	requestCount int
}

// stringer for iplogs struct
func (ipl *ipLogs) String() string {
	return fmt.Sprintf("%d request(s) have been made to the following endpoints: \n %v \n From the following browser(s): \n %v", ipl.requestCount, ipl.requests, ipl.browsers)
}

// check if a given browser exists in our logs for a given ip
func (ipl *ipLogs) checkBrowser(b string) {
	if _, ok := ipl.browsers[b]; ok {
		ipl.browsers[b]++
		return
	}

	ipl.browsers[b] = 1
}

// check if a request exists in our logs for a given ip
func (ipl *ipLogs) checkRequest(r Request) {
	if _, ok := ipl.requests[r]; ok {
		ipl.requests[r]++
		return
	}

	ipl.requests[r] = 1
}

// For given a given request line from STDIN create and return a Request
func handleRequestLine(s []string) Request {
	var r Request

	switch len(s) {
	case 0:
		r = Request{}
	case 1:
		r = Request{requestPath: s[0]}
	default:
		r = Request{method: s[0], requestPath: s[1]}
	}

	return r
}
