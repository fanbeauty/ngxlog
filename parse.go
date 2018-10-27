package ngxlog

import (
	"regexp"
	"fmt"
	"strings"
)

type StringParse interface {
	ParseString(line string) (entry *Entry, err error)
}

type Parser struct {
	format string
	regexp *regexp.Regexp
}

func NewParser(format string) *Parser {
	re := regexp.MustCompile(`\\\$([A-Za-z0-9_]+)(\\?(.))`).ReplaceAllString(
		regexp.QuoteMeta(format+" "), "(?P<$1>[^$3]*)$2")
	return &Parser{format, regexp.MustCompile(fmt.Sprintf("^%v", strings.Trim(re, " ")))}
}

func (parser *Parser) ParseString(line string) (entry *Entry, err error) {
	re := parser.regexp
	fields := re.FindStringSubmatch(line)
	if fields == nil {
		err = fmt.Errorf("access log line '%v' does not match given format '%v'", line, re)
		return
	}

	entry = NewEmptyEntry()
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		entry.SetField(name, fields[i])
	}
	return
}
