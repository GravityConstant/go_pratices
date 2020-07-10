package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
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
	Ratio     int64  // 转换比率，负数为除，整数位乘
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
			tmp := binary.BigEndian.Uint16(oneData)
			if d.Sign {
				res[i] = int64(int16(tmp))
			} else {
				res[i] = int64(tmp)
			}
		case 4:
			tmp := binary.BigEndian.Uint32(oneData)
			if d.Sign {
				res[i] = int64(int32(tmp))
			} else {
				res[i] = int64(tmp)
			}
		}
		if d.Ratio > 0 {
			restr[i] = fmt.Sprintf("%d%s", res[i]*d.Ratio, d.Formatter)
		} else {
			restr[i] = fmt.Sprintf("%.1f%s", float64(res[i])/math.Abs(float64(d.Ratio)), d.Formatter)
		}

	}
	// check last element and handle
	if data.Last.Len > 0 {
		oneDataLen := dataLen - data.Last.Len
		oneData := data.Data[oneDataLen:]
		fmt.Printf("last data: %#v, len: %d\n", oneData, oneDataLen)
		switch data.Last.Len {
		case 2:
			tmp := binary.BigEndian.Uint16(oneData)
			if data.Last.Sign {
				res[resLen-1] = int64(int16(tmp))
			} else {
				res[resLen-1] = int64(tmp)
			}
		case 4:
			tmp := binary.BigEndian.Uint32(oneData)
			if data.Last.Sign {
				res[resLen-1] = int64(int32(tmp))
			} else {
				res[resLen-1] = int64(tmp)
			}
		}
		if data.Last.Ratio > 0 {
			restr[resLen-1] = fmt.Sprintf("%d%s", res[resLen-1]*data.Last.Ratio, data.Last.Formatter)
		} else {
			restr[resLen-1] = fmt.Sprintf("%.1f%s", float64(res[resLen-1])/math.Abs(float64(data.Last.Ratio)), data.Last.Formatter)
		}
	}

	return res, restr
}

func main() {
	handler := modbus.NewTCPClientHandler("127.0.0.1:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x01
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	if err != nil {
		fmt.Println("tcp connect error.", err)
		return
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	getAll(client)
}

// 读取所有信息
func getAll(client modbus.Client) {
	results, err := client.ReadHoldingRegisters(0, 5)
	if err != nil {
		fmt.Println("read get all error.", err)
		return
	}
	// fmt.Printf("get all: %v\n", results)

	dh := DataHandle{}
	dh.Data = results
	dh.Index = []Element{
		Element{2, "%RH", false, -10},
		Element{2, "℃", true, -10},
		Element{2, "Lux", false, 1},
		Element{4, "Lux", false, 1},
	}
	res, restr := handleData(dh)
	fmt.Println("get all:", res, restr)
}
