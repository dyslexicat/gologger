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

// error checking
func check(err error) {
	if err != nil {
		panic(err)
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

// for a given ip and logs map check whether a key is in the map and check relevant browser and request maps for
// given browser and request
func checkIPLogs(ip string, logs Logs, browser string, r Request) {
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

	logs[ip].requestCount++
}

func main() {
	f, err := os.Open("./logs/access.log")
	check(err)

	defer f.Close()

	uniqueFlag := flag.Bool("unique", false, "print unique IPs in log")
	//banFlag := flag.Bool("ban mode", false, "ban mode according to a rule set")
	//browserFlag := flag.String("browsers", "", "comma separated list of browsers")

	flag.Parse()

	logs := make(Logs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "\"")
		ip := getIPaddress(split[0])
		browser := split[5]
		requestLine := strings.Fields(split[1])
		r := handleRequestLine(requestLine)

		checkIPLogs(ip, logs, browser, r)
	}

	if *uniqueFlag {
		printUniqueIPs(logs)
	}

	for key, value := range logs {
		fmt.Println(key, value)
	}
}
