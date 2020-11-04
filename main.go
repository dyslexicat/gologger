package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Logs is..
type Logs map[string]*ipLogs

// get the IP address from a given string from our logs
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

// if the unique flag is set, print all the unique ips that sent a request to our server
func printUniqueIPs(logs Logs) {
	for key := range logs {
		fmt.Println(key)
	}
}

// if the current line contains a rule inside we ban
func handleBan(line string, ruleset []string, ip string) {
	for _, rule := range ruleset {
		if strings.Contains(line, rule) {
			fmt.Println(ip)
		}

	}
}

// for a given ip check whether that ip is in our logs
func checkIPLogs(ip string) bool {
	if _, ok := logs[ip]; ok {
		return true
	}

	return false
}

// updating our logs according to a given ip and then checking to see if we have seen given browser and request before
func handleLogging(ip string, browser string, request Request) {
	if checkIPLogs(ip) {
		logs[ip].checkBrowser(browser)
		logs[ip].checkRequest(request)
	} else {
		var p = new(ipLogs)
		logs[ip] = p
		p.browsers = make(map[string]int)
		p.requests = make(map[Request]int)

		p.checkBrowser(browser)
		p.checkRequest(request)
	}

	logs[ip].requestCount++
}

var logs = make(Logs)
var rules []string

// for a given line from our scanner separate ip, browser, and request then handle the case according to the banMode
func handleScanLine(line string, banMode string) {
	split := strings.Split(line, "\"")
	ip := getIPaddress(split[0])
	browser := split[5]
	requestLine := strings.Fields(split[1])
	request := handleRequestLine(requestLine)

	if banMode != "" && checkIPLogs(ip) {
		return
	}

	switch banMode {
	case "browser":
		handleBan(browser, rules, ip)
	case "request-path":
		handleBan(request.requestPath, rules, ip)
	case "request-method":
		handleBan(request.method, rules, ip)
	}

	handleLogging(ip, browser, request)
}

func main() {
	uniqueFlag := flag.Bool("u", false, "Prints unique IPs in logs. Usage: ./main --u")
	verbose := flag.Bool("v", false, "Output more information about things. Usage: ./main --v")

	var banMode string
	flag.StringVar(&banMode, "b", "", "Use together with the --w flag. Available modes: browser / request-path / request-method")

	var wordlist string
	flag.StringVar(&wordlist, "w", "", "Comma separated list of rules. Examples: --w=python,perl")

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "Set the concurrency level. Currently not implemented.")

	flag.Parse()

	if banMode != "" {
		rules = strings.Split(wordlist, ",")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		handleScanLine(line, banMode)
	}

	if *uniqueFlag {
		printUniqueIPs(logs)
		return
	}

	if *verbose && banMode == "" {
		for key, value := range logs {
			fmt.Println(key)
			for req := range value.requests {
				fmt.Println("\t", req)
			}
		}
		return
	}
}
