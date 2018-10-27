package ngxlog

import (
	"fmt"
	"strconv"
)

//type Fields map[string]string
type Fields map[string]string

type Entry struct {
	fields Fields
}

func NewEmptyEntry() *Entry {
	return &Entry{make(Fields)}
}

func (entry *Entry) Fields() Fields {
	return entry.fields
}

func NewEntry(fields Fields) *Entry {
	return &Entry{fields}
}

func (entry *Entry) SetField(name string, value string) {
	entry.fields[name] = value
}

func (entry *Entry) GetField(name string) (string, error) {
	value, ok := entry.fields[name]
	if !ok {
		err := fmt.Errorf("field '%v' does not found in record %+v", name, *entry)
		return "", err
	}
	return value, nil
}

func (entry *Entry) Int64Field(name string) (value int64, err error) {
	tmp, err := entry.GetField(name)
	if err == nil {
		value, err = strconv.ParseInt(tmp, 0, 64)
	}
	return
}

func (entry *Entry) Float64Field(name string) (value float64) {
	tmp, err := entry.GetField(name)
	if err == nil {
		value, err = strconv.ParseFloat(tmp, 64)
	}
	return value
}

//存放结果集
type ReqDataSet []*Entry

func (ret ReqDataSet) Len() int {
	return len(ret)
}

func (ret ReqDataSet) Swap(i, j int) {
	ret[i], ret[j] = ret[j], ret[i]
}

func (ret ReqDataSet) Less(i, j int) bool {
	return ret[i].Float64Field("request_time") < ret[j].Float64Field("request_time")
}
