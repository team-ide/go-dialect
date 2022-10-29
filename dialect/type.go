package dialect

import (
	"strconv"
	"strings"
)

type FuncTypeInfo struct {
	Name   string `json:"name,omitempty"`
	Format string `json:"format,omitempty"`
}

type ColumnTypeInfo struct {
	Name       string `json:"name,omitempty"`
	TypeFormat string `json:"typeFormat,omitempty"`
	HasLength  bool   `json:"hasLength,omitempty"`
	HasDecimal bool   `json:"hasDecimal,omitempty"`
	IsNumber   bool   `json:"isNumber,omitempty"`
	IsString   bool   `json:"isString,omitempty"`
	IsDateTime bool   `json:"isDateTime,omitempty"`
	IsBytes    bool   `json:"isBytes,omitempty"`
	MinLength  int    `json:"minLength,omitempty"`
	MaxLength  int    `json:"maxLength,omitempty"`
	FormatFunc func(length int, decimal int) (columnType string)
}

func (this_ *ColumnTypeInfo) FormatColumnType(length int, decimal int) (columnType string) {
	if this_.FormatFunc != nil {
		return this_.FormatFunc(length, decimal)
	}
	columnType = this_.TypeFormat
	lStr := ""
	dStr := ""
	if length >= 0 {
		lStr = strconv.Itoa(length)
	}
	if decimal >= 0 {
		dStr = strconv.Itoa(decimal)
	}
	if length == 0 && decimal == 0 {
		lStr = ""
		dStr = ""
	}
	columnType = strings.ReplaceAll(columnType, "$l", lStr)
	columnType = strings.ReplaceAll(columnType, "$d", dStr)
	columnType = strings.ReplaceAll(columnType, " ", "")
	columnType = strings.ReplaceAll(columnType, ",)", ")")
	columnType = strings.TrimSuffix(columnType, "(,)")
	columnType = strings.TrimSuffix(columnType, "()")
	return
}
