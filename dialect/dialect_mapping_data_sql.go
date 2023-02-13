package dialect

import (
	"errors"
	"strings"
)

func (this_ *mappingDialect) AppendSqlValue(param *ParamModel, sqlInfo *string, column *ColumnModel, value interface{}, args *[]interface{}) {
	if param != nil && param.AppendSqlValue != nil && *param.AppendSqlValue {
		*sqlInfo += this_.SqlValuePack(param, column, value)
	} else {
		*sqlInfo += "?"
		*args = append(*args, value)
	}
}

func (this_ *mappingDialect) DataListInsertSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, batchSqlList []string, batchValuesList [][]interface{}, err error) {
	if len(dataList) == 0 {
		return
	}

	var batchIndexCache = make(map[string]int)

	var columnCache = map[string]*ColumnModel{}
	for _, column := range columnList {
		columnCache[column.ColumnName] = column
	}

	for _, data := range dataList {

		var values []interface{}
		insertColumns := ""
		insertValues := ""
		for _, column := range columnList {
			name := column.ColumnName
			value, ok := data[name]
			if !ok {
				continue
			}
			insertColumns += this_.ColumnNamePack(param, name) + ", "
			this_.AppendSqlValue(param, &insertValues, column, value, &values)
			insertValues += ", "
		}
		insertColumns = strings.TrimSuffix(insertColumns, ", ")
		insertValues = strings.TrimSuffix(insertValues, ", ")

		var insertSql = "INSERT INTO "
		insertSql += this_.OwnerTablePack(param, ownerName, tableName)
		if insertColumns != "" {
			insertSql += "(" + insertColumns + ")"
		}
		if insertValues != "" {
			insertSql += " VALUES (" + insertValues + ")"
		}

		sqlList = append(sqlList, this_.ReplaceSqlVariable(insertSql, values))
		valuesList = append(valuesList, values)

		// 批量 插入 SQL

		if this_.dialectType == TypeOracle {
			batchSqlList = append(batchSqlList, insertSql)
			batchValuesList = append(batchValuesList, values)
		} else {
			index, ok := batchIndexCache[insertColumns]
			if ok {
				batchSqlList[index] += ",\n(" + insertValues + ")"
				batchValuesList[index] = append(batchValuesList[index], values...)
			} else {
				var batchSql = "INSERT INTO "
				batchSql += this_.OwnerTablePack(param, ownerName, tableName)
				if insertColumns != "" {
					batchSql += "(" + insertColumns + ")"
				}
				if insertValues != "" {
					batchSql += " VALUES (" + insertValues + ")"
				}
				index = len(batchSqlList)
				batchIndexCache[insertColumns] = index
				batchSqlList = append(batchSqlList, batchSql)
				batchValuesList = append(batchValuesList, values)
			}
		}

	}
	for index := range batchSqlList {
		batchSqlList[index] = this_.ReplaceSqlVariable(batchSqlList[index], batchValuesList[index])
	}
	return
}

func (this_ *mappingDialect) DataListUpdateSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataList) == 0 {
		return
	}
	if len(dataList) != len(dataWhereList) {
		err = errors.New("更新数据与更新条件数量不一致")
		return
	}
	var columnCache = map[string]*ColumnModel{}
	for _, column := range columnList {
		columnCache[column.ColumnName] = column
	}

	for index, data := range dataList {
		dataWhere := dataWhereList[index]
		if len(dataWhere) == 0 {
			err = errors.New("更新数据条件丢失")
			return
		}

		var updateSql = "UPDATE "
		var values []interface{}

		updateSql += this_.OwnerTablePack(param, ownerName, tableName)
		updateSql += " SET "

		for name, value := range data {
			column := columnCache[name]
			updateSql += "" + this_.ColumnNamePack(param, name) + "="
			this_.AppendSqlValue(param, &updateSql, column, value, &values)
			updateSql += ", "
		}
		updateSql = strings.TrimSuffix(updateSql, ", ")

		updateSql += " WHERE "
		for name, value := range dataWhere {
			column := columnCache[name]
			updateSql += "" + this_.ColumnNamePack(param, name) + "="
			this_.AppendSqlValue(param, &updateSql, column, value, &values)
			updateSql += " AND "
		}
		updateSql = strings.TrimSuffix(updateSql, " AND ")

		updateSql = this_.ReplaceSqlVariable(updateSql, values)
		sqlList = append(sqlList, updateSql)
		valuesList = append(valuesList, values)
	}
	return
}
func (this_ *mappingDialect) DataListDeleteSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataWhereList) == 0 {
		return
	}
	var columnCache = map[string]*ColumnModel{}
	for _, column := range columnList {
		columnCache[column.ColumnName] = column
	}

	for _, dataWhere := range dataWhereList {
		if len(dataWhere) == 0 {
			err = errors.New("更新数据条件丢失")
			return
		}

		var deleteSql = "DELETE FROM "
		var values []interface{}

		deleteSql += this_.OwnerTablePack(param, ownerName, tableName)

		deleteSql += " WHERE "

		for name, value := range dataWhere {
			column := columnCache[name]
			deleteSql += "" + this_.ColumnNamePack(param, name) + "="
			this_.AppendSqlValue(param, &deleteSql, column, value, &values)
			deleteSql += " AND "
		}
		deleteSql = strings.TrimSuffix(deleteSql, " AND ")

		deleteSql = this_.ReplaceSqlVariable(deleteSql, values)
		sqlList = append(sqlList, deleteSql)
		valuesList = append(valuesList, values)
	}
	return
}

