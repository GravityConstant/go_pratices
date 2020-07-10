package self_rpc

import (
	"fmt"
	"net"
	"reflect"
)

// Server
type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

// new
func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

// register
func (s *Server) Register(rpcName string, f interface{}) {
	if _, ok := s.funcs[rpcName]; ok {
		return
	}
	fVal := reflect.ValueOf(f)
	s.funcs[rpcName] = fVal
}

// server
func (s *Server) Run() {
	// listen
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Printf("监听 %s err: %v", s.addr, err)
		return
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			return
		}
		serSession := NewSession(conn)
		// rpc read
		b, err := serSession.Read()
		if err != nil {
			fmt.Println("read error: ", err)
			return
		}
		// decode
		rpcData, err := decode(b)
		if err != nil {
			fmt.Println("decode error: ", err)
			return
		}
		// get func name
		f, ok := s.funcs[rpcData.Name]
		if !ok {
			fmt.Printf("函数%s不存在\n", rpcData.Name)
			return
		}
		// parse args
		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for _, arg := range rpcData.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}
		// reflect call method
		out := f.Call(inArgs)
		// traverse result, put in slice
		outArgs := make([]interface{}, 0, len(out))
		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}
		// encode
		respRPCData := RPCData{rpcData.Name, outArgs}
		bytes, err := encode(respRPCData)
		if err != nil {
			fmt.Println("encode error. ", err)
			return
		}
		// write to client
		err = serSession.Write(bytes)
		if err != nil {
			fmt.Println("write error. ", err)
			return
		}
	}
}
