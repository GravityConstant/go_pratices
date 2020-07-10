package self_rpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSession_ReadWriter(t *testing.T) {
	addr := "127.0.0.1:8000"
	my_data := "hello"

	wg := sync.WaitGroup{}
	wg.Add(2)

	// write
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, _ := lis.Accept()
		s := Session{conn: conn}
		err = s.Write([]byte(my_data))
		if err != nil {
			t.Fatal(err)
		}
	}()

	// read
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s := Session{conn: conn}
		data, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != my_data {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}()
	wg.Wait()
}
