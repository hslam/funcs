# funcs
[![GoDoc](https://godoc.org/github.com/hslam/funcs?status.svg)](https://godoc.org/github.com/hslam/funcs)
[![Build Status](https://travis-ci.org/hslam/funcs.svg?branch=master)](https://travis-ci.org/hslam/funcs)
[![linux](https://github.com/hslam/funcs/workflows/linux/badge.svg)](https://github.com/hslam/funcs/actions)
[![codecov](https://codecov.io/gh/hslam/funcs/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/funcs)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/funcs)](https://goreportcard.com/report/github.com/hslam/funcs)
[![GitHub release](https://img.shields.io/github/release/hslam/funcs.svg)](https://github.com/hslam/funcs/releases/latest)
[![LICENSE](https://img.shields.io/github/license/hslam/funcs.svg?style=flat-square)](https://github.com/hslam/funcs/blob/master/LICENSE)
Call a function by the function's name in golang.
## Get started

### Install
```
go get github.com/hslam/funcs
```
### Import
```
import "github.com/hslam/funcs"
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
func (s *Service) Method(params ...interface{}) error {
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
Funcs.SetLog(true)
```
if a function has 2 input parameters ,You can get the function's first input parameter and second input parameter.
```
Funcs.GetFuncIn("Service.Method",0)
Funcs.GetFuncIn("Service.Method",1)
//and so on
```
#### Example
```
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
```

#### Output
```
[funcs] 2020/01/09 10:12:19.065121 StructName:Arith,NumMethod:1
[funcs] 2020/01/09 10:12:19.065246 MethodIndex:0,CallName:Arith.Divide,NumIn:2,NumOut:1
num of args : 2
9 / 2, quo is 4, rem is 1
```
### Licence
This package is licenced under a MIT licence (Copyright (c) 2019 Meng Huang)


### Authors
funcs was written by Meng Huang.


