package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"strings"
)

const (
	prefixEncoding = `=?utf-8?B?`
	suffixEncoding = `?=`
)

var (
	subject = `=?utf-8?B?5qyi6L+O5oKo5L2/55So5rOJ5bee5biC5Lit5LyB5L+h5oGv5oqA?=
 =?utf-8?B?5pyv5pyJ6ZmQ5YWs5Y+456aP5bee5YiG5YWs5Y+46YKu566x77yB?=`
	from = `=?utf-8?B?5rOJ5bee5biC5Lit5LyB5L+h5oGv5oqA5pyv5pyJ6ZmQ5YWs?= =?utf-8?B?5Y+456aP5bee5YiG5YWs5Y+46YKu566x57O757uf566h55CG5ZGY?= <admin@2009400.cn>`
)

func main() {
	fmt.Println(decodeHeader(from))
	fmt.Println(decodeHeader(subject))
	fmt.Println(decodeQuotedPrintable("mail.html"))
}

func base64Decode(s string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func decodeHeader(s string) string {
	var res string

	s = filterLF(s)
	ss := strings.Split(s, " ")
	for _, val := range ss {
		tmp := val
		if strings.Index(val, prefixEncoding) == 0 {
			tmp, _ = base64Decode(strings.TrimSuffix(strings.TrimPrefix(val, prefixEncoding), suffixEncoding))
		}
		res += tmp
	}

	return res
}

func filterLF(input string) string {
	var res []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		tmp := scanner.Text()
		res = append(res, strings.TrimSpace(tmp))
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return strings.Join(res, " ")
}

func decodeQuotedPrintable(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(quotedprintable.NewReader(bytes.NewReader(bs)))
	if err != nil {
		return ""
	}
	return string(b)
}
