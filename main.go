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
	"text/template"
)

const (
	gfwlistURL = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

	surgeAllowFile = "conf/DOMAIN-SET/gfwlist-allow.txt"
	surgeBlockFile = "conf/DOMAIN-SET/gfwlist-block.txt"

	pacTemplate = "tpl/pac.js"
	pacFile     = "conf/gfwlist.pac"
)

func main() {
	gfw := gfwlist()
	a, b := domainList(gfw)
	saveSurge(a, b)
	savePAC(a, b)
}

func savePAC(allow, block [][]byte) {
	bufa := bytes.NewBuffer([]byte("\n"))
	for _, b := range allow {
		bufa.WriteString("    \"")
		bufa.Write(b)
		bufa.WriteString("\",\n")
	}

	bufb := bytes.NewBuffer([]byte("\n"))
	for _, b := range block {
		bufb.WriteString("    \"")
		bufb.Write(b)
		bufb.WriteString("\",\n")
	}

	rule := struct {
		Allow, Block string
	}{bufa.String(), bufb.String()}

	t, err := template.ParseFiles(pacTemplate)
	checkErr(err)

	var buf bytes.Buffer
	err = t.Execute(&buf, rule)
	checkErr(err)
	writeFile(pacFile, buf.Bytes())
}

func saveSurge(allow, block [][]byte) {
	buf := bytes.NewBuffer(nil)
	for _, b := range allow {
		buf.WriteString(".")
		buf.Write(b)
		buf.WriteString("\n")
	}
	writeFile(surgeAllowFile, buf.Bytes())

	buf = bytes.NewBuffer(nil)
	for _, b := range block {
		buf.WriteString(".")
		buf.Write(b)
		buf.WriteString("\n")
	}
	writeFile(surgeBlockFile, buf.Bytes())
}

var currentLine = 0

func domainList(gfw []byte) (allow, block [][]byte) {
	scanner := bufio.NewScanner(bytes.NewBuffer(gfw))
	for scanner.Scan() {
		currentLine++
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
	checkErr(scanner.Err())
	return allow, block
}

var (
	domainReg = regexp.MustCompile(`^[\w\-]+\.[\w\-]+(\.[\w\-]+)*$`)
	ipReg     = regexp.MustCompile(`^\d+(\.\d+)+$`)
	starReg   = regexp.MustCompile(`^\w*\*\w*\.`)
)

// https://github.com/gfwlist/gfwlist/wiki/Syntax
// !  注释
// @@ 白名单, 除了白名单, 其它全是黑名单
// || 匹配域名和子域名(任意scheme)
// |  匹配前缀(可指定scheme)
func parseDomain(b []byte) (domain []byte, allow bool) {
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
		allow = true
	} else {
		allow = false
	}

	// remove prefix
	b = bytes.TrimPrefix(b, []byte("||"))
	b = bytes.TrimPrefix(b, []byte("|"))
	b = bytes.TrimPrefix(b, []byte("https://"))
	b = bytes.TrimPrefix(b, []byte("http://"))
	b = bytes.TrimPrefix(b, []byte("*."))
	b = bytes.TrimPrefix(b, []byte("."))
	b = starReg.ReplaceAll(b, []byte{})

	// remove suffix
	b = bytes.SplitN(b, []byte("/"), 2)[0]
	b = bytes.SplitN(b, []byte("%2F"), 2)[0]

	// IP address
	if ipReg.Match(b) {
		return nil, false
	}

	if domainReg.Match(b) {
		b = []byte(string(b)) // copy bytes
		return b, allow
	}

	fmt.Printf("miss: %d, %s\n", currentLine, origin)
	return nil, false
}

func gfwlist() []byte {
	res, err := http.Get(gfwlistURL)
	checkErr(err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	decoded := make([]byte, len(body))
	_, err = base64.StdEncoding.Decode(decoded, body)
	checkErr(err)
	return decoded
}

func writeFile(path string, b []byte) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	checkErr(err)
	defer f.Close()

	_, err = f.Write(b)
	checkErr(err)
	f.Sync()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}