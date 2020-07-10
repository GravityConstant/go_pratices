package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/lisijie/gopub/app/libs"
	"gopkg.in/gomail.v2"
)

var (
	host     string = "smtp.qiye.163.com"
	port     int    = 25
	username string = "php01@2009400.cn"
	password string = "ZQxx400@"
	from     string = "php01@2009400.cn"
)

func SendMail(subject, content string, to, cc []string) error {
	toList := make([]string, 0, len(to))
	ccList := make([]string, 0, len(cc))

	for _, v := range to {
		v = strings.TrimSpace(v)
		if libs.IsEmail([]byte(v)) {
			exists := false
			for _, vv := range toList {
				if v == vv {
					exists = true
					break
				}
			}
			if !exists {
				toList = append(toList, v)
			}
		}
	}
	for _, v := range cc {
		v = strings.TrimSpace(v)
		if libs.IsEmail([]byte(v)) {
			exists := false
			for _, vv := range ccList {
				if v == vv {
					exists = true
					break
				}
			}
			if !exists {
				ccList = append(ccList, v)
			}
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(from, "Bob"))
	fmtToList := make([]string, 0)
	for _, val := range toList {
		fmtToList = append(fmtToList, m.FormatAddress(val, "Alice"))
	}
	m.SetHeader("To", fmtToList...)
	if len(ccList) > 0 {
		m.SetHeader("Cc", ccList...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewPlainDialer(host, port, username, password)

	return d.DialAndSend(m)
}

// 获取当前机器的IP
func GetCurrentMachineAddr() (ipAddr string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Println(ipnet.IP.String())
				ipAddr = ipnet.IP.String()
				return
			}
		}
	}
	return
}

func main() {
	subject := `subject`
	content := `This is content`
	to := []string{
		`php01@2009400.cn`,
	}
	cc := []string{}
	if err := SendMail(subject, content, to, cc); err != nil {
		log.Panic(err)
	} else {
		log.Println("send mail success")
	}
}
