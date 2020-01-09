package funcs

import (
	"errors"
	"log"
	"os"
	"reflect"
	"sync"
)

//LogPrefix is the prefix of log.
const LogPrefix = "funcs"

var (
	//ErrNumParams is the error of params number.
	ErrNumParams = errors.New("The number of params is not adapted")
	//ErrObject is the error of nil.
	ErrObject = errors.New("The object is nil")
	//DefalutFuncs is the defalut Funcs.
	DefalutFuncs = New()
	logger       = log.New(os.Stdout, "["+LogPrefix+"] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
)

//Funcs defines the set of Func.
type Funcs struct {
	m     sync.Map
	isLog bool
}

//Func defines a method of struct .
type Func struct {
	StructName string
	MethodName string
	Value      reflect.Value
	Type       reflect.Type
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
	tp := reflect.TypeOf(obj)
	vf := reflect.ValueOf(obj)
	if name == "" {
		name = reflect.Indirect(vf).Type().Name()
	}
	var pname string
	if len(name) > 0 {
		pname = name + "."
	}
	nm := vf.NumMethod()
	f.logPrintf("StructName:%s,NumMethod:%d", name, nm)
	for i := 0; i < nm; i++ {
		callName := pname + vf.Type().Method(i).Name
		f.logPrintf("MethodIndex:%d,CallName:%s,NumIn:%d,NumOut:%d", i, callName, vf.Method(i).Type().NumIn(), vf.Method(i).Type().NumOut())
		Func := &Func{
			StructName: pname,
			MethodName: vf.Type().Method(i).Name,
			Value:      vf.Method(i),
			Type:       tp.Method(i).Type,
		}
		f.m.Store(callName, Func)
	}
	return nil
}

// Call calls the function v with the input arguments in.
// For example, Call("v",arg1,arg2) represents the Go call v(arg1,arg2).
// Call panics if v's Kind is not Func.
// As in Go, each input argument must be assignable to the
// type of the function's corresponding input parameter.
func Call(name string, params ...interface{}) (err error) {
	return DefalutFuncs.Call(name, params...)
}

// Call calls the function v with the input arguments in.
// For example, Call("v",arg1,arg2) represents the Go call v(arg1,arg2).
// Call panics if v's Kind is not Func.
// As in Go, each input argument must be assignable to the
// type of the function's corresponding input parameter.
func (f *Funcs) Call(name string, params ...interface{}) (err error) {
	var F *Func
	if F = f.GetFunc(name); F == nil {
		err = errors.New(name + " is not existed")
		return
	}
	if len(params) != F.NumIn() {
		err = ErrNumParams
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	if F.Value.Type().NumOut() > 0 {
		if F.Value.Type().Out(0).Name() == "error" {
			vs := F.Value.Call(in)
			if vs[0].IsNil() {
				return nil
			}
			return vs[0].Interface().(error)
		}
	}
	F.Value.Call(in)
	return
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
	return reflect.New(F.Type.In(index).Elem()).Interface()
}

//NumIn returns the number of input parameter.
func (f *Func) NumIn() int {
	return f.Value.Type().NumIn()
}

//SetLog can enable Log in the DefalutFuncs.
func SetLog(enable bool) {
	DefalutFuncs.SetLog(enable)
}

//SetLog can enable Log in the Funcs.
func (f *Funcs) SetLog(enable bool) {
	f.isLog = enable
}

func (f *Funcs) logPrintf(format string, args ...interface{}) {
	if f.isLog {
		logger.Printf(format, args...)
	}
}
