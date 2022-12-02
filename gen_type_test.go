package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/team-ide/go-dialect/dialect"
	"os"
	"strings"
	"testing"
)

func TestTypeParseGen(t *testing.T) {
	err := dataTypeParse(`数据库类型.xlsx`, "dialect/mapping.column.type.go")
	if err != nil {
		panic(err)
	}
}

type databaseModel struct {
	Name      string
	DataTypes []*dialect.ColumnTypeInfo
}

func dataTypeParse(path string, outPath string) (err error) {
	xlsxFForRead, err := xlsx.OpenFile(path)
	if err != nil {
		err = errors.New("excel [" + path + "] open error, " + err.Error())
		return
	}
	sheets := xlsxFForRead.Sheets

	var databases []*databaseModel

	for _, sheet := range sheets {
		database := &databaseModel{}
		database.Name = sheet.Name

		var titles []string

		var RowMergeEnd = -1
		var RowMergeCell = -1
		var RowMergeValue string
		for rowIndex, row := range sheet.Rows {

			if rowIndex == 0 {
				for _, cell := range row.Cells {
					title := cell.Value
					title = strings.TrimSpace(title)
					titles = append(titles, title)
				}
				continue
			}
			var dataType = map[string]string{}
			for cellIndex, cell := range row.Cells {
				if cellIndex >= len(titles) {
					break
				}
				title := titles[cellIndex]
				if title == "" {
					continue
				}
				value := cell.Value
				value = strings.TrimSpace(value)
				if cell.VMerge > 0 {
					RowMergeCell = cellIndex
					RowMergeEnd = rowIndex + cell.VMerge
					RowMergeValue = value
				}
				if RowMergeCell == cellIndex {
					if rowIndex <= RowMergeEnd {
						value = RowMergeValue
					} else {
						RowMergeEnd = -1
						RowMergeValue = ""
					}
				}
				dataType[title] = value
			}
			if dataType["名称"] == "" {
				continue
			}
			database.DataTypes = append(database.DataTypes, formatDataType(dataType))
		}

		databases = append(databases, database)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return
	}
	_, err = outFile.WriteString(`package dialect

import "strings"

`)
	if err != nil {
		return
	}
	for _, one := range databases {
		fmt.Println("-------- database [" + one.Name + "] start --------")

		var code string
		code += "// " + one.Name + " 数据库 字段类型" + "\n"
		columnTypeListName := ""
		var isMysql bool
		var isShenTong bool
		if strings.EqualFold(one.Name, "Mysql") {
			columnTypeListName = "mysqlColumnTypeList"
			isMysql = true
		} else if strings.EqualFold(one.Name, "Oracle") {
			columnTypeListName = "oracleColumnTypeList"
		} else if strings.EqualFold(one.Name, "达梦") {
			columnTypeListName = "dmColumnTypeList"
		} else if strings.EqualFold(one.Name, "金仓") {
			columnTypeListName = "kingBaseColumnTypeList"
		} else if strings.EqualFold(one.Name, "神通") {
			columnTypeListName = "shenTongColumnTypeList"
			isShenTong = true
		} else if strings.EqualFold(one.Name, "Sqlite") {
			columnTypeListName = "sqliteColumnTypeList"
		} else if strings.EqualFold(one.Name, "GBase") {
			columnTypeListName = "gBaseColumnTypeList"
		} else if strings.EqualFold(one.Name, "Postgresql") {
			columnTypeListName = "postgresqlColumnTypeList"
		} else if strings.EqualFold(one.Name, "DB2") {
			columnTypeListName = "db2ColumnTypeList"
		}
		code += "var " + columnTypeListName + " = []*ColumnTypeInfo{" + "\n"
		for _, dataType := range one.DataTypes {
			code += "\t" + "{"
			code += "Name: `" + dataType.Name + "`, "
			code += "Format: `" + dataType.Format + "`, "
			if len(dataType.Matches) > 0 {
				code += "Matches: []string{`" + strings.Join(dataType.Matches, "`, `") + "`}, "
			}
			if dataType.IsNumber {
				code += "IsNumber: true, "
			}
			if dataType.IsInteger {
				code += "IsInteger: true, "
			}
			if dataType.IsFloat {
				code += "IsFloat: true, "
			}
			if dataType.IsBoolean {
				code += "IsBoolean: true, "
			}
			if dataType.IsString {
				code += "IsString: true, "
			}
			if dataType.IsBytes {
				code += "IsBytes: true, "
			}
			if dataType.IsEnum {
				code += "IsEnum: true, "
			}
			if dataType.IsDateTime {
				code += "IsDateTime: true, "
			}
			if dataType.Comment != "" {
				code += "Comment: `" + dataType.Comment + "`, "
			}
			var hasOtherMethod bool
			if dataType.Name == "DATETIME" || dataType.Name == "TIMESTAMP" {
				if isShenTong {
				} else {
					code = strings.TrimSuffix(code, " ")
					hasOtherMethod = true
					code += `
		ColumnDefaultPack: func(param *ParamModel, column *ColumnModel) (columnDefaultPack string, err error) {
			if strings.Contains(strings.ToLower(column.ColumnDefault), "current_timestamp") ||
				strings.Contains(strings.ToLower(column.ColumnDefault), "0000-00-00 00:00:00") {
				columnDefaultPack = "CURRENT_TIMESTAMP"
			}
`
					if isMysql {
						code += `
			if strings.Contains(strings.ToLower(column.ColumnExtra), "on update current_timestamp") {
				columnDefaultPack += " ON UPDATE CURRENT_TIMESTAMP"
			}
`
					}
					code += `
			return
		},
`
				}
			} else if dataType.IsEnum {
				if isMysql {
					hasOtherMethod = true
					code = strings.TrimSuffix(code, " ")
					code += `
		FullColumnByColumnType: func(columnType string, column *ColumnModel) (err error) {
			if strings.Contains(columnType, "(") {
				setStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
				setStr = strings.ReplaceAll(setStr, "'", "")
				column.ColumnEnums = strings.Split(setStr, ",")
			}
			return
		},
`
				} else {
				}
			} else {

			}
			if hasOtherMethod {
				code += "\t" + "}," + "\n"
			} else {
				code = code[0 : len(code)-2]
				code += "}," + "\n"
			}
		}
		code += "}" + "\n\n"
		fmt.Println(code)
		_, err = outFile.WriteString(code)
		if err != nil {
			return
		}
		fmt.Println("-------- database [" + one.Name + "] end --------")
	}
	return
}

