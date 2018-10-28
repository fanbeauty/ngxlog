# ngxlog
nginx 日志分析工具
### 编译生成可执行文件
```code
cd ngxlog/main/
go build evaluate.go
```
### 用法
```
Usage of ./evaluate:
  -interval int (指定间隔时间，统计系统稳定性)
    	Specify the interval to calculate system stability. default 60 seconds (default 60)
  -logs string (要分析的ngx日志文件，已逗号分割)
    	Specify the log file. use ',' separated if you want to analyse mlti log (default "-")
  -project string (对应日志文件的业务名称，以逗号分割)
    	Specify the business thread. please direct to logs (default "-")
  -to string (邮件接收人，如有多个，用逗号分隔)
    	Specify the email address. If you want to specify multiple recipients, use a ',' separated string (default "meichaofan@tengyue360.com")
```
### example
```
以5分钟为间隔，统计双师、APP业务系统的稳定性和其它指标
./evaluate -logs api.access.log.20181025,api-shuangshi.access.log.20181025 -project 腾跃APP业务,腾跃双师业务 -to meichaofan@tengyue360.com -interval 300
```
