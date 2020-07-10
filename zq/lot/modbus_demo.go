package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goburrow/modbus"
)

type DataHandle struct {
	Data  []byte
	Index []Element
	Last  Element
}

// 获取到的元素的信息
type Element struct {
	Len       int    // 所占字节数
	Formatter string // 格式化数据，如℃
	Sign      bool   // 结果是否是有符号的
}

func Bytes2Uint16(buf []byte) uint16 {
	if len(buf) > 2 {
		return 0
	}
	return uint16(buf[0])<<8 + uint16(buf[1])
}

func Bytes2Int16(buf []byte) int16 {
	if len(buf) > 2 {
		return 0
	}
	return int16(buf[0])<<8 + int16(buf[1])
}

func Bytes2Uint32(buf []byte) uint32 {
	if len(buf) > 4 {
		return 0
	}
	return uint32(buf[0])<<24 + uint32(buf[1])<<16 + uint32(buf[2])<<8 + uint32(buf[3])
}

func Bytes2Int32(buf []byte) int32 {
	if len(buf) > 4 {
		return 0
	}
	return int32(buf[0])<<24 + int32(buf[1])<<16 + int32(buf[2])<<8 + int32(buf[3])
}

func main() {
	modbusDemo()
}

func modbusDemo() {
	handler := modbus.NewRTUClientHandler("COM2")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	if err != nil {
		fmt.Println("connect error.", err)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	// 3s读一次数据
	tick := time.Tick(3 * time.Second)
	go func() {
		for range tick {
			// humiture(client)
			// getLux(client)
			// getLuxLong(client)
			getAll(client)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	fmt.Println("exit.")
}

// 字节转为人可读的信息
func handleData(data DataHandle) ([]int64, []string) {
	dataLen := len(data.Data)
	resLen := len(data.Index)
	if data.Last.Len > 0 {
		resLen += 1
	}
	res := make([]int64, resLen)
	restr := make([]string, resLen)

	j := 0 // 截取的byte的起始
	for i, d := range data.Index {
		oneData := data.Data[j : j+d.Len]
		oneDataLen := len(oneData)
		j = j + d.Len
		fmt.Printf("one data: %#v\n", oneData)
		switch oneDataLen {
		case 2:
			if d.Sign {
				res[i] = int64(Bytes2Int16(oneData))
			} else {
				res[i] = int64(Bytes2Uint16(oneData))
			}
		case 4:
			if d.Sign {
				res[i] = int64(Bytes2Int32(oneData))
			} else {
				res[i] = int64(Bytes2Uint32(oneData))
			}
		}
		restr[i] = fmt.Sprintf("%d%s", res[i], d.Formatter)
	}
	// last element
	oneDataLen := dataLen - data.Last.Len
	oneData := data.Data[oneDataLen:]
	fmt.Printf("last data: %#v, len: %d\n", oneData, oneDataLen)
	switch data.Last.Len {
	case 2:
		if data.Last.Sign {
			res[resLen-1] = int64(Bytes2Int16(oneData))
		} else {
			res[resLen-1] = int64(Bytes2Uint16(oneData))
		}
	case 4:
		if data.Last.Sign {
			res[resLen-1] = int64(Bytes2Int32(oneData))
		} else {
			res[resLen-1] = int64(Bytes2Uint32(oneData))
		}
	}
	restr[resLen-1] = fmt.Sprintf("%d%s", res[resLen-1], data.Last.Formatter)

	return res, restr
}

// 温湿度获取
func humiture(client modbus.Client) {
	results, err := client.ReadHoldingRegisters(0, 2)
	if err != nil {
		fmt.Println("read humiture error.", err)
		return
	}
	// fmt.Printf("humiture: %v\n", results)

	dh := DataHandle{}
	dh.Data = results
	dh.Index = []Element{
		Element{2, "RH", false},
		Element{2, "℃", true},
	}
	// fmt.Printf("humiture dh.last: %#v\n", dh.Last)
	res, restr := handleData(dh)
	fmt.Println("humiture:", res, restr)
}

// 读取光照度0-65535
func getLux(client modbus.Client) {
	results, err := client.ReadHoldingRegisters(6, 1)
	if err != nil {
		fmt.Println("read lux error.", err)
		return
	}
	// fmt.Printf("lux: %v\n", results)

	dh := DataHandle{}
	dh.Data = results
	dh.Index = []Element{
		Element{2, "Lux", false},
	}
	// fmt.Printf("lux dh.last: %#v\n", dh.Last)
	res, restr := handleData(dh)
	fmt.Println("lux:", res, restr)
}

// 读取光照度0-200000
func getLuxLong(client modbus.Client) {
	results, err := client.ReadHoldingRegisters(2, 2)
	if err != nil {
		fmt.Println("read lux long error.", err)
		return
	}
	// fmt.Printf("lux: %v\n", results)

	dh := DataHandle{}
	dh.Data = results
	dh.Index = []Element{
		Element{4, "Lux", false},
	}
	// fmt.Printf("lux dh.last: %#v\n", dh.Last)
	res, restr := handleData(dh)
	fmt.Println("lux long:", res, restr)
}

// 读取所有信息
func getAll(client modbus.Client) {
	results, err := client.ReadHoldingRegisters(0, 7)
	if err != nil {
		fmt.Println("read get all error.", err)
		return
	}
	fmt.Printf("get all: %v\n", results)

	dh := DataHandle{}
	dh.Data = results
	dh.Index = []Element{
		Element{2, "RH", false},
		Element{2, "℃", true},
	}
	dh.Last = Element{2, "Lux", false}
	fmt.Printf("get all dh.last: %#v\n", dh.Last)
	res, restr := handleData(dh)
	fmt.Println("get all:", res, restr)
}
