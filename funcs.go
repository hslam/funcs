package funcs

import (
	"errors"
	"reflect"
	"log"
	"sync"
)

var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
	DefalutFuncs *Funcs
	IsLog =false
)

type Funcs struct {
	m sync.Map
}
type Func struct {
	Value 		reflect.Value
	Type 		reflect.Type
}
func init() {
	DefalutFuncs=new(Funcs)
}
func New() *Funcs {
	return new(Funcs)
}
func Register(obj interface{}) (err error) {
	return DefalutFuncs.Register(obj)
}
func (f *Funcs)Register(obj interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New( " is not callable.")
		}
	}()
	typ:=reflect.TypeOf(obj)
	vf := reflect.ValueOf(obj)
	name := reflect.Indirect(vf).Type().Name()
	if len(name)>1{
		name=name+"."
	}
	vft := vf.Type()
	mNum := vf.NumMethod()
	logPrintln("NumMethod:", mNum)
	for i := 0; i < mNum; i++ {
		mName := name+vft.Method(i).Name
		logPrintln("index:", i, " MethodName:", mName)
		method := typ.Method(i)
		mtype := method.Type
		logPrintln(mtype,method.Name,mtype.NumIn,mtype.In(1),mtype.In(2),mtype.NumOut(),mtype.Out(0))
		Func:=&Func{
			Value:vf.Method(i),
			Type:method.Type,
		}
		//f[mName]=Func
		f.m.Store(mName,Func)
	}
	return nil
}
func Call(name string, params ...interface{}) ( err error) {
	return DefalutFuncs.Call(name,params...)
}

func (f *Funcs)Call(name string, params ...interface{}) ( err error) {
	var F *Func
	if F = f.GetFunc(name); F==nil {
		err = errors.New(name + " is not existed.")
		return
	}
	if len(params) != F.Value.Type().NumIn() {
		err = ErrParamsNotAdapted
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	F.Value.Call(in)
	return
}
func (f *Funcs)GetFunc(name string) (F *Func) {
	if v, ok := f.m.Load(name); !ok {
		return
	}else {
		F=v.(*Func)
	}
	return
}
func logPrintln(args ...interface{}) {
	if IsLog{
		log.Println(args...)
	}
}