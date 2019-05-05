package funcs

import (
	"errors"
	"reflect"
)

var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
)

type FuncMap map[string]*Func

type Func struct {
	Value 		reflect.Value
	Type 		reflect.Type
}

var DefalutFuncs FuncMap

func init() {
	DefalutFuncs=make(FuncMap)
}
func New() FuncMap {
	Funcs:=make(FuncMap)
	return Funcs
}
func Register(obj interface{}) (err error) {
	return DefalutFuncs.Register(obj)
}
func (f FuncMap)Register(obj interface{}) (err error) {
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
	//fmt.Println("NumMethod:", mNum)
	for i := 0; i < mNum; i++ {
		mName := name+vft.Method(i).Name
		//fmt.Println("index:", i, " MethodName:", mName)
		method := typ.Method(i)
		//mtype := method.Type
		//fmt.Println(mtype,method.Name,mtype.NumIn,mtype.In(1),mtype.In(2),mtype.NumOut(),mtype.Out(0))
		Func:=&Func{
			Value:vf.Method(i),
			Type:method.Type,
		}
		f[mName]=Func
	}
	return nil
}
func Call(name string, params ...interface{}) ( err error) {
	return DefalutFuncs.Call(name,params...)
}

func (f FuncMap)Call(name string, params ...interface{}) ( err error) {
	if _, ok := f[name]; !ok {
		err = errors.New(name + " does not exist.")
		return
	}
	if len(params) != f[name].Value.Type().NumIn() {
		err = ErrParamsNotAdapted
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	f[name].Value.Call(in)
	return
}
