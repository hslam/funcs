# funcs
##A simple package to call a function by the function's name
Call a Struct and its Method given a string with the Struct and its Method's name in Golang

## Usage
### Here is how you use it:
First you need to create an instance of the funcs:
```
import "hslam.com/mgit/Mort/funcs"

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
if err = Funcs.Call("Service.Method", params...);err != nil {
    log.Fatalln("Call Service.Method error: ", err)
}
```


### Example
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
	Pro int // product
	Quo int // quotient
	Rem int // remainder
}

type Arith struct {
}

func (this *Arith) Multiply(req *ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil
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
	Funcs.Register(new(Arith))
	var(
		err error
		req ArithRequest
		res ArithResponse
	)
	req = ArithRequest{A:9,B:2}
	if err = Funcs.Call("Arith.Multiply", &req, &res);err != nil {
		log.Fatalln("arith multiply error: ", err)
	}else {
		fmt.Printf("%d * %d = %d\n", req.A, req.B, res.Pro)
	}
	if err = Funcs.Call("Arith.Divide", &req, &res);err != nil {
		log.Fatalln("arith divide error: ", err)
	}else {
		fmt.Printf("%d / %d, quo is %d, rem is %d\n", req.A, req.B, res.Quo, res.Rem)
	}
}

```

### Run Result
```
9 * 2 = 18
9 / 2, quo is 4, rem is 1
```
### Licence
This package is licenced under a MIT licence (Copyright (c) 2017 Arne Küderle)


### Authors
funcs was written by Mort Huang.


