package jsonutil

import (
	"fmt"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

type JsonPath struct {
	object interface{}
}

func NewJsonPath(data string) (*JsonPath, error) {
	ob, err := oj.ParseString(data)
	if err != nil {
		return nil, err
	}
	return &JsonPath{object: ob}, nil
}

func (j *JsonPath) Get(path string) interface{} {
	x, err := jp.ParseString(path)
	if err != nil {
		log.Warn("Invalid JSON path: ", err)
		return nil
	}
	val := x.Get(j.object)
	if val == nil {
		return nil
	}
	return val[0]
}

func (j *JsonPath) GetString(path string) string {
	val := j.Get(path)
	if val == nil {
		return ""
	}
	return val.(string)
}

func (j *JsonPath) GetStringf(path string, args ...interface{}) string {
	return j.GetString(fmt.Sprintf(path, args...))
}

func (j *JsonPath) GetStringPtr(path string) *string {
	val := j.Get(path)
	if val == nil {
		return nil
	}
	str := val.(string)
	return &str
}

func (j *JsonPath) GetStringPtrf(path string, args ...interface{}) *string {
	return j.GetStringPtr(fmt.Sprintf(path, args...))
}

func (j *JsonPath) GetFloat64(path string) float64 {
	val := j.Get(path)
	if val == nil {
		return 0
	}
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", val), 64)
	if err != nil {
		return 0
	}
	return f
}

func (j *JsonPath) GetFloat64f(path string, args ...interface{}) float64 {
	return j.GetFloat64(fmt.Sprintf(path, args...))
}

func (j *JsonPath) GetInt(path string) int {
	val := j.Get(path)
	if val == nil {
		return 0
	}
	i, err := strconv.Atoi(fmt.Sprintf("%v", val))
	if err != nil {
		return 0
	}
	return i
}

func (j *JsonPath) GetIntf(path string, args ...interface{}) int {
	return j.GetInt(fmt.Sprintf(path, args...))
}

func (j *JsonPath) GetIntPtr(path string) *int {
	val := j.Get(path)
	if val == nil {
		return nil
	}
	v := int(val.(int64))
	return &v
}

func (j *JsonPath) GetIntPtrf(path string, args ...interface{}) *int {
	return j.GetIntPtr(fmt.Sprintf(path, args...))
}

func (j *JsonPath) GetBool(path string) bool {
	val := j.Get(path)
	if val == nil {
		return false
	}
	return val.(bool)
}

func (j *JsonPath) GetBoolf(path string, args ...interface{}) bool {
	return j.GetBool(fmt.Sprintf(path, args...))
}
