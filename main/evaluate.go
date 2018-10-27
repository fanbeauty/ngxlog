package main

import (
	"flag"
	"io"
	"os"
	"bufio"
	"ngxlog"
	"ngxlog/module"
	"sort"
	"strings"
	"strconv"
	"ngxlog/mail"
)

const (
	format = `$remote_addr - $remote_user [$time_local] "$request"` +
		` "$http_referer" "$http_user_agent" "$http_cookie" "$http_x_forwarded_for"` +
		` $status $body_bytes_sent $request_time`
	subject = "TengYue Infra Team"
)

//define flag
var logFile string
var interval int64
var emailAddress string
var project string

//define status count
var reqCount int64
var over1MinuteReq int
var xx5Count int64
var xx4Count int64
var durationCount int64
var failure float64
var stability float64
var over1MinuteReqPercentage float64
var p99 float64

//dataSet
var reqDataSet ngxlog.ReqDataSet

//定义时间区间
var start int64
var end int64

func init() {
	flag.StringVar(&logFile, "log", "-", "Specify the log file. Read from STDIN if the file name is '-'")
	flag.Int64Var(&interval, "interval", 60, "Specify the interval to calculate system stability. default 60 seconds")
	flag.StringVar(&emailAddress, "to", "meichaofan@tengyue360.com", "Specify the email address. If you want to specify multiple recipients, use a ',' separated string")
	flag.StringVar(&project, "project", "-", "Specify the business thread.the value is must")
}

func main() {
	flag.Parse()
	if project == "-" {
		panic("you must specify the business name")
		os.Exit(1)
	}
	var logReader io.Reader
	if logFile == "-" {
		logReader = os.Stdin
	} else {
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		logReader = file
		defer file.Close()
	}
	statistics(logReader)

	data := map[string]string{
		"subject":     subject,
		"project":     project,
		"date":        module.GetYesDate(),
		"stability":   strconv.FormatFloat(stability, 'f', 6, 64),
		"p99":         strconv.FormatFloat(p99, 'f', 6, 64),
		"over1minute": strconv.FormatFloat(over1MinuteReqPercentage, 'f', 6, 64),
		"qps":         strconv.FormatInt(reqCount, 10),
	}

	//content := fmt.Sprintf("稳定性:%.6f\nP99:%.6f\n>>1s:%.6f\nQPS:%v\n", stability, p99, over1MinuteReqPercentage, reqCount)
	to := strings.Split(emailAddress, ",")
	//发送邮件
	//module.Send(subject, content, to)

	r := mail.NewRequest(to, subject)
	r.Send("../templates/template.html", data)
}

//统计
func statistics(f io.Reader) {
	scanner := bufio.NewScanner(f)
	parse := ngxlog.NewParser(format)

	//先取得第一条日志
	scanner.Scan()
	reqCount++
	firstEntry, _ := parse.ParseString(scanner.Text())
	reqDataSet = append(reqDataSet, firstEntry)
	start = module.TransferNgxStartTs2UnixTs(firstEntry)
	end = start + interval
	statReqTsMoreThan1Second(firstEntry)
	statStability(firstEntry)

	//统计后续的日志
	for scanner.Scan() {
		entry, _ := parse.ParseString(scanner.Text())
		reqDataSet = append(reqDataSet, entry)

		statReqTsMoreThan1Second(entry)
		statStability(entry)
		reqCount++
	}
	//计算 request_time > 1s 的占比
	over1MinuteReqPercentage = float64(over1MinuteReq) / float64(reqCount)
	//计算稳定性
	stability = 1 - failure*float64(interval)/float64(86400)
	statP99(reqDataSet)
}

//统计request_time > 1s
func statReqTsMoreThan1Second(entry *ngxlog.Entry) {
	if requestTime := entry.Float64Field("request_time"); requestTime >= 1 {
		over1MinuteReq++
	}
}

//统计稳定性数据
func statStability(entry *ngxlog.Entry) {
	//稳定性计算
	curTime := module.TransferNgxTs2UnixTs(entry)
	if curTime >= start && curTime < end {
		if module.Is4xx(entry) {
			xx4Count++
		} else if module.Is5xx(entry) {
			xx5Count++
		}
		durationCount++
	} else {
		//新的开始
		if durationCount-xx4Count != 0 {
			failure += float64(xx5Count) / (float64(durationCount) - float64(xx4Count))
		}

		/*
		timeLayout := "2006-01-02 15:04:05"  //转化所需模板
		startTs := time.Unix(start, 0).Format(timeLayout)
		endTs := time.Unix(end, 0).Format(timeLayout)
		fmt.Printf("start:%v - end:%v 4xx:%v 5xx:%v qps:%v\n", startTs, endTs, xx4Count, xx5Count, durationCount)
		*/

		//初始化数据
		start = end
		end = start + interval
		xx5Count = 0
		xx4Count = 0
		durationCount = 0
	}
}

//计算P99
func statP99(set ngxlog.ReqDataSet) {
	sort.Stable(set)
	threshold99 := len(set) * 99 / 100
	requestTime := 0.0
	for i := 0; i <= threshold99; i++ {
		requestTime += reqDataSet[i].Float64Field("request_time")
	}
	p99 = requestTime / float64(threshold99)
}
