# funcs
## A simple go package to call a function by the function's name
Call a Struct and its Method given a string with the Struct and its Method's name in Golang

## Get started

### Install
```
go get hslam.com/mgit/Mort/funcs
```
### Import
```
import "hslam.com/mgit/Mort/funcs"
```
### Usage
#### Here is how you use it:
First you need to create an instance of the funcs:
```
Funcs:=funcs.New()
```
Second you need to register your Struct:
```
type Service struct {
}
func (this *Service) Method(params ...interface{}) error {
    //to do
    return nil
}
Funcs.Register(new(Service))
```
And now you can call your function by name.

Function's Name Format : "StructName.MethodName"
```
if err := Funcs.Call("Service.Method", params...);err != nil {
    log.Fatalln("Call Service.Method error: ", err)
}
```
Logging.
```
Funcs.EnabledLog()
```

#### Example
```
package main
import (
	"hslam.com/mgit/Mort/funcs"
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
	req := ArithRequest{A:9,B:2}	//	req := ArithRequest{A:9,B:0}
	var res ArithResponse
	if err := Funcs.Call("Arith.Divide", &req, &res);err != nil {
		log.Fatalln("Call Arith.Divide error: ", err)
		return
	}
	fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
}
```

#### Output
```
[funcs] 2017/08/11-20:27:38 || NumMethod: 1
[funcs] 2017/08/11-20:27:38 || index: 0  MethodName: Arith.Divide
[funcs] 2017/08/11-20:27:38 || FuncName: Divide || NumInParam: 2 || NumOutResult: 1
9 / 2, quo is 4, rem is 1
```
### Licence
This package is licenced under a MIT licence (Copyright (c) 2017 Mort Huang)


### Authors
funcs was written by Mort Huang.