type Where struct {
	Name                    string `json:"name"`
	Value                   string `json:"value"`
	Before                  string `json:"before"`
	After                   string `json:"after"`
	CustomSql               string `json:"customSql"`
	SqlConditionalOperation string `json:"sqlConditionalOperation"`
	AndOr                   string `json:"andOr"`
}

type Order struct {
	Name    string `json:"name"`
	AscDesc string `json:"ascDesc"`
}

func (this_ *mappingDialect) DataListSelectSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, whereList []*Where, orderList []*Order) (sql string, values []interface{}, err error) {
	selectColumns := ""
	var columnCache = map[string]*ColumnModel{}
	for _, column := range columnList {
		selectColumns += this_.ColumnNamePack(param, column.ColumnName) + ","
		columnCache[column.ColumnName] = column
	}
	selectColumns = strings.TrimSuffix(selectColumns, ",")
	if selectColumns == "" {
		selectColumns = "*"
	}
	sql = "SELECT " + selectColumns + " FROM "

	sql += this_.OwnerTablePack(param, ownerName, tableName)

	//构造查询用的finder
	if len(whereList) > 0 {
		sql += " WHERE"
		for index, where := range whereList {
			column := columnCache[where.Name]
			sql += " " + this_.ColumnNamePack(param, where.Name)
			value := where.Value
			switch where.SqlConditionalOperation {
			case "like":
				sql += " LIKE "
				value = "%" + value + "%"
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "not like":
				sql += " NOT LIKE "
				value = "%" + value + "%"
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "like start":
				sql += " LIKE "
				value = "" + value + "%"
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "not like start":
				sql += " NOT LIKE "
				value = "" + value + "%"
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "like end":
				sql += " LIKE "
				value = "%" + value + ""
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "not like end":
				sql += " NOT LIKE "
				value = "%" + value + ""
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "is null":
				sql += " IS NULL"
			case "is not null":
				sql += " IS NOT NULL"
			case "is empty":
				sql += " = "
				value = ""
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "is not empty":
				sql += " <> "
				this_.AppendSqlValue(param, &sql, column, value, &values)
			case "between":
				sql += " BETWEEN "
				this_.AppendSqlValue(param, &sql, column, where.Before, &values)
				sql += " AND "
				this_.AppendSqlValue(param, &sql, column, where.After, &values)
			case "not between":
				sql += " NOT BETWEEN "
				this_.AppendSqlValue(param, &sql, column, where.Before, &values)
				sql += " AND "
				this_.AppendSqlValue(param, &sql, column, where.After, &values)
			case "in":
				sql += " IN ("
				this_.AppendSqlValue(param, &sql, column, value, &values)
				sql += ")"
			case "not in":
				sql += " NOT IN ("
				this_.AppendSqlValue(param, &sql, column, value, &values)
				sql += ")"
			default:
				sql += " " + where.SqlConditionalOperation + " "
				this_.AppendSqlValue(param, &sql, column, value, &values)
			}
			// params_ = append(params_, where.Value)
			if index < len(whereList)-1 {
				sql += " " + where.AndOr + " "
			}
		}
	}
	if len(orderList) > 0 {
		sql += " ORDER BY"
		for index, order := range orderList {
			sql += " " + this_.ColumnNamePack(param, order.Name)
			if order.AscDesc != "" {
				sql += " " + order.AscDesc
			}
			// params_ = append(params_, where.Value)
			if index < len(orderList)-1 {
				sql += ","
			}
		}

	}
	sql = this_.ReplaceSqlVariable(sql, values)
	return
}
