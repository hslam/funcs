package main

import (
	"errors"
	"fmt"
	"github.com/hslam/funcs"
	"log"
)

type ArithRequest struct {
	A int
	B int
}

type ArithResponse struct {
	Quo int // quotient
	Rem int // remainder
}

type Arith struct {
}

func (a *Arith) Divide(req *ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

func main() {
	Funcs := funcs.New()
	Funcs.SetLog(true)

	//Funcs.RegisterName("Arith",new(Arith))
	Funcs.Register(new(Arith))

	f := Funcs.GetFunc("Arith.Divide")
	fmt.Printf("num of args : %d\n", f.NumIn())

	//req := &ArithRequest{A: 9, B: 2}
	req := Funcs.GetFuncIn("Arith.Divide", 0).(*ArithRequest)
	req.A = 9
	req.B = 2

	//res := &ArithResponse{}
	res := Funcs.GetFuncIn("Arith.Divide", 1).(*ArithResponse)

	if err := Funcs.Call("Arith.Divide", req, res); err != nil {
		log.Fatalln("Call Arith.Divide error: ", err)
		return
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}
