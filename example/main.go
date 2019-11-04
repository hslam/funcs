package main
import (
	"hslam.com/git/x/funcs"
	"fmt"
	"errors"
	"log"
)
type ArithRequest struct {
	A int
	B int
}
type ArithResponse struct {
	Quo int		// quotient
	Rem int		// remainder
}

type Arith struct {
}
func (this *Arith) Divide(req *ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

func main() {
	Funcs:=funcs.New()
	Funcs.EnabledLog()
	Funcs.Register(new(Arith))
	req := &ArithRequest{A:9,B:2}	//req := Funcs.GetFuncIn("Arith.Divide",0).(*ArithRequest);req.A=9;req.B=2
	res :=new(ArithResponse)	//res :=Funcs.GetFuncIn("Arith.Divide",1).(*ArithResponse)
	if err := Funcs.Call("Arith.Divide", req, res);err != nil {
		log.Fatalln("Call Arith.Divide error: ", err)
		return
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}
