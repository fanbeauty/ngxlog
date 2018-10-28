package ngxlog

type Rets []Ret;

type Ret struct {
	Subject     string;
	Project     string;
	Date        string;
	Stability   string;
	P99         string;
	Over1minute string;
	Qps         string;
}
