package main

import (
	"io"
	"io/ioutil"
	"log"

	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

/*
	函数名称：mk_dir
	函数作用：新建目录
	输入参数：dir_path（目录路径）
	输出参数：新建目录路径
*/

func mk_dir(dir_path string) string {
	var path string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	//fmt.Println(path)
	dir, _ := os.Getwd()                            //当前的目录
	err := os.Mkdir(dir+path+dir_path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		log.Println(err)
	}
	return dir + path + dir_path
}

/*
	函数名称：checkFileIsExist
	函数作用：检查文件是否存在，不存在则新建文件
	输入参数：filename（文件名）
	输出参数：是否新建成功
*/
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

/*
	函数名称：write_to_file
	函数作用：写内容到文件
	输入参数：filename（文件名），content（内容）
	输出参数：无
*/

func write_to_file(filename string, content string) {
	var f *os.File
	var err error
	if checkFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		log.Println("文件存在")
	} else {
		f, err = os.Create(filename) //创建文件
		log.Println("文件不存在")
	}
	check_error(err)
	_, err = io.WriteString(f, content) //写入文件(字符串)
	check_error(err)
}

/*
	函数名称：check_error
	函数作用：捕抓错误
	输入参数：error
	输出参数：无
*/
func check_error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func main() {
	arg_num := len(os.Args)
	if arg_num != 3 {
		log.Println("参数有错,请以此输入邮箱名称和密码，如./download.exe 123@qq.com 123123123")
		return
	}

	log.Println("Connecting to server...")

	// 连接服务器
	c, err := client.DialTLS("imap.qq.com:993", nil)
	check_error(err)
	log.Println("连接服务器")

	// 结束后退出登录
	defer c.Logout()

	// 登录
	//args[1]是用户名，args[2]是imap密码
	if err := c.Login(os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("登陆邮箱")

	// 获取邮箱列表
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("邮箱列表:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// 选择收件箱
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	from := uint32(1)
	to := mbox.Messages
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	attrs := []string{"BODY[]", imap.FlagsMsgAttr}
	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, attrs, messages)
	}()

	for msg := range messages {
		//判断是否未读邮件
		if len(msg.Flags) == 0 {
			r := msg.GetBody("BODY[]")
			if r == nil {
				log.Fatal("Server didn't returned message body")
			}

			mr, err := mail.CreateReader(r)
			if err != nil {
				log.Fatal(err)
			}
			header := mr.Header
			//获取邮件日期
			mail_date, errs := header.Date()
			check_error(errs)
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
				//提取邮件的附件
				switch h := p.Header.(type) {
				case mail.TextHeader:
					log.Println("just mail")
				case mail.AttachmentHeader:
					filename, _ := h.Filename()
					log.Println("Got attachment: %v", filename)
					dir := mk_dir(mail_date.Format("01-02-2006"))
					filename = dir + "\\" + filename
					content, _ := ioutil.ReadAll(p.Body)
					write_to_file(filename, string(content))

				}
			}
		} else {
			log.Println("已读邮件")
		}

	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

}
