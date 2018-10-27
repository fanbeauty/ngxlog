package module

import "ngxlog"

//插入排序
func InsertSort(req ngxlog.ReqDataSet, entry *ngxlog.Entry) (ngxlog.ReqDataSet) {
	req = append(req, entry)
	tmp := entry.Float64Field("request_time")
	length := len(req)
	j := length - 2
	for j >= 0 && req[j].Float64Field("request_time") >= tmp {
		req[j+1] = req[j]
		j--
	}
	req[j+1] = entry
	return req
}