func formatDataType(dataType map[string]string) (info *dialect.ColumnTypeInfo) {
	info = &dialect.ColumnTypeInfo{}
	name := dataType["名称"]
	format := name
	if strings.Contains(name, "(") {
		nameStart := name[0:strings.Index(name, "(")]
		nameEnd := name[strings.Index(name, ")")+1:]
		inStr := name[strings.Index(name, "("):strings.Index(name, ")")]
		inStr = strings.ReplaceAll(inStr, "(", "")
		inStr = strings.ReplaceAll(inStr, ")", "")

		ss := strings.Split(inStr, ",")
		format = nameStart + "("
		for _, s := range ss {
			s = strings.TrimSpace(s)
			if strings.EqualFold(s, "p") ||
				strings.EqualFold(s, "precision") ||
				strings.Contains(s, "精度") {
				format += "$p, "
			} else if strings.EqualFold(s, "s") ||
				strings.EqualFold(s, "scale") ||
				strings.Contains(s, "标度") ||
				strings.Contains(s, "刻度") {
				format += "$s, "
			} else {
				format += "$l, "
			}
		}
		format = strings.TrimSuffix(format, ", ")

		format += ")"
		if strings.Contains(nameEnd, "(") {
			endStart := nameEnd[0:strings.Index(nameEnd, "(")]
			endEnd := nameEnd[strings.Index(nameEnd, ")")+1:]
			inStr = nameEnd[strings.Index(nameEnd, "("):strings.Index(nameEnd, ")")]
			inStr = strings.ReplaceAll(inStr, "(", "")
			inStr = strings.ReplaceAll(inStr, ")", "")

			ss = strings.Split(inStr, ",")
			format += endStart + "("
			for _, s := range ss {
				s = strings.TrimSpace(s)
				if strings.EqualFold(s, "p") ||
					strings.EqualFold(s, "precision") ||
					strings.Contains(s, "精度") {
					if strings.Contains(s, "小数秒精度") {
						format += "$s, "
					} else {
						format += "$p, "
					}
				} else if strings.EqualFold(s, "s") ||
					strings.EqualFold(s, "scale") ||
					strings.Contains(s, "标度") ||
					strings.Contains(s, "刻度") {
					format += "$s, "
				} else {
					format += "$l, "
				}
			}
			format = strings.TrimSuffix(format, ", ")
			format += ")" + endEnd
			name = nameStart + endStart + endEnd
		} else {
			format += nameEnd
			name = nameStart + nameEnd
		}

	}
	var typeText = dataType["类型"]
	if strings.Contains(typeText, "整型") {
		info.IsInteger = true
		info.IsNumber = true
	} else if strings.Contains(typeText, "浮点") {
		info.IsFloat = true
		info.IsNumber = true
	} else if strings.Contains(typeText, "定点") {
		info.IsFloat = true
		info.IsNumber = true
	} else if strings.Contains(typeText, "数值") {
		info.IsNumber = true
	} else if strings.Contains(typeText, "字符") {
		info.IsString = true
	} else if strings.Contains(typeText, "二进制") {
		info.IsBytes = true
	} else if strings.Contains(typeText, "布尔") {
		info.IsBoolean = true
	} else if strings.Contains(typeText, "日期") {
		info.IsDateTime = true
	} else if strings.Contains(typeText, "枚举") {
		info.IsEnum = true
	}
	info.Name = name
	info.Format = format
	info.Comment = dataType["说明"]
	matchStr := dataType["匹配"]
	matches := strings.Split(matchStr, "\n")
	for _, match := range matches {
		match = strings.TrimSpace(match)
		if match == "" {
			continue
		}
		if strings.EqualFold(match, "if not found") {
			info.IfNotFound = true
			continue
		}
		info.Matches = append(info.Matches, match)
	}
	return
}
