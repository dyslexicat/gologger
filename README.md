# gologger

### Disclaimer: current tool was written to only work with NGINX access logs. If you need to make it work with another log format you would need to edit the handleScanLine function.

## purpose
As soon as I created a website on Digitalocean, I was bombarded with automated attacks which was annoying to say the least.

I built this tool to both practice writing some Go code and at least ban the IPs where these attacks are coming from which is not the ideal solution but it is something.

I am planning on introducing concurrency to the code (again for practice) but even as it is the code is quite fast.
It reads from STDIN and outputs the relevant information to STDOUT which makes it quite functional.

## usage
For example, you can use the following command to read through all your log files

> **for FILE in access.log.\*; do cat $FILE; done | ./gologger -b=request-path -w=asp,admin,conf,cgi > ips_to_be_banned**

which would output all the IPs that sent one or more requests to endpoints that includes words from the provided wordlist. In the example above, any request that contain "asp", "admin", "conf" or "cgi" in the endpoint would match.

**If you have any other requests please feel free to ask**
 