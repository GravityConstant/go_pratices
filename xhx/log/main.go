package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type logFileWriter struct {
	file *os.File
	//write count
	size int64
}

var wg sync.WaitGroup

func main() {
	nowStr := time.Now().Format("2006-01-02-15-04-05")
	fmt.Println(nowStr)
}
func demo() {
	//log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("./mylog"+strconv.FormatInt(time.Now().Unix(), 10), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		log.Fatal("log  init failed")
	}

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileWriter := logFileWriter{file, info.Size()}
	log.SetOutput(&fileWriter)
	log.Info("start.....")
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go logTest(i)
	}
	log.Warn("waitting...")
	wg.Wait()
}
func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	n, e := p.file.Write(data)
	p.size += int64(n)
	//文件最大 64K byte
	if p.size > 1024*64 {
		p.file.Close()
		fmt.Println("log file full")
		p.file, _ = os.OpenFile("./mylog"+strconv.FormatInt(time.Now().Unix(), 10), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		p.size = 0
	}
	return n, e
}

func logTest(id int) {
	for i := 0; i < 100; i++ {
		log.Info("Thread:", id, " value:", i)
		time.Sleep(10 * time.Millisecond)
	}
	wg.Done()
}
