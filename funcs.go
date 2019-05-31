package funcs

import (
	"errors"
	"reflect"
	"fmt"
	"sync"
	"time"
)

const LogName  = "funcs"
var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
	DefalutFuncs *Funcs
)

type Funcs struct {
	m sync.Map
	isLog bool
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
	var pname string
	if len(name)>1{
		pname=name+"."
	}
	vft := vf.Type()
	mNum := vf.NumMethod()
	f.logPrintln("StructName:", name,"|| NumMethod:", mNum)
	for i := 0; i < mNum; i++ {
		mName := pname+vft.Method(i).Name
		method := typ.Method(i)
		mtype := method.Type
		//logPrintln(mtype,method.Name,mtype.NumIn(),mtype.In(0),mtype.In(1),mtype.In(2),mtype.NumOut(),mtype.Out(0))
		f.logPrintln("MethodIndex:", i, "|| CallName:", mName,"|| MethodName:",method.Name,"|| NumInParam:",mtype.NumIn()-1,"|| NumOutResult:",mtype.NumOut())
		Func:=&Func{
			Value:vf.Method(i),
			Type:mtype,
		}
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
	if F.Value.Type().NumOut()>0 {
		if F.Value.Type().Out(0).Name()=="error"{
			vs:=F.Value.Call(in)
			if vs[0].IsNil(){
				return nil
			}else {
				return vs[0].Interface().(error)
			}
		}
	}
	F.Value.Call(in)
	return
}
func GetFunc(name string) (F *Func) {
	return DefalutFuncs.GetFunc(name)
}
func (f *Funcs)GetFunc(name string) (F *Func) {
	if v, ok := f.m.Load(name); !ok {
		return
	}else {
		F=v.(*Func)
	}
	return
}
func EnabledLog() {
	DefalutFuncs.EnabledLog()
}
func (f *Funcs)EnabledLog() {
	f.isLog=true
}
func (f *Funcs)logPrintln(args ...interface{}) {
	if f.isLog{
		logargs:=make([]interface{},1)
		logargs[0]="["+LogName+"] "+time.Now().Format("2006/01/02-15:04:05")+" ||"
		logargs=append(logargs,args...)
		fmt.Println(logargs...)
	}
}
