package self_rpc

import (
	"net"
	"reflect"
)

// client
type Client struct {
	conn net.Conn
}

// new
func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

// call rpc
func (c *Client) callRPC(rpcName string, fPtr interface{}) {
	// get original
	fn := reflect.ValueOf(fPtr).Elem()
	// fn params
	f := func(args []reflect.Value) []reflect.Value {
		// handle params
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		// connect to server
		cliSession := NewSession(c.conn)
		// encode data
		reqRPC := RPCData{Name: rpcName, Args: inArgs}
		b, err := encode(reqRPC)
		if err != nil {
			panic(err)
		}
		// write data
		err = cliSession.Write(b)
		if err != nil {
			panic(err)
		}
		// read data
		respBytes, err := cliSession.Read()
		if err != nil {
			panic(err)
		}
		// decode
		respRPC, err := decode(respBytes)
		if err != nil {
			panic(err)
		}
		// handle data
		outArgs := make([]reflect.Value, 0, len(respRPC.Args))
		for i, arg := range respRPC.Args {
			// nil transform
			if arg == nil {
				// reflect.Zero()会返回类型的零值的value
				// .out()会返回函数输出的参数类型
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(arg))
		}
		return outArgs
	}
	v := reflect.MakeFunc(fn.Type(), f)
	fn.Set(v)
}
