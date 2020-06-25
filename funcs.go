// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package funcs implements function call by its name.
package funcs

import (
	"errors"
	"log"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
)

// Value is the reflection interface to a Go value.
type Value reflect.Value

//LogPrefix is the prefix of log.
const LogPrefix = "funcs"

var (
	//ZeroValue is the Value of zero
	ZeroValue = Value{}
	//ErrNumParams is the error of params number.
	ErrNumParams = errors.New("The number of params is not adapted")
	//ErrObject is the error of nil.
	ErrObject = errors.New("The object is nil")
	//DefalutFuncs is the defalut Funcs.
	DefalutFuncs = New()
	logger       = log.New(os.Stdout, "["+LogPrefix+"] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
)

//Funcs defines the struct of Funcs.
type Funcs struct {
	m     sync.Map
	isLog bool
}

//Func defines the struct of func.
type Func struct {
	structName  string
	methodName  string
	structValue reflect.Value
	methodType  reflect.Type
	methodFunc  reflect.Value
	numIn       int
	numOut      int
	errorOut    bool
	count       int64
}

// New returns a new blank Funcs instance.
func New() *Funcs {
	return new(Funcs)
}

// Register publishes the set of struct's methods in the DefalutFuncs.
func Register(obj interface{}) (err error) {
	return DefalutFuncs.RegisterName("", obj)
}

// Register publishes the set of struct's methods in the Funcs.
func (f *Funcs) Register(obj interface{}) (err error) {
	return f.RegisterName("", obj)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the struct's concrete type.
func RegisterName(name string, obj interface{}) error {
	return DefalutFuncs.RegisterName(name, obj)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the struct's concrete type.
func (f *Funcs) RegisterName(name string, obj interface{}) (err error) {
	if obj == nil {
		return ErrObject
	}
	tf := reflect.TypeOf(obj)
	vf := reflect.ValueOf(obj)
	if name == "" {
		name = reflect.Indirect(vf).Type().Name()
	}
	nm := vf.NumMethod()
	f.logPrintf("StructName:%s,NumMethod:%d", name, nm)
	for i := 0; i < nm; i++ {
		Func := &Func{
			structName:  name,
			methodName:  vf.Type().Method(i).Name,
			structValue: vf,
			methodType:  tf.Method(i).Type,
			methodFunc:  vf.Type().Method(i).Func,
		}
		Func.numIn = vf.Method(i).Type().NumIn()
		Func.numOut = vf.Method(i).Type().NumOut()
		if Func.numOut > 0 && vf.Method(i).Type().Out(0).Name() == "error" {
			Func.errorOut = true
		}
		callName := Func.structName + "." + Func.methodName
		f.logPrintf("MethodIndex:%d,CallName:%s,NumIn:%d,NumOut:%d", i, callName, vf.Method(i).Type().NumIn(), vf.Method(i).Type().NumOut())
		f.m.Store(callName, Func)
	}
	return nil
}

// Call calls the function with the input arguments.
// For example, Call("v",arg1,arg2) represents the Go call v(arg1,arg2).
// Call panics if v's Kind is not Func.
// As in Go, each input argument must be assignable to the
// type of the function's corresponding input parameter.
func Call(name string, params ...interface{}) (err error) {
	return DefalutFuncs.Call(name, params...)
}

// Call calls the function with the input arguments.
// For example, Call("v",arg1,arg2) represents the Go call v(arg1,arg2).
// Call panics if v's Kind is not Func.
// As in Go, each input argument must be assignable to the
// type of the function's corresponding input parameter.
func (f *Funcs) Call(name string, params ...interface{}) (err error) {
	in := make([]Value, len(params))
	for k, param := range params {
		in[k] = Value(reflect.ValueOf(param))
	}
	return f.ValueCall(name, in...)
}

// ValueCall calls the function with the Value of input arguments.
func ValueCall(name string, in ...Value) (err error) {
	return DefalutFuncs.ValueCall(name, in...)
}

// ValueCall calls the function with the Value of input arguments.
func (f *Funcs) ValueCall(name string, in ...Value) (err error) {
	var F *Func
	if F = f.GetFunc(name); F == nil {
		err = errors.New(name + " is not existed")
		return
	}

	return F.ValueCall(in...)
}

//GetFunc returns Func by name in the DefalutFuncs.
func GetFunc(name string) (F *Func) {
	return DefalutFuncs.GetFunc(name)
}

//GetFunc returns Func by name in the Funcs.
func (f *Funcs) GetFunc(name string) (F *Func) {
	if v, ok := f.m.Load(name); ok {
		F = v.(*Func)
	}
	return
}

//GetFuncIn returns index'th input parameter by name and index in the DefalutFuncs.
func GetFuncIn(name string, i int) interface{} {
	return DefalutFuncs.GetFuncIn(name, i)
}

//GetFuncIn returns index'th input parameter by name and index in the Funcs.
func (f *Funcs) GetFuncIn(name string, i int) interface{} {
	index := i + 1
	F := f.GetFunc(name)
	if F == nil || index < 1 || index > F.NumIn() {
		return nil
	}
	return reflect.New(F.methodType.In(index).Elem()).Interface()
}

//GetFuncValueIn returns the Value of index'th input parameter by name and index in the DefalutFuncs.
func GetFuncValueIn(name string, i int) Value {
	return DefalutFuncs.GetFuncValueIn(name, i)
}

//GetFuncValueIn returns the Value of index'th input parameter by name and index in the Funcs.
func (f *Funcs) GetFuncValueIn(name string, i int) Value {
	index := i + 1
	F := f.GetFunc(name)
	if F == nil || index < 1 || index > F.NumIn() {
		return ZeroValue
	}
	return Value(reflect.New(F.methodType.In(index).Elem()))
}

//SetLog enables Log in the DefalutFuncs.
func SetLog(enable bool) {
	DefalutFuncs.SetLog(enable)
}

//SetLog enables Log in the Funcs.
func (f *Funcs) SetLog(enable bool) {
	f.isLog = enable
}

func (f *Funcs) logPrintf(format string, args ...interface{}) {
	if f.isLog {
		logger.Printf(format, args...)
	}
}

// Call calls the function with the input arguments.
func (f *Func) Call(params ...interface{}) (err error) {
	in := make([]Value, len(params))
	for k, param := range params {
		in[k] = Value(reflect.ValueOf(param))
	}
	f.ValueCall(in...)
	return
}

// ValueCall calls the function with the Value of input arguments.
func (f *Func) ValueCall(in ...Value) (err error) {
	if len(in) != f.NumIn() {
		err = ErrNumParams
		return
	}
	atomic.AddInt64(&f.count, 1)
	defer func() { atomic.AddInt64(&f.count, -1) }()
	ins := make([]reflect.Value, len(in)+1)
	ins[0] = f.structValue
	for k, param := range in {
		ins[k+1] = reflect.Value(param)
	}
	vs := f.methodFunc.Call(ins)
	if f.errorOut {
		if !vs[0].IsNil() {
			return vs[0].Interface().(error)
		}
	}
	return
}

//GetValueIn returns the Value of index'th input parameter by index.
func (f *Func) GetValueIn(i int) Value {
	index := i + 1
	if index < 1 || index > f.NumIn() {
		return Value{}
	}
	return Value(reflect.New(f.methodType.In(index).Elem()))
}

//GetIn returns index'th input parameter by index.
func (f *Func) GetIn(i int) interface{} {
	index := i + 1
	if index < 1 || index > f.NumIn() {
		return nil
	}
	return reflect.New(f.methodType.In(index).Elem()).Interface()
}

//NumIn returns the number of input parameter.
func (f *Func) NumIn() int {
	return f.numIn
}

//NumOut returns the number of output parameter.
func (f *Func) NumOut() int {
	return f.numOut
}

//NumCalls returns the number of calls.
func (f *Func) NumCalls() (n int64) {
	return atomic.LoadInt64(&f.count)
}

// ValueOf returns a new Value.
func ValueOf(param interface{}) Value {
	return Value(reflect.ValueOf(param))
}

// ReflectValueOf returns a new reflect.Value.
func ReflectValueOf(param interface{}) reflect.Value {
	return reflect.ValueOf(param)
}

// Interface returns v's current value as an interface{}.
func (v Value) Interface() (i interface{}) {
	return reflect.Value(v).Interface()
}

// Kind returns v's Kind.
// If v is the zero Value (IsValid returns false), Kind returns Invalid.
func (v Value) Kind() reflect.Kind {
	return reflect.Value(v).Kind()
}
