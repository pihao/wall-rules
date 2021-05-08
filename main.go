package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	fmt.Println("hi")

	gfw := gfwlist()
	a, b := domainList(gfw)
	saveSurge(a, b)
}

func saveSurge(allow, block [][]byte) {
	buf := bytes.NewBuffer(nil)
	for _, b := range allow {
		buf.Write(b)
		buf.WriteString("\n")
	}
	writeFile("a.txt", buf.Bytes())

	buf = bytes.NewBuffer(nil)
	for _, b := range block {
		buf.Write(b)
		buf.WriteString("\n")
	}
	writeFile("b.txt", buf.Bytes())
}

var line = 0

func domainList(gfw []byte) (allow, block [][]byte) {
	scanner := bufio.NewScanner(bytes.NewBuffer(gfw))
	for scanner.Scan() {
		line++
		domain, whitelist := parseDomain(scanner.Bytes())
		if domain == nil {
			continue
		}

		if whitelist {
			allow = append(allow, domain)
		} else {
			block = append(block, domain)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return allow, block
}

var (
	domainReg = regexp.MustCompile(`^[\w\-]+\.[\w\-]+(\.[\w\-]+)*$`)
	ipReg     = regexp.MustCompile(`^\d+(\.\d+)+$`)
)

// https://github.com/gfwlist/gfwlist/wiki/Syntax
// !  注释
// @@ 白名单, 除了白名单, 其它全是黑名单
// || 匹配域名和子域名(任意scheme)
// |  匹配前缀(可指定scheme)
func parseDomain(b []byte) (domain []byte, whitelist bool) {
	origin := b

	// remove space
	b = bytes.TrimSpace(b)
	// remove NUL char
	b = bytes.Trim(b, "\x00")

	if len(b) == 0 {
		return nil, false
	}

	if bytes.HasPrefix(b, []byte("!")) || // 注释
		bytes.HasPrefix(b, []byte("[")) || // 特殊的第一行: [AutoProxy 0.2.9]
		bytes.HasPrefix(b, []byte("/")) { // 正则
		return nil, false
	}

	// 白名单
	if bytes.HasPrefix(b, []byte("@@")) {
		b = bytes.TrimPrefix(b, []byte("@@"))
		whitelist = true
	} else {
		whitelist = false
	}

	// remote prefix
	b = bytes.TrimPrefix(b, []byte("||"))
	b = bytes.TrimPrefix(b, []byte("|"))
	b = bytes.TrimPrefix(b, []byte("https://"))
	b = bytes.TrimPrefix(b, []byte("http://"))
	b = bytes.TrimPrefix(b, []byte("*."))
	b = bytes.TrimPrefix(b, []byte("."))

	// remove path suffix
	b = bytes.SplitN(b, []byte("/"), 2)[0]
	b = bytes.SplitN(b, []byte("%2F"), 2)[0]

	// IP address
	if ipReg.Match(b) {
		return nil, false
	}

	if domainReg.Match(b) {
		// copy bytes
		b = []byte(string(b))
		return b, whitelist
	}

	fmt.Printf("invalid domain: %d, [%s]\n", line, origin)
	return nil, false
}

func gfwlist() []byte {
	// res, err := http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	res, err := http.Get("http://127.0.0.1:8080/gfwlist.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// return body
	decoded := make([]byte, len(body))
	_, err = base64.StdEncoding.Decode(decoded, body)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

func writeFile(path string, b []byte) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
}
