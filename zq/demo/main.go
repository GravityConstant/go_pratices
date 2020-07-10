package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func isProcessExist(appName string) (bool, string, int) {
	appary := make(map[string]int)
	cmd := exec.Command("cmd", "/C", "tasklist")
	output, _ := cmd.Output()
	//fmt.Printf("fields: %v\n", output)
	n := strings.Index(string(output), "System")
	if n == -1 {
		fmt.Println("no find")
		os.Exit(1)
	}
	data := string(output)[n:]
	fields := strings.Fields(data)
	for k, v := range fields {
		if v == appName {
			appary[appName], _ = strconv.Atoi(fields[k+1])

			return true, appName, appary[appName]
		}
	}

	return false, appName, -1
}

// var appName = "NWAY400.exe"
var appName = "notepad.exe"
var duration = 10 // 秒

// var appName = "notepad.exe"

func main() {
	log.Println("Now will be restart Nway400.exe every 30 minutes")
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println("open log.txt failed. ", err)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			f.WriteString(fmt.Sprintf("check NWAY400.exe failed. %s\r\n", err.(error).Error()))
			f.Close()
		}
	}()
	//now.Hour() == 2 && now.Minute() == 0 &&

	now := time.Now()
	fmt.Println("init:", now)
	time.Sleep(time.Duration(60 - now.Second() + 2))
	tick := time.Tick(time.Minute)
	for t := range tick {
		now := time.Now()
		fmt.Println(now)
		if now.Second() == 2 {
			log.Println("execute restart Nway400.exe. ", t)
			if exist, _, _ := isProcessExist(appName); exist {
				if err := stopProcess(); err != nil {
					f.WriteString(fmt.Sprintf("stop NWAY400.exe failed. %s\r\n", err.(error).Error()))
				}
			}
			time.Sleep(time.Duration(duration) * time.Second)
			if err := startProcess(); err != nil {
				f.WriteString(fmt.Sprintf("start NWAY400.exe failed. %s\r\n", err.(error).Error()))
			}
		}
	}
	log.Printf("\r\n------------------------------------------------------\r\n")
}

func startProcess() error {
	c := exec.Command("cmd.exe", "/C", "start", appName)
	c.Start()
	err := c.Wait()
	if err != nil {
		return err
	}
	return nil
}

func stopProcess() error {
	c := exec.Command("cmd.exe", "/C", "taskkill", "/IM", appName)
	c.Start()
	err := c.Wait()
	if err != nil {
		return err
	}
	return nil
}
