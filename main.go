package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

// Logs is..
type Logs map[string]*ipLogs

// error checking
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
func printUniqueIPs(logs map[string]*ipLogs) {
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

// for a given line from our scanner separate ip, browser, and request then handle the case according to the banFlag
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
	f, err := os.Open("./logs/access.log")
	check(err)

	defer f.Close()

	uniqueFlag := flag.Bool("unique", false, "print unique IPs in log")
	banFlag := flag.Bool("ban", false, "ban ips for a given ip list")
	verbose := flag.Bool("verbose", false, "output more information about things")

	var banMode string
	flag.StringVar(&banMode, "banmode", "", "ban mode according to a rule set")

	var wordlist string
	flag.StringVar(&wordlist, "wordlist", "", "comma separated list of endpoints")

	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "set the concurrency level")

	flag.Parse()

	if banMode != "" {
		rules = strings.Split(wordlist, ",")
	}

	scanner := bufio.NewScanner(f)

	if *banFlag {
		for scanner.Scan() {
			//line := scanner.Text()
			fmt.Print("banning this shit")
		}
	} else {
		for scanner.Scan() {
			line := scanner.Text()
			handleScanLine(line, banMode)
		}

	}

	if *uniqueFlag {
		printUniqueIPs(logs)
		return
	}

	if *verbose && banMode == "" {
		for key, value := range logs {
			fmt.Println(key, value)
		}
		return
	}
}
