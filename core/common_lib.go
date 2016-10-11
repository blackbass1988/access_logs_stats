package core

import (
	"regexp"
)

func NewRow(rawString string, rex *regexp.Regexp) (row *Row, err error) {
	row = new(Row)
	row.Fields = make(map[string]string)
	row.Raw = rawString

	matches := rex.FindStringSubmatch(rawString)

	if len(matches) == 0 {
		return nil, ERR_EMPTY_RESULT
	}

	for i, name := range rex.SubexpNames() {
		if len(name) > 0 {
			row.Fields[name] = matches[i]
		}
	}
	return row, err
}
