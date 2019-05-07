package main

import (
	"hslam.com/mgit/Mort/funcs"
	"fmt"
	"hslam.com/mgit/Mort/funcs/test/pb"
	"errors"
	"log"
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
	funcs.Register(new(Arith))
	var err error
	req := pb.ArithRequest{A:9,B:2}
	var res pb.ArithResponse
	err = funcs.Call("Arith.Multiply", &req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	err = funcs.Call("Arith.Divide", &req, &res)
	if err != nil {
		log.Fatalln("arith error: ", err)
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
	}

