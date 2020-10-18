package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getIPaddress(s string) string {
	var result string
	for index, char := range s {
		if unicode.IsSpace(char) {
			result = s[:index]
			break
		}
	}

	return result
}

func printUniqueIPs(logs map[string]*ipLogs) {
	for key := range logs {
		fmt.Println(key)
	}
}

func main() {
	f, err := os.Open("./logs/access.log")
	check(err)

	defer f.Close()

	unique := flag.Bool("unique", false, "print unique IPs in log")
	flag.Parse()

	logs := make(map[string]*ipLogs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "\"")
		ip := getIPaddress(split[0])
		browser := split[5]
		requestLine := strings.Fields(split[1])

		r := handleRequestLine(requestLine)

		if _, ok := logs[ip]; !ok {
			var p = new(ipLogs)
			logs[ip] = p
			p.browsers = make(map[string]int)
			p.requests = make(map[Request]int)

			p.checkBrowser(browser)
			p.checkRequest(r)
		} else {
			logs[ip].checkBrowser(browser)
			logs[ip].checkRequest(r)
		}

		(*logs[ip]).requestCount++
	}

	if *unique {
		printUniqueIPs(logs)
	}

	for key, value := range logs {
		fmt.Println(key, value.requestCount)
	}
}
