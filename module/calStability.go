package module

import "ngxlog"

//4xx
func Is4xx(entry *ngxlog.Entry) bool {
	status, _ := entry.Int64Field("status")
	return status >= 400 && status < 500
}

//5xx
func Is5xx(entry *ngxlog.Entry) bool {
	status, _ := entry.Int64Field("status")
	return status >= 500
}
