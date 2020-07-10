package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"

	"zq/mail/parsemail"

	"github.com/bytbox/go-pop3"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	address = "pop.qiye.163.com:110"
	user    = "php01@2009400.cn"
	pass    = "ZQxx400@"
)

var (
	prefixEncoding = `=?utf-8?B?`
	suffixEncoding = `?=`
)

func setEncoding(charset string) {
	prefixEncoding = charset
	fmt.Println("current charset:", charset)
}

func main() {

	rb := getEmailMsgFromRemote()
	useParsEmail(rb)

	// filename := `d:\goProjects\src\mail_gb2312.txt`
	// rb := readFile(filename)
	// useParsEmail(rb)
}

func getEmailMsgFromRemote() *strings.Reader {
	s := cmdSendMail(9)
	return strings.NewReader(s)
}

// d:\goProjects\src\mail_utf8.txt
func readFile(filename string) *bytes.Reader {
	dataBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("read mail file error.", err)
	}

	return bytes.NewReader(dataBytes)
}

func useParsEmail(reader io.Reader) {
	email, err := parsemail.Parse(reader) // returns Email struct and error
	if err != nil {
		log.Fatal("parse email error.", err)
	}

	subject := email.Header.Get("Subject")
	if strings.Index(subject, "=?") == 0 {
		charset := subject[:strings.Index(subject, "B?")+2]
		setEncoding(charset)
	}
	fmt.Println("Subject:", DecodeHeader(email.Header.Get("Subject")))
	fmt.Println("From:", DecodeHeader(email.Header.Get("From")))
	fmt.Println("To:", DecodeHeader(email.Header.Get("To")))
	fmt.Println("Cc:", DecodeHeader(email.Header.Get("Cc")))
	fmt.Println(email.HTMLBody)
	// fmt.Println(email.TextBody)
	if len(email.TextBody) > 0 {
		textBody, _ := base64Decode(email.TextBody)
		tb, _ := GB23122UTF8([]byte(textBody))
		fmt.Println(string(tb))
	}

	if len(email.Attachments) > 0 {
		for _, item := range email.Attachments {
			filename := DecodeHeader(item.Filename)
			fmt.Println(filename)
			fmt.Println(item.ContentType)
			msBytes := ReadBytes(item.Data)
			ioutil.WriteFile(filename, msBytes, 0666)
		}

	}

}

func myselfParseEmail(reader io.Reader) {
	m, err := mail.ReadMessage(reader)

	sendTime, err := m.Header.Date()
	fmt.Println("Date:", sendTime.Format("2006-01-02 15:04:05"))
	fmt.Println("From:", DecodeHeader(m.Header.Get("From")))
	fmt.Println("To:", DecodeHeader(m.Header.Get("To")))
	subjectBytes, _ := GB23122UTF8([]byte(DecodeHeader(m.Header.Get("Subject"))))
	fmt.Printf("Subject: %s", string(subjectBytes))

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		log.Fatal(err)
	}
	DecodeQuotedPrintable(body)
	// fmt.Printf("%s\n", qp)

}

func DecodeHeader(s string) string {
	var res string

	// s = filterLF(s)
	ss := strings.Split(s, " ")
	for _, val := range ss {
		tmp := val
		if strings.Index(val, prefixEncoding) == 0 {
			tmp, _ = base64Decode(strings.TrimSuffix(strings.TrimPrefix(val, prefixEncoding), suffixEncoding))
		}
		res += tmp
	}

	var resBytes []byte
	switch prefixEncoding {
	case "=?GB2312?B?":
	case "=?gb2312?B?":
		resBytes, _ = GB23122UTF8([]byte(res))
	}

	if len(resBytes) > 0 {
		res = string(resBytes)
	}

	return res
}

func DecodeQuotedPrintable(bs []byte) string {

	b, err := ioutil.ReadAll(quotedprintable.NewReader(bytes.NewReader(bs)))
	if err != nil {
		return ""
	}
	return string(b)
}

func base64Decode(s string) (string, error) {
	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func cmdSendMail(msg int) string {
	client, err := pop3.Dial(address)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	defer func() {
		client.Quit()
	}()

	if err = client.User(user); err != nil {
		log.Printf("Error: %v\n", err)
		return ""
	}

	if err = client.Pass(pass); err != nil {
		log.Printf("Error: %v\n", err)
		return ""
	}

	var content string

	if content, err = client.Retr(msg); err != nil {
		log.Printf("Error: %v\n", err)
		return ""
	}

	return content
}

func GB23122UTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func checkErr(err error) {
	log.Panic(err)
}

func ReadBytes(reader io.Reader) []byte {
	var arr [1024]byte
	var buf []byte
	var i int = 0
	for {
		i++
		n, err := reader.Read(arr[:])
		if err == io.EOF {
			fmt.Println("file read finished")
			break
		}
		if err != nil {
			fmt.Println("file read failed")
			os.Exit(-1)
		}

		buf = append(buf, arr[:n]...)

	}
	fmt.Println("read 1024 times is:", i)

	return buf
}
