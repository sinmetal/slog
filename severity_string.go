// Code generated by "stringer -type Severity severity.go"; DO NOT EDIT.

package slog

import "fmt"

const _Severity_name = "DEBUGINFOWARNINGERROR"

var _Severity_index = [...]uint8{0, 5, 9, 16, 21}

func (i Severity) String() string {
	if i < 0 || i >= Severity(len(_Severity_index)-1) {
		return fmt.Sprintf("Severity(%d)", i)
	}
	return _Severity_name[_Severity_index[i]:_Severity_index[i+1]]
}
