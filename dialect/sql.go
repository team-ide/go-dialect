package dialect

import (
	"errors"
	"regexp"
	"strings"
)

func FormatCountSql(selectSql string) (countSql string, err error) {
	countSql = strings.TrimSpace(selectSql)
	if countSql == "" {
		return
	}

	//查询order by 的位置
	//Query the position of order by
	locOrderBy := findOrderByIndex(countSql)
	//如果存在order by
	//If there is order by
	if len(locOrderBy) > 0 {
		countSql = countSql[:locOrderBy[0]]
	}
	s := strings.ToLower(countSql)
	gbi := -1
	locGroupBy := findGroupByIndex(countSql)
	if len(locGroupBy) > 0 {
		gbi = locGroupBy[0]
	}
	var sqlBuilder strings.Builder
	sqlBuilder.Grow(100)
	//特殊关键字,包装SQL
	//Special keywords, wrap SQL
	if strings.Contains(s, " distinct ") || strings.Contains(s, " union ") || gbi > -1 {
		sqlBuilder.WriteString("SELECT COUNT(*)  frame_row_count FROM (")
		sqlBuilder.WriteString(countSql)
		sqlBuilder.WriteString(") temp_frame_noob_table_name WHERE 1=1 ")
	} else {
		locFrom := findSelectFromIndex(countSql)
		//没有找到FROM关键字,认为是异常语句
		//The FROM keyword was not found, which is considered an abnormal statement
		if len(locFrom) == 0 {
			err = errors.New("->selectCount-->findFromIndex没有FROM关键字,语句错误")
			return
		}
		sqlBuilder.WriteString("SELECT COUNT(*) ")
		sqlBuilder.WriteString(countSql[locFrom[0]:])
	}
	countSql = sqlBuilder.String()
	return
}

var orderByExpr = "(?i)\\s(order)\\s+by\\s"
var orderByRegexp, _ = regexp.Compile(orderByExpr)

func findOrderByIndex(selectSql string) []int {
	loc := orderByRegexp.FindStringIndex(selectSql)
	return loc
}

var groupByExpr = "(?i)\\s(group)\\s+by\\s"
var groupByRegexp, _ = regexp.Compile(groupByExpr)

func findGroupByIndex(selectSql string) []int {
	loc := groupByRegexp.FindStringIndex(selectSql)
	return loc
}

var fromExpr = "(?i)(^\\s*select)(\\(.*?\\)|[^()]+)*?(from)"
var fromRegexp, _ = regexp.Compile(fromExpr)

func findSelectFromIndex(selectSql string) []int {
	//匹配出来的是完整的字符串,用最后的FROM即可
	loc := fromRegexp.FindStringIndex(selectSql)
	if len(loc) < 2 {
		return loc
	}
	//最后的FROM前推4位字符串
	loc[0] = loc[1] - 4
	return loc
}
