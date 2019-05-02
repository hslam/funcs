package main

import (
	"hslam.com/mgit/Mort/funcs"
	"fmt"
	"hslam.com/mgit/Mort/funcs/test/pb"
	"errors"
	"log"
	proto "github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
)

type Arith struct {
}


// 乘法运算方法
func (this *Arith) Multiply(req *pb.ArithRequest, res *pb.ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
}

// 除法运算方法
func (this *Arith) Divide(req *pb.ArithRequest, res *pb.ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

func main() {
	var err error
	//client send req to Server
	req := pb.ArithRequest{A:9,B:2}

	req_bytes, _ := proto.Marshal(&req)
	rpc_req:=pb.Request{Method:"Arith.Multiply",Data:req_bytes}
	rpc_req_bytes, _ := proto.Marshal(&rpc_req)

	//Server rec req
	var rpc_req_decode pb.Request
	if err := proto.Unmarshal(rpc_req_bytes, &rpc_req_decode); err != nil {
	}
	var req_decode pb.ArithRequest
	if err := proto.Unmarshal(rpc_req_decode.Data, &req_decode); err != nil {
	}
	//Server to do
	var res pb.ArithResponse
	fmt.Println(rpc_req_decode.Method,req_decode)
	funcs.Register(new(Arith))
	err = funcs.Call(rpc_req_decode.Method, &req_decode, &res) // 乘法运算
	//Sever send res to client
	res_bytes, _ := proto.Marshal(&res)
	rpc_res:= pb.Response{Data:res_bytes}

	rpc_res_bytes, _ := proto.Marshal(&rpc_res)
	//client rec res
	var rpc_res_decode pb.Response
	if err := proto.Unmarshal(rpc_res_bytes, &rpc_res_decode); err != nil {
	}
	var res_decode pb.ArithResponse
	if err := proto.Unmarshal(rpc_res_decode.Data, &res_decode); err != nil {
	}
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res_decode.Pro)
	//if err != nil {
	//	log.Fatalln("arith error: ", err)
	//}
	//fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	//
	//err = funcs.Call("Arith.Divide", &req, &res)
	//if err != nil {
	//	log.Fatalln("arith error: ", err)
	//}
	//fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)

	}


type Request struct {
	pb *pb.Request
}


func (req *Request) Encode(w io.Writer) (int, error) {
	pb := &pb.Request{
	}
	p, err := proto.Marshal(pb)
	if err != nil {
		return -1, err
	}
	return w.Write(p)
}


func (req *Request) Decode(r io.Reader) (int, error) {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		return -1, err
	}

	pb := new(pb.Request)
	if err := proto.Unmarshal(data, pb); err != nil {
		return -1, err
	}
	return len(data), nil
}
type Response struct {
	MethodName string
	pb *pb.Response
}


func (resp *Response) Encode(w io.Writer) (int, error) {
	b, err := proto.Marshal(resp.pb)
	if err != nil {
		return -1, err
	}
	return w.Write(b)
}


func (resp *Response) Decode(r io.Reader) (int, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return -1, err
	}
	resp.pb = new(pb.Response)
	if err := proto.Unmarshal(data, resp.pb); err != nil {
		return -1, err
	}
	return len(data), nil
}
