// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package funcs implements function call by its name.
package funcs

import (
	"errors"
	"reflect"
	"testing"
)

//ArithRequest is the first input parameter.
type ArithRequest struct {
	A int
	B int
}

//ArithResponse is the second input parameter.
type ArithResponse struct {
	Pro int // product
	Quo int // quotient
	Rem int // remainder
}

//Arith is the Service struct.
type Arith struct {
}

//Multiply is the Arith's Method.
func (a *Arith) Multiply(req *ArithRequest, res *ArithResponse) {
	res.Pro = req.A * req.B
}

//Divide is the Arith's Method.
func (a *Arith) Divide(req *ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("divide by zero")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}
func TestDefalutFuncs(t *testing.T) {
	SetLog(true)
	RegisterName("Arith", new(Arith))
	Register(new(Arith))

	f := GetFunc("Arith.Divide")
	if f.NumIn() != 2 {
		t.Errorf("%d\n", f.NumIn())
	}

	//req := &ArithRequest{A: 9, B: 2}
	req := GetFuncIn("Arith.Divide", 0).(*ArithRequest)
	req.A = 9
	req.B = 2

	//res := &ArithResponse{}
	res := GetFuncIn("Arith.Divide", 1).(*ArithResponse)

	if arg := GetFuncIn("Arith.Divide", 2); arg != nil {
		t.Error("arg is not nil", arg)
	}

	if err := Call("Arith.Divide", req, res); err != nil {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
	if res.Quo != 4 {
		t.Errorf("%d / %d, quo is %d\n", req.A, req.B, res.Quo)
	}
	if res.Rem != 1 {
		t.Errorf("%d / %d, rem is %d\n", req.A, req.B, res.Rem)
	}
}

func TestRegister(t *testing.T) {
	if err := Register(new(Arith)); err != nil {
		t.Error("Register error", err)
	}
	if err := Register(nil); err != ErrObject {
		t.Error("Register error", err)
	}
	Funcs := New()
	if err := Funcs.Register(new(Arith)); err != nil {
		t.Error("Register error", err)
	}
	if err := Funcs.Register(nil); err != ErrObject {
		t.Error("Register error", err)
	}
}

func TestRegisterName(t *testing.T) {
	if err := RegisterName("Arith", new(Arith)); err != nil {
		t.Error("Register error", err)
	}
	if err := RegisterName("Arith", nil); err != ErrObject {
		t.Error("Register error", err)
	}
	Funcs := New()
	if err := Funcs.RegisterName("Arith", new(Arith)); err != nil {
		t.Error("Register error", err)
	}
	if err := Funcs.RegisterName("Arith", nil); err != ErrObject {
		t.Error("Register error", err)
	}
}

func TestCall(t *testing.T) {
	Register(new(Arith))

	//req := &ArithRequest{A: 9, B: 2}
	req := GetFuncIn("Arith.Divide", 0).(*ArithRequest)
	req.A = 9
	req.B = 0
	f := GetFunc("Arith.Divide")
	//res := &ArithResponse{}
	res := f.GetIn(1).(*ArithResponse)

	if err := Call("", req, res); err.Error() != " is not existed" {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
	if err := Call("Arith.Divide", req); err != ErrNumParams {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
	if err := Call("Arith.Divide", req, res); err.Error() != "divide by zero" {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
	if err := Call("Arith.Multiply", req, res); err != nil {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
	if err := f.Call("Arith.Multiply", req, res); err != nil {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
}

func TestValueCall(t *testing.T) {
	Register(new(Arith))
	f := GetFunc("Arith.Multiply")
	//req := &ArithRequest{A: 9, B: 2}
	req := f.GetValueIn(0).Interface().(*ArithRequest)
	req.A = 9
	req.B = 0

	//res := &ArithResponse{}
	res := GetFuncValueIn("Arith.Multiply", 1).Interface().(*ArithResponse)

	if err := ValueCall("Arith.Multiply", ValueOf(req), ValueOf(res)); err != nil {
		t.Errorf("Call Arith.Divide error: %s\n", err.Error())
		return
	}
	if Value(ReflectValueOf(res)).Kind() != reflect.Ptr {
		t.Error("Value.Kind() error\n")
	}
}

func TestGetFuncValueIn(t *testing.T) {
	Register(new(Arith))
	if v := GetFuncValueIn("Arith.Multiply", 2); v != ZeroValue {
		t.Errorf("GetFuncValueIn error")
	}
}

func TestGetValueIn(t *testing.T) {
	Register(new(Arith))
	f := GetFunc("Arith.Multiply")
	if v := f.GetValueIn(2); v != ZeroValue {
		t.Errorf("GetValueIn error")
	}
}

func TestGetIn(t *testing.T) {
	Register(new(Arith))
	f := GetFunc("Arith.Multiply")
	if v := f.GetIn(2); v != nil {
		t.Errorf("GetValueIn error")
	}
}

func TestNumOut(t *testing.T) {
	Register(new(Arith))
	f := GetFunc("Arith.Divide")
	if n := f.NumIn(); n != 2 {
		t.Errorf("NumIn error")
	}
	if n := f.NumOut(); n != 1 {
		t.Errorf("NumOut error")
	}
	if n := f.NumCalls(); n != 0 {
		t.Errorf("NumCalls error")
	}
}
